package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hs-zavet/sso-oauth/internal/config"
	"github.com/hs-zavet/sso-oauth/internal/repo/sqldb"
)

type Account struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Subscription uuid.UUID `json:"subscription,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type accountSQL interface {
	New() sqldb.AccountQ
	Insert(ctx context.Context, input sqldb.AccountInsertInput) error
	Delete(ctx context.Context) error
	Count(ctx context.Context) (int, error)
	Select(ctx context.Context) ([]sqldb.AccountModel, error)
	Get(ctx context.Context) (sqldb.AccountModel, error)

	FilterID(id uuid.UUID) sqldb.AccountQ
	FilterEmail(email string) sqldb.AccountQ
	FilterRole(role string) sqldb.AccountQ
	FilterSubscription(subscription uuid.UUID) sqldb.AccountQ

	Update(ctx context.Context, input sqldb.AccountUpdateInput) error

	Page(limit, offset uint64) sqldb.AccountQ
	Transaction(fn func(ctx context.Context) error) error
}

type AccountsRepo struct {
	sql sqldb.AccountQ
}

func NewAccounts(cfg *config.Config) (AccountsRepo, error) {
	db, err := sql.Open("postgres", cfg.Database.SQL.URL)
	if err != nil {
		return AccountsRepo{}, err
	}

	accounts := sqldb.NewAccounts(db)

	return AccountsRepo{
		sql: accounts,
	}, nil
}

type AccountCreateRequest struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	Subscription uuid.UUID `json:"subscription,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func (a AccountsRepo) Create(ctx context.Context, input AccountCreateRequest) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	err := a.sql.New().Insert(ctxSync, sqldb.AccountInsertInput{
		ID:           input.ID,
		Email:        input.Email,
		Role:         input.Role,
		Subscription: input.Subscription,
		CreatedAt:    input.CreatedAt,
	})
	if err != nil {
		return err
	}
	return nil
}

type AccountUpdateRequest struct {
	Role         string    `json:"role"`
	Subscription uuid.UUID `json:"subscription,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (a AccountsRepo) Update(ctx context.Context, ID uuid.UUID, input AccountUpdateRequest) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	err := a.sql.New().FilterID(ID).Update(ctxSync, sqldb.AccountUpdateInput{
		Role:         input.Role,
		Subscription: input.Subscription,
		UpdatedAt:    input.UpdatedAt,
	})
	if err != nil {
		return err
	}
	return nil
}

func (a AccountsRepo) Delete(ctx context.Context, ID uuid.UUID) error {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	err := a.sql.New().FilterID(ID).Delete(ctxSync)
	if err != nil {
		return err
	}
	return nil
}

func (a AccountsRepo) GetByID(ctx context.Context, ID uuid.UUID) (Account, error) {
	ctxSync, cancel := context.WithTimeout(ctx, dataCtxTimeAisle)
	defer cancel()

	account, err := a.sql.New().FilterID(ID).Get(ctxSync)
	if err != nil {
		return Account{}, err
	}

	res := Account{
		ID:           account.ID,
		Email:        account.Email,
		Role:         account.Role,
		Subscription: account.Subscription,
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    *account.UpdatedAt,
	}

	if account.UpdatedAt != nil {
		res.UpdatedAt = *account.UpdatedAt
	}

	return res, nil
}

func (a AccountsRepo) GetByEmail(ctx context.Context, email string) (Account, error) {
	ctxSync, cancel := context.WithTimeout(context.Background(), dataCtxTimeAisle)
	defer cancel()

	account, err := a.sql.New().FilterEmail(email).Get(ctxSync)
	if err != nil {
		return Account{}, err
	}

	res := Account{
		ID:           account.ID,
		Email:        account.Email,
		Role:         account.Role,
		Subscription: account.Subscription,
		CreatedAt:    account.CreatedAt,
	}

	if account.UpdatedAt != nil {
		res.UpdatedAt = *account.UpdatedAt
	}

	return res, nil
}

func (a AccountsRepo) Transaction(fn func(ctx context.Context) error) error {
	return a.sql.Transaction(fn)
}
