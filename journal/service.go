// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package journal

import (
	"context"

	"github.com/absmach/magistrala"
	"github.com/absmach/magistrala/pkg/auth"
	"github.com/absmach/magistrala/pkg/errors"
	svcerr "github.com/absmach/magistrala/pkg/errors/service"
	"github.com/absmach/magistrala/pkg/policies"
)

type service struct {
	idProvider magistrala.IDProvider
	auth       auth.AuthClient
	repository Repository
}

func NewService(idp magistrala.IDProvider, repository Repository, authClient auth.AuthClient) Service {
	return &service{
		idProvider: idp,
		auth:       authClient,
		repository: repository,
	}
}

func (svc *service) Save(ctx context.Context, journal Journal) error {
	id, err := svc.idProvider.ID()
	if err != nil {
		return err
	}
	journal.ID = id

	return svc.repository.Save(ctx, journal)
}

func (svc *service) RetrieveAll(ctx context.Context, token string, page Page) (JournalsPage, error) {
	if err := svc.authorize(ctx, token, page.EntityID, page.EntityType.AuthString()); err != nil {
		return JournalsPage{}, err
	}

	return svc.repository.RetrieveAll(ctx, page)
}

func (svc *service) authorize(ctx context.Context, token, entityID, entityType string) error {
	user, err := svc.auth.Identify(ctx, &magistrala.IdentityReq{Token: token})
	if err != nil {
		return errors.Wrap(svcerr.ErrAuthentication, err)
	}

	permission := policies.ViewPermission
	objectType := entityType
	object := entityID
	subject := user.GetId()

	// If the entity is a user, we need to check if the user is an admin
	if entityType == policies.UserType {
		permission = policies.AdminPermission
		objectType = policies.PlatformType
		object = policies.MagistralaObject
		subject = user.GetUserId()
	}

	req := &magistrala.AuthorizeReq{
		Domain:      user.GetDomainId(),
		SubjectType: policies.UserType,
		SubjectKind: policies.UsersKind,
		Subject:     subject,
		Permission:  permission,
		ObjectType:  objectType,
		Object:      object,
	}

	res, err := svc.auth.Authorize(ctx, req)
	if err != nil {
		return errors.Wrap(svcerr.ErrAuthorization, err)
	}
	if !res.GetAuthorized() {
		return svcerr.ErrAuthorization
	}

	return nil
}
