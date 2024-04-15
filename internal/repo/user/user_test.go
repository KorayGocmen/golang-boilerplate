package user_repo

import (
	"os"
	"testing"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	databasetest "github.com/koraygocmen/golang-boilerplate/internal/database/databasetest"
	User "github.com/koraygocmen/golang-boilerplate/internal/model/user"
	"github.com/koraygocmen/null"
	_ "github.com/lib/pq"
)

var (
	dbTest = databasetest.Get()
)

func TestMain(m *testing.M) {
	code := m.Run()

	// Purge and exit.
	dbTest.Purge()
	os.Exit(code)
}

func dbClean() {
	ctx := context.Background()

	dbTest.DB.Reset(ctx)
	dbTest.DB.Up(ctx)
	dbTest.DB.Seed(ctx)
}

func TestCreate(t *testing.T) {
	dbClean()

	userRepo := New(dbTest.DB.GORM)

	// Create a user.
	u := &User.User{
		Email:         null.StringFrom("koray@test.com"),
		EmailVerified: null.BoolFrom(true),
		Password:      null.StringFrom("123456"),
		PasswordHash:  null.StringFrom("hash"),
		GivenNames:    null.StringFrom("Koray"),
		Surname:       null.StringFrom("Gocmen"),
	}
	if err := userRepo.Create(context.Background(), u); err != nil {
		t.Fatalf("want: create error nil; got: %v", err)
	}

	// Get the user.
	var uGot User.User
	err := dbTest.DB.GORM.
		Raw(`SELECT * FROM public.user WHERE id = ?`, u.ID).
		Scan(&uGot).
		Error
	if err != nil {
		t.Fatalf("select error: %v", err)
	}

	// Compare.
	if u.ID != uGot.ID {
		t.Fatal("want: id match; got: id does not match")
	}

	if uGot.CreatedAt.IsZero() {
		t.Fatalf("want: created at to be set; got: created at zero")
	}

	if uGot.DeletedAt.Valid {
		t.Fatal("want: deleted at not valid; got: deleted at valid")
	}

	if u.Email != uGot.Email {
		t.Fatal("want: email match; got: email does not match", u.Email, uGot.Email)
	}

	if u.EmailVerified != uGot.EmailVerified {
		t.Fatal("want: email verified match; got: email verified does not match", u.EmailVerified, uGot.EmailVerified)
	}

	if u.PasswordHash != uGot.PasswordHash {
		t.Fatal("want: password hash match; got: password hash does not match", u.PasswordHash, uGot.PasswordHash)
	}

	if u.GivenNames != uGot.GivenNames {
		t.Fatal("want: given names match; got: given names does not match", u.GivenNames, uGot.GivenNames)
	}

	if u.Surname != uGot.Surname {
		t.Fatal("want: surname match; got: surname does not match", u.Surname, uGot.Surname)
	}
}

func TestSave(t *testing.T) {
	dbClean()

	userRepo := New(dbTest.DB.GORM)

	// Create a user.
	u := &User.User{
		Email:         null.StringFrom("koray@test.com"),
		EmailVerified: null.BoolFrom(true),
		Password:      null.StringFrom("123456"),
		PasswordHash:  null.StringFrom("hash"),
		GivenNames:    null.StringFrom("Koray"),
		Surname:       null.StringFrom("Gocmen"),
	}
	if err := userRepo.Create(context.Background(), u); err != nil {
		t.Fatalf("want: create error nil; got: %v", err)
	}

	// Update the user.
	u.Email = null.StringFrom("koray2@test.com")
	u.EmailVerified = null.BoolFrom(false)
	u.Password = null.StringFrom("1234567")
	u.PasswordHash = null.StringFrom("hash2")
	u.GivenNames = null.StringFrom("Koray2")
	u.Surname = null.StringFrom("Gocmen2")
	if err := userRepo.Save(context.Background(), u); err != nil {
		t.Fatalf("want: update error nil; got: %v", err)
	}

	// Get the user.
	var uGot User.User
	err := dbTest.DB.GORM.
		Raw(`SELECT * FROM public.user WHERE id = ?`, u.ID).
		Scan(&uGot).
		Error
	if err != nil {
		t.Fatalf("select error: %v", err)
	}

	// Compare.
	if u.ID != uGot.ID {
		t.Fatal("want: id match; got: id does not match")
	}

	if uGot.CreatedAt.IsZero() {
		t.Fatalf("want: created at to be set; got: created at zero")
	}

	if uGot.DeletedAt.Valid {
		t.Fatal("want: deleted at not valid; got: deleted at valid")
	}

	if u.Email != uGot.Email {
		t.Fatal("want: email match; got: email does not match", u.Email, uGot.Email)
	}

	if u.EmailVerified != uGot.EmailVerified {
		t.Fatal("want: email verified match; got: email verified does not match", u.EmailVerified, uGot.EmailVerified)
	}

	if u.PasswordHash != uGot.PasswordHash {
		t.Fatal("want: password hash match; got: password hash does not match", u.PasswordHash, uGot.PasswordHash)
	}

	if u.GivenNames != uGot.GivenNames {
		t.Fatal("want: given names match; got: given names does not match", u.GivenNames, uGot.GivenNames)
	}

	if u.Surname != uGot.Surname {
		t.Fatal("want: surname match; got: surname does not match", u.Surname, uGot.Surname)
	}
}

