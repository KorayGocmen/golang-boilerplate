package user

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/koraygocmen/null"
)

func TestJSON(t *testing.T) {
	src := &User{
		ID:            1,
		CreatedAt:     time.Now().UTC(),
		Email:         null.StringFrom("koray@test.com"),
		EmailVerified: null.BoolFrom(true),
		Password:      null.StringFrom("123456"),
		PasswordHash:  null.StringFrom("hash"),
		GivenNames:    null.StringFrom("Koray"),
		Surname:       null.StringFrom("Gocmen"),
	}

	marshalled, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("want: no error when marshalling; got: %v", err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(marshalled, &m); err != nil {
		t.Fatalf("want: no error when unmarshalling; got: %v", err)
	}

	if len(m) != 7 {
		t.Fatalf("want: 7 fields; got: %v", len(m))
	}

	var dst User
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

	if src.Email.String != dst.Email.String {
		t.Fatalf("want: email to match; got: %v", dst.Email)
	}

	if src.EmailVerified.Bool != dst.EmailVerified.Bool {
		t.Fatalf("want: email verified to match; got: %v", dst.EmailVerified)
	}

	if dst.Password.Valid {
		t.Fatalf("want: password nil; got: %v", dst.Password)
	}

	if dst.PasswordHash.Valid {
		t.Fatalf("want: password hash nil; got: %v", dst.PasswordHash)
	}

	if src.GivenNames.String != dst.GivenNames.String {
		t.Fatalf("want: given names to match; got: %v", dst.GivenNames)
	}

	if src.Surname.String != dst.Surname.String {
		t.Fatalf("want: surname to match; got: %v", dst.Surname)
	}
}

func TestPasswordHashCreate(t *testing.T) {
	user := &User{}
	if err := user.PasswordHashCreate(); err == nil {
		t.Fatalf("want: error when password is nil; got: err nil")
	}

	user.Password = null.StringFrom("password")
	if !user.Password.Valid {
		t.Fatalf("want: password valid; got: password invalid")
	}

	if err := user.PasswordHashCreate(); err != nil {
		t.Fatalf("want: no error when password hash create; got: %v", err)
	}
}

func TestPasswordHashCompare(t *testing.T) {
	user := &User{}
	if user.PasswordHashCompare("wrong") {
		t.Fatalf("want: invalid password hash compare; got: valid")
	}

	user.Password = null.StringFrom("password")
	if err := user.PasswordHashCreate(); err != nil {
		t.Fatalf("want: no error when password hash create; got: %v", err)
	}

	if user.PasswordHashCompare("wrong") {
		t.Fatalf("want: invalid password hash compare; got: valid")
	}

	if !user.PasswordHashCompare("password") {
		t.Fatalf("want: valid password hash compare; got: invalid")
	}
}
