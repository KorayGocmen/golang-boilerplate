package user_session

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/koraygocmen/golang-boilerplate/pkg/generate"
	"github.com/koraygocmen/null"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Purpose string

var (
	PurposeSessionCreate Purpose = "SESSION_CREATE"
	PurposePasswordReset Purpose = "PASSWORD_RESET"
	PurposeImpersonate   Purpose = "IMPERSONATE"

	Purposes = map[Purpose]bool{
		PurposeSessionCreate: true,
		PurposePasswordReset: true,
		// Don't allow impersonation by non-admins.
	}
)

func ToPurpose(purpose string) Purpose {
	return Purpose(strings.ToUpper(strings.TrimSpace(purpose)))
}

type UserSession struct {
	ID        int64          `gorm:"type:integer; primaryKey;" json:"id"`
	CreatedAt time.Time      `gorm:"type:timestamp; autoCreateTime;" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp; index;" json:"-"`
	UserID    int64          `gorm:"type:integer; not null; index;" json:"userId"`
	Token     string         `gorm:"-" json:"token"`
	TokenHash string         `gorm:"type:text; not null; uniqueIndex;" json:"-"`
	ClientIP  string         `gorm:"type:text; not null; index;" json:"-"`
	Purpose   Purpose        `gorm:"type:text; index;" json:"purpose"`
	ExpireAt  null.Time      `gorm:"type:timestamp;" json:"-"`
}

// TokenHashCreate creates a token hash from the seed provided.
func (u *UserSession) TokenHashCreate() error {
	seed := generate.AlphaCode(25, false)

	tokenSeed, err := bcrypt.GenerateFromPassword([]byte(seed), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("bcrypt generate from seed error: %w", err)
		return err
	}

	hashWriter := sha256.New()
	hashWriter.Write(tokenSeed)
	token := fmt.Sprintf("%x", hashWriter.Sum(nil))

	tokenHashRaw, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("bcrypt generate from token error: %w", err)
		return err
	}
	tokenHash := string(tokenHashRaw)

	u.Token = token
	u.TokenHash = tokenHash
	return nil
}

// TokenHashCompare compares the existing token hash with the provided token.
func (u *UserSession) TokenHashCompare(token string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.TokenHash), []byte(token)) == nil
}

// SetExpireAt sets the expire at time to 2 weeks from now.
func (u *UserSession) ExpireAtCreate() {
	u.ExpireAt = null.TimeFrom(time.Now().UTC().Add(14 * 24 * time.Hour)) // 2 weeks
}

// Expired returns true if the session is expired.
func (u *UserSession) IsExpired() bool {
	return u.ExpireAt.Valid && u.ExpireAt.Time.Before(time.Now().UTC())
}
