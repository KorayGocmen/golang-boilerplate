package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/koraygocmen/null"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/koraygocmen/golang-boilerplate/pkg/generate"
)

type GivenNames string
type Surname string
type VerificationStatus string
type Status string

var (
	VerificationStatusPending    VerificationStatus = "PENDING"
	VerificationStatusVerified   VerificationStatus = "VERIFIED"
	VerificationStatusRejected   VerificationStatus = "REJECTED"
	VerificationStatusInProgress VerificationStatus = "IN_PROGRESS"

	VerificationStatuses = map[VerificationStatus]bool{
		VerificationStatusPending:    true,
		VerificationStatusVerified:   true,
		VerificationStatusRejected:   true,
		VerificationStatusInProgress: true,
	}

	AccountStatusActive    Status = "ACTIVE"
	AccountStatusSuspended Status = "SUSPENDED"
	AccountStatusDeleted   Status = "DELETED"

	AccountStatuses = map[Status]bool{
		AccountStatusActive:    true,
		AccountStatusSuspended: true,
		AccountStatusDeleted:   true,
	}

	AccountStatusOK = map[Status]bool{
		AccountStatusActive:    true,
		AccountStatusSuspended: false,
		AccountStatusDeleted:   false,
	}
)

func ToGivenNames(givenNames string) GivenNames {
	return GivenNames(strings.ToUpper(strings.TrimSpace(givenNames)))
}

func ToSurname(surname string) Surname {
	return Surname(strings.ToUpper(strings.TrimSpace(surname)))
}

type User struct {
	ID            int64          `gorm:"type:integer; primaryKey;" json:"id"`
	CreatedAt     time.Time      `gorm:"type:timestamp; autoCreateTime;" json:"createdAt"`
	DeletedAt     gorm.DeletedAt `gorm:"type:timestamp; index;" json:"-"`
	Email         null.String    `gorm:"type:text; index;" json:"email"`
	EmailVerified null.Bool      `gorm:"type:boolean; default:false;" json:"emailVerified"`
	Password      null.String    `gorm:"-" json:"-"`
	PasswordHash  null.String    `gorm:"type:text; uniqueIndex;" json:"-"`
	PasswordSet   bool           `gorm:"-" json:"passwordSet"`
	GivenNames    null.String    `gorm:"type:text; index;" json:"givenNames"`
	Surname       null.String    `gorm:"type:text; index;" json:"surname"`
}

// Gorm hooks.

// AfterFind gorm hook.
func (u *User) AfterFind(tx *gorm.DB) error {
	return u.after(tx)
}

// AfterCreate gorm hook.
func (u *User) AfterCreate(tx *gorm.DB) error {
	return u.after(tx)
}

// Internal.
func (u *User) after(tx *gorm.DB) error {
	if u == nil {
		return nil
	}

	u.PasswordSet = u.PasswordHash.Valid
	return nil
}

// Model methods.

// Create the password hash from the password.
func (u *User) PasswordHashCreate() error {
	if !u.Password.Valid {
		err := fmt.Errorf("password hash create error: user password is required to create password hash")
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password.String), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("password hash create error: bcrypt generate from password error: %w", err)
		return err
	}

	u.PasswordHash = null.StringFrom(string(passwordHash))
	return nil
}

// Compare the password hash with the password.
func (u *User) PasswordHashCompare(password string) bool {
	if !u.PasswordHash.Valid {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash.String), []byte(password)) == nil
}

func (User) ReferralCodeCreate() string {
	// User referral code is alphanumeric and 5 characters.
	return generate.AlphaCode(5, true)
}
