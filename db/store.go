package db

import (
	"context"
	"database/sql"
	"filesearch/model"
)

type Store interface {
	CreateUserTx(ctx context.Context, arg model.User) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}
type SQLStore struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
	}
}

var txKey = struct{}{}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
