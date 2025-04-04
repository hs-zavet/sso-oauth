package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/app/models"
	"github.com/hs-zavet/sso-oauth/internal/repo"
	"github.com/hs-zavet/tokens/identity"
)

func (a App) AccountCreate(ctx context.Context, email string) error {
	ID := uuid.New()
	CreatedAt := time.Now().UTC()

	err := a.accounts.Create(ctx, repo.AccountCreateRequest{
		ID:           ID,
		Email:        email,
		Role:         string(identity.User),
		Subscription: uuid.Nil,
		CreatedAt:    CreatedAt,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a App) AccountUpdateSubscription(ctx context.Context, ID uuid.UUID, subscriptionID uuid.UUID) error {
	UpdatedAt := time.Now().UTC()

	err := a.accounts.Update(ctx, ID, repo.AccountUpdateRequest{
		Subscription: subscriptionID,
		UpdatedAt:    UpdatedAt,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a App) AccountUpdateRole(ctx context.Context, ID uuid.UUID, role string) error {
	UpdatedAt := time.Now().UTC()

	err := a.accounts.Update(ctx, ID, repo.AccountUpdateRequest{
		Role:      role,
		UpdatedAt: UpdatedAt,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a App) AccountGetByID(ctx context.Context, ID uuid.UUID) (models.Account, error) {
	account, err := a.accounts.GetByID(ctx, ID)
	if err != nil {
		return models.Account{}, err
	}

	role, err := identity.ParseIdentityType(account.Role)
	if err != nil {
		return models.Account{}, err
	}

	return models.Account{
		ID:           account.ID,
		Email:        account.Email,
		Role:         role,
		Subscription: account.Subscription,
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    account.UpdatedAt,
	}, nil
}

func (a App) AccountGetByEmail(ctx context.Context, email string) (models.Account, error) {
	account, err := a.accounts.GetByEmail(ctx, email)
	if err != nil {
		return models.Account{}, err
	}

	role, err := identity.ParseIdentityType(account.Role)
	if err != nil {
		return models.Account{}, err
	}

	return models.Account{
		ID:           account.ID,
		Email:        account.Email,
		Role:         role,
		Subscription: account.Subscription,
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    account.UpdatedAt,
	}, nil
}
