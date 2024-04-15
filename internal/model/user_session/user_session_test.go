package user_session

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/koraygocmen/golang-boilerplate/pkg/date"
)

func TestJSON(t *testing.T) {
	src := UserSession{
		ID:        1,
		CreatedAt: time.Now().UTC(),
		UserID:    1,
		Token:     "token",
		TokenHash: "hash",
		ClientIP:  "127.0.0.1",
		Purpose:   PurposeSessionCreate,
	}

	marshalled, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("want: no error when marshalling; got: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(marshalled, &m); err != nil {
		t.Fatalf("want: no error when unmarshalling; got: %v", err)
	}

	if len(m) != 5 {
		t.Fatalf("want: 5 fields; got: %v", len(m))
	}

	var dst UserSession
	if err := json.Unmarshal(marshalled, &dst); err != nil {
		t.Fatalf("want: no error when unmarshalling; got: %v", err)
	}

	if src.ID != dst.ID {
		t.Fatalf("want: id to match; got: %v", dst.ID)
	}

	if src.CreatedAt.Format("2006-01-02") != dst.CreatedAt.Format("2006-01-02") {
		t.Fatalf("want: created at to match; got: %v", dst.CreatedAt)
	}

	if dst.DeletedAt.Valid {
		t.Fatalf("want: deleted nil; got: %v", dst.DeletedAt)
	}

	if src.UserID != dst.UserID {
		t.Fatalf("want: user id to match; got: %v", dst.UserID)
	}

	if src.Token != dst.Token {
		t.Fatalf("want: token to match; got: %v", dst.Token)
	}

	if dst.TokenHash != "" {
		t.Fatalf("want: token hash to be not nil; got: %v", dst.TokenHash)
	}

	if dst.ClientIP != "" {
		t.Fatalf("want: client ip to be nil; got: %v", dst.ClientIP)
	}

	if src.Purpose != dst.Purpose {
		t.Fatalf("want: purpose to match; got: %v", dst.Purpose)
	}
}

func TestTokenHashCreate(t *testing.T) {
	userSession := UserSession{
		ID:        1,
		CreatedAt: time.Now().UTC(),
		UserID:    1,
		ClientIP:  "127.0.0.1",
	}

	// Test successful token hash creation.
	if err := userSession.TokenHashCreate(); err != nil {
		t.Fatalf("want: no error when creating token hash; got: %v", err)
	}

	if userSession.TokenHash == "" {
		t.Fatalf("want: token hash to be set; got: %v", userSession.TokenHash)
	}
}

func TestTokenHashCompare(t *testing.T) {
	userSession := UserSession{
		ID:        1,
		CreatedAt: time.Now().UTC(),
		UserID:    1,
		ClientIP:  "127.0.0.1",
	}

	// Test successful token hash creation.
	if err := userSession.TokenHashCreate(); err != nil {
		t.Fatalf("want: no error when creating token hash; got: %v", err)
	}

	token := userSession.Token
	if !userSession.TokenHashCompare(token) {
		t.Fatalf("want: success when comparing token hash; got: failure")
	}

	if userSession.TokenHashCompare("wrong") {
		t.Fatalf("want: failure when comparing token hash with wrong token; got: success")
	}

	// Test token hash creation with empty token.
	userSession.Token = ""
	if userSession.TokenHashCompare("") {
		t.Fatalf("want: failure when comparing token hash with nil token; got: success")
	}

	userSession.TokenHash = ""
	if userSession.TokenHashCompare("") {
		t.Fatalf("want: failure when comparing token hash with nil token and token hash; got: success")
	}
}

func TestExpireAtCreate(t *testing.T) {
	userSession := UserSession{
		ID:        1,
		CreatedAt: time.Now().UTC(),
		UserID:    1,
		ClientIP:  "",
	}

	// Test successful expire at set.
	userSession.ExpireAtCreate()

	if !userSession.ExpireAt.Valid {
		t.Fatalf("want: expire at to be valid; got: %v", userSession.ExpireAt)
	}

	if userSession.ExpireAt.Time.IsZero() {
		t.Fatalf("want: expire at time to be set; got: %v", userSession.ExpireAt.Time)
	}

	expected := userSession.CreatedAt.Add(14 * 24 * time.Hour)
	if !date.Equal(userSession.ExpireAt.Time, expected) {
		t.Fatalf("want: expire at time to be 14 days in the future; got: %v", userSession.ExpireAt.Time)
	}
}

func TestIsExpired(t *testing.T) {
	userSession := UserSession{
		ID:        1,
		CreatedAt: time.Now().UTC(),
		UserID:    1,
		ClientIP:  "",
	}

	if userSession.IsExpired() {
		t.Fatalf("want: not expired; got: expired")
	}

	// Test successful expire at set.
	userSession.ExpireAtCreate()

	if userSession.IsExpired() {
		t.Fatalf("want: not expired; got: expired")
	}

	userSession.ExpireAt.Time = time.Now().UTC().Add(-1 * time.Hour)
	if !userSession.IsExpired() {
		t.Fatalf("want: expired; got: not expired")
	}
}
