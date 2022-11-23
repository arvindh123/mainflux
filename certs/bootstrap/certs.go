// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"context"
	"net/http"

	"github.com/mainflux/mainflux/bootstrap"
	"github.com/mainflux/mainflux/certs"
	"github.com/mainflux/mainflux/pkg/errors"

	"os"
	"sync"

	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
)

var (
	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("unauthorized access to Bootstrap service")

	// ErrUnexpectedBSResponse indicates unexpected response from Bootstrap service.
	ErrUnexpectedBSResponse = errors.New("unexpected Bootstrap service response")

	//ErrUnableToAccess indicates boostrap service is not accessible
	ErrUnableToAccess = errors.New("unable to access bootstrap service")

	ErrFailedToLogin = errors.New("failed to login")

	ErrFailedToReadResponseBody = errors.New("failed to read Bootstrap response body ")
)

type bootstrapClient struct {
	updateURL string
	token     string
	email     string
	pass      string
	sdk       mfsdk.SDK
	mu        sync.Mutex
}

// New returns new Bootstrap service client
func New(updateURL, email, pass string, sdk mfsdk.SDK) certs.BootstrapClient {
	token := os.Getenv("MF_USERS_TOKEN")

	return &bootstrapClient{
		updateURL: updateURL,
		email:     email,
		pass:      pass,
		sdk:       sdk,
		token:     token,
		mu:        sync.Mutex{},
	}
}

func (c *bootstrapClient) fetchToken() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.token
}

func (c *bootstrapClient) login() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	user := mfsdk.User{Email: c.email, Password: c.pass}
	token, err := c.sdk.CreateToken(user)
	if err != nil {
		return "", err
	}
	c.token = token

	return c.token, nil
}

func (c *bootstrapClient) UpdateCerts(ctx context.Context, thingID, clientCert, clientKey, caCert string) error {
	if err := c.sdk.UpdateBootstrapCerts(c.fetchToken(), thingID, clientCert, clientKey, caCert); err != nil {
		err := sdkError(err)
		if err != ErrUnauthorizedAccess {
			return err
		}
		if _, err := c.login(); err != nil {
			return errors.Wrap(ErrFailedToLogin, err)
		}
		return c.sdk.UpdateBootstrapCerts(c.fetchToken(), thingID, clientCert, clientKey, caCert)
	}
	return nil
}

func sdkError(err error) error {
	if sdkErr, ok := err.(errors.SDKError); ok {
		switch sdkErr.StatusCode() {
		case http.StatusForbidden, http.StatusUnauthorized:
			return ErrUnauthorizedAccess
		case http.StatusNotFound:
			if errors.Contains(sdkErr, bootstrap.ErrUpdateCert) {
				return errors.ErrNotFound
			}
		}
	}
	return errors.Wrap(ErrUnexpectedBSResponse, err)
}