func TestList(t *testing.T) {
	dbClean()

	userRepo := New(dbTest.DB.GORM)

	users := []*User.User{
		{
			Email:         null.StringFrom("test@test.com"),
			EmailVerified: null.BoolFrom(true),
			Password:      null.StringFrom("123456"),
			PasswordHash:  null.StringFrom("hash1"),
			GivenNames:    null.StringFrom("Koray"),
			Surname:       null.StringFrom("Gocmen"),
		},
		{
			Email:         null.StringFrom("test2@test.com"),
			EmailVerified: null.BoolFrom(true),
			Password:      null.StringFrom("123456"),
			PasswordHash:  null.StringFrom("hash2"),
			GivenNames:    null.StringFrom("Koray"),
			Surname:       null.StringFrom("Gocmen"),
		},
	}
	for _, u := range users {
		if err := userRepo.Create(context.Background(), u); err != nil {
			t.Fatalf("want: create error nil; got: %v", err)
		}
	}

	// Get the users.
	usersGot, err := userRepo.List(context.Background(), 10, 1)
	if err != nil {
		t.Fatalf("want: list error nil; got: %v", err)
	}

	if len(usersGot) != len(users) {
		t.Fatalf("want: users length match; got: users length does not match")
	}
}

func TestGetByID(t *testing.T) {
	dbClean()

	userRepo := New(dbTest.DB.GORM)

	// Create a user.
	u1 := &User.User{
		Email:         null.StringFrom("koray1@test.com"),
		EmailVerified: null.BoolFrom(true),
		Password:      null.StringFrom("123456"),
		PasswordHash:  null.StringFrom("hash1"),
		GivenNames:    null.StringFrom("Koray"),
		Surname:       null.StringFrom("Gocmen"),
	}
	u2 := &User.User{
		Email:         null.StringFrom("koray2@test.com"),
		EmailVerified: null.BoolFrom(true),
		Password:      null.StringFrom("123456"),
		PasswordHash:  null.StringFrom("hash2"),
		GivenNames:    null.StringFrom("Koray"),
		Surname:       null.StringFrom("Gocmen"),
	}
	for _, u := range []*User.User{u1, u2} {
		if err := userRepo.Create(context.Background(), u); err != nil {
			t.Fatalf("want: create error nil; got: %v", err)
		}

		// Get the user.
		uGot, err := userRepo.GetByID(context.Background(), u.ID)
		if err != nil {
			t.Fatalf("want: get by id error nil; got: %v", err)
		}

		// Compare.
		if u.ID != uGot.ID {
			t.Fatal("want: id match; got: id does not match")
		}
	}
}

func TestGetByEmail(t *testing.T) {
	dbClean()

	userRepo := New(dbTest.DB.GORM)
	users := []*User.User{
		{
			Email:         null.StringFrom("koray@test.com"),
			EmailVerified: null.BoolFrom(true),
			Password:      null.StringFrom("123456"),
			PasswordHash:  null.StringFrom("hash1"),
			GivenNames:    null.StringFrom("Koray"),
			Surname:       null.StringFrom("Gocmen"),
		},
		{
			Email:         null.StringFrom("koray2@test.com"),
			EmailVerified: null.BoolFrom(true),
			Password:      null.StringFrom("123456"),
			PasswordHash:  null.StringFrom("hash2"),
			GivenNames:    null.StringFrom("Koray"),
			Surname:       null.StringFrom("Gocmen"),
		},
		{
			Email:         null.StringFrom("koray3@test.com"),
			EmailVerified: null.BoolFrom(true),
			Password:      null.StringFrom("123456"),
			PasswordHash:  null.StringFrom("hash3"),
			GivenNames:    null.StringFrom("Koray"),
			Surname:       null.StringFrom("Gocmen"),
		},
	}
	for _, u := range users {
		if err := userRepo.Create(context.Background(), u); err != nil {
			t.Fatalf("want: create success; got: %v", err)
		}
	}

	for _, u := range users {
		// Get the user.
		uGot, err := userRepo.GetByEmail(
			context.Background(),
			u.Email.String,
		)
		if err != nil {
			t.Fatalf("select error: %v", err)
		}

		if uGot == nil {
			t.Fatalf("want: user get by phone number successful not nil; got: nil")
		}

		// Compare.
		if u.ID != uGot.ID {
			t.Fatal("want: id match; got: id does not match", u.ID, uGot.ID)
		}
	}
}
