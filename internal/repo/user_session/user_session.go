package user_session

import (
	"errors"
	"fmt"
	"time"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	UserSession "github.com/koraygocmen/golang-boilerplate/internal/model/user_session"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Function types.
type CreateFn func(ctx context.Ctx, userSession *UserSession.UserSession) error
type GetByIDFn func(ctx context.Ctx, id int64) (*UserSession.UserSession, error)
type ListActiveFn func(ctx context.Ctx, userID int64) ([]*UserSession.UserSession, error)
type DeleteFn func(ctx context.Ctx, userSession *UserSession.UserSession) error

// Repo.
// Repo definition and repo related fields.
type Repo struct {
	Create     CreateFn
	GetByID    GetByIDFn
	ListActive ListActiveFn
	Delete     DeleteFn
}

func New(tx *gorm.DB) *Repo {
	return &Repo{
		Create:     create(tx),
		GetByID:    getByID(tx),
		ListActive: listActive(tx),
		Delete:     delete(tx),
	}
}

// Functions.
func create(tx *gorm.DB) CreateFn {
	return func(ctx context.Ctx, userSession *UserSession.UserSession) error {
		err := tx.WithContext(ctx).
			Omit(clause.Associations).
			Create(userSession).
			Error
		if err != nil {
			err = fmt.Errorf("user session repo create error: %w", err)
			return err
		}
		return nil
	}
}

func getByID(tx *gorm.DB) GetByIDFn {
	return func(ctx context.Ctx, id int64) (*UserSession.UserSession, error) {
		var userSession UserSession.UserSession
		err := tx.WithContext(ctx).
			Where(`"id" = ?`, id).
			First(&userSession).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("user session repo get by id error: %w", err)
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return &userSession, nil
	}
}

func listActive(tx *gorm.DB) ListActiveFn {
	return func(ctx context.Ctx, userID int64) ([]*UserSession.UserSession, error) {
		var userSessions []*UserSession.UserSession
		err := tx.WithContext(ctx).
			Where(`"user_id" = ?`, userID).
			Where(`"expire_at" > ?`, time.Now().UTC()).
			Find(&userSessions).
			Error
		if err != nil {
			err = fmt.Errorf("user session repo list active error: %w", err)
			return nil, err
		}

		return userSessions, nil
	}
}

func delete(tx *gorm.DB) DeleteFn {
	return func(ctx context.Ctx, userSession *UserSession.UserSession) error {
		err := tx.WithContext(ctx).
			Omit(clause.Associations).
			Delete(userSession).
			Error
		if err != nil {
			err = fmt.Errorf("user session repo delete error: %w", err)
			return err
		}
		return nil
	}
}
