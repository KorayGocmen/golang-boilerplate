package user_repo

import (
	"errors"
	"fmt"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	User "github.com/koraygocmen/golang-boilerplate/internal/model/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Function types.
type CreateFn func(ctx context.Ctx, user *User.User) error
type SaveFn func(ctx context.Ctx, user *User.User) error
type ListFn func(ctx context.Ctx, pageSize, pageNum int) ([]*User.User, error)
type TotalFn func(ctx context.Ctx) (int64, error)
type GetByIDFn func(ctx context.Ctx, id int64) (*User.User, error)
type GetByEmailFn func(ctx context.Ctx, email string) (*User.User, error)

// Repo.
// Repo definition and repo related fields.
type Repo struct {
	Create     CreateFn
	Save       SaveFn
	List       ListFn
	Total      TotalFn
	GetByID    GetByIDFn
	GetByEmail GetByEmailFn
}

func New(tx *gorm.DB) *Repo {
	return &Repo{
		Create:     create(tx),
		Save:       save(tx),
		List:       list(tx),
		Total:      total(tx),
		GetByID:    getByID(tx),
		GetByEmail: getByEmail(tx),
	}
}

// Functions.
func create(tx *gorm.DB) CreateFn {
	return func(ctx context.Ctx, user *User.User) error {
		err := tx.WithContext(ctx).
			Omit(clause.Associations).
			Create(user).
			Error
		if err != nil {
			err = fmt.Errorf("user repo create error: %w", err)
			return err
		}
		return nil
	}
}

func save(tx *gorm.DB) SaveFn {
	return func(ctx context.Ctx, user *User.User) error {
		err := tx.WithContext(ctx).
			Omit(clause.Associations).
			Save(user).
			Error
		if err != nil {
			err = fmt.Errorf("user repo save error: %w", err)
			return err
		}
		return nil
	}
}

func list(tx *gorm.DB) ListFn {
	return func(ctx context.Ctx, pageSize, pageNum int) ([]*User.User, error) {
		var users []*User.User
		err := tx.WithContext(ctx).
			Order(`"id" ASC`).
			Limit(pageSize).
			Offset(pageSize * (pageNum - 1)).
			Find(&users).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("user repo list error: %w", err)
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*User.User{}, nil
		}

		return users, nil
	}
}

func total(tx *gorm.DB) TotalFn {
	return func(ctx context.Ctx) (int64, error) {
		var total int64
		err := tx.WithContext(ctx).
			Model(&User.User{}).
			Count(&total).
			Error
		if err != nil {
			err = fmt.Errorf("user repo total error: %w", err)
			return 0, err
		}
		return total, nil
	}
}

func getByID(tx *gorm.DB) GetByIDFn {
	return func(ctx context.Ctx, id int64) (*User.User, error) {
		var user User.User
		err := tx.WithContext(ctx).
			Where("id = ?", id).
			First(&user).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("user repo get by id error: %w", err)
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return &user, nil
	}
}

func getByEmail(tx *gorm.DB) GetByEmailFn {
	return func(ctx context.Ctx, email string) (*User.User, error) {
		var user User.User
		err := tx.WithContext(ctx).
			Where(`"email" = ?`, email).
			First(&user).
			Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("user repo get by email error: %w", err)
			return nil, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return &user, nil
	}
}
