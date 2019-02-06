package repository

import (
	"database/sql"
)

type ar struct {
	db *sql.DB
}

// AccountRepository - interfice to provide operations with account models
type AccountRepository interface {
}

// NewAccountRepository - create new account repository object
func NewAccountRepository(db *sql.DB) AccountRepository {
	ar := ar{
		db: db,
	}

	return &ar
}
