package service

import (
	"fmt"
	"time"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/logger"
	"github.com/koraygocmen/golang-boilerplate/internal/repo"

	UserService "github.com/koraygocmen/golang-boilerplate/internal/service/user"
	UserSessionService "github.com/koraygocmen/golang-boilerplate/internal/service/user_session"
)

// Service.
type ServiceType struct {
	Transaction func(ctx context.Ctx, timeout time.Duration) *Transaction
}

var (
	Service = &ServiceType{
		Transaction: transaction(),
	}
)

// Transaction.
type Transaction struct {
	tx    *repo.Transaction
	timer *time.Timer
	start time.Time
	end   time.Time
	err   error

	Commit   func() error
	Rollback func(err error) error
	Ping     func() error

	User        *UserService.Service
	UserSession *UserSessionService.Service
}

func transaction() func(ctx context.Ctx, timeout time.Duration) *Transaction {
	return func(ctx context.Ctx, timeout time.Duration) *Transaction {
		tx := repo.New(ctx)

		userService := UserService.New(tx)
		userSessionService := UserSessionService.New(tx, userService)

		transaction := &Transaction{
			tx: tx,

			User:        userService,
			UserSession: userSessionService,
		}

		// Set the commit, rollback and ping functions.
		transaction.Commit = commit(transaction)
		transaction.Rollback = rollback(transaction)
		transaction.Ping = ping(transaction)

		// Start the transaction.
		transaction.start = time.Now().UTC()
		transaction.timer = time.AfterFunc(timeout, func() {
			err := fmt.Errorf("service transaction timeout error: transaction timed out after: %v", timeout)
			logger.Logger.Emerf(ctx, `err="%v"`, err)
			ctx = context.WithValue(ctx, context.KeyError, err)
			transaction.Rollback(err)
		})

		return transaction
	}
}

func commit(transaction *Transaction) func() error {
	return func() error {
		if !transaction.end.IsZero() {
			err := fmt.Errorf("service transaction commit error: transaction already ended: %w", transaction.err)
			return err
		}

		transaction.timer.Stop()
		transaction.end = time.Now().UTC()
		return transaction.tx.Commit()
	}
}

func rollback(transaction *Transaction) func(err error) error {
	return func(err error) error {
		if !transaction.end.IsZero() {
			err := fmt.Errorf("service transaction rollback error: transaction already ended: %w", transaction.err)
			return err
		}

		transaction.timer.Stop()
		transaction.end = time.Now().UTC()
		transaction.err = err
		return transaction.tx.Rollback()
	}
}

func ping(transaction *Transaction) func() error {
	return func() error {
		return transaction.tx.Ping()
	}
}
