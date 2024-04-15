package repo

import (
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/database"
	UserRepo "github.com/koraygocmen/golang-boilerplate/internal/repo/user"
	UserSessionRepo "github.com/koraygocmen/golang-boilerplate/internal/repo/user_session"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Transaction creates a db transaction and all
// repos are created with this transaction to
// ensure that all operations are done in the
// same transaction.
type Transaction struct {
	tx *gorm.DB

	Commit   func() error
	Rollback func() error
	Ping     func() error

	User        *UserRepo.Repo
	UserSession *UserSessionRepo.Repo
}

func New(ctx context.Ctx) *Transaction {
	tx := database.DB.GORM.Session(&gorm.Session{Context: ctx}).Begin()

	transaction := &Transaction{
		tx: tx,

		User:        UserRepo.New(tx),
		UserSession: UserSessionRepo.New(tx),
	}

	// Set the commit, rollback and ping functions.
	transaction.Commit = commit(transaction)
	transaction.Rollback = rollback(transaction)
	transaction.Ping = ping(transaction)

	return transaction
}

func commit(transaction *Transaction) func() error {
	return func() error {
		return transaction.tx.Commit().Error
	}
}

func rollback(transaction *Transaction) func() error {
	return func() error {
		return transaction.tx.Rollback().Error
	}
}

func ping(transaction *Transaction) func() error {
	return func() error {
		return transaction.tx.
			Session(&gorm.Session{Logger: gormlogger.Default.LogMode(gormlogger.Silent)}).
			Exec("SELECT 1;").
			Error
	}
}
