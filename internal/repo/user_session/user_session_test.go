package user_session

import (
	"os"
	"testing"
	"time"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/database/databasetest"
	User "github.com/koraygocmen/golang-boilerplate/internal/model/user"
	UserSession "github.com/koraygocmen/golang-boilerplate/internal/model/user_session"
	UserRepo "github.com/koraygocmen/golang-boilerplate/internal/repo/user"
	"github.com/koraygocmen/null"
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

func populate() ([]*User.User, []*UserSession.UserSession, error) {
	userRepo := UserRepo.New(dbTest.DB.GORM)
	userSessionRepo := New(dbTest.DB.GORM)

	users := []*User.User{
		{
			Email:         null.StringFrom("koray@test.com"),
			EmailVerified: null.BoolFrom(true),
			Password:      null.StringFrom("123456"),
			PasswordHash:  null.StringFrom("hash1"),
			GivenNames:    null.StringFrom("KORAY"),
			Surname:       null.StringFrom("GOCMEN"),
		},
		{
			Email:         null.StringFrom("taylan@test.com"),
			EmailVerified: null.BoolFrom(true),
			Password:      null.StringFrom("123456"),
			PasswordHash:  null.StringFrom("hash2"),
			GivenNames:    null.StringFrom("TAYLAN"),
			Surname:       null.StringFrom("GOCMEN"),
		},
	}
	for _, u := range users {
		if err := userRepo.Create(context.Background(), u); err != nil {
			return nil, nil, err
		}
	}

	userSessions := []*UserSession.UserSession{
		{
			UserID:    users[0].ID,
			Token:     "token1",
			TokenHash: "tokenhash1",
			ClientIP:  "127.0.0.1",
			Purpose:   UserSession.PurposeSessionCreate,
			ExpireAt:  null.TimeFrom(time.Now().UTC().Add(-time.Hour)),
		},
		{
			UserID:    users[0].ID,
			Token:     "token2",
			TokenHash: "tokenhash2",
			ClientIP:  "127.0.0.1",
			Purpose:   UserSession.PurposePasswordReset,
			ExpireAt:  null.TimeFrom(time.Now().UTC().Add(time.Hour)),
		},
		{
			UserID:    users[1].ID,
			Token:     "token3",
			TokenHash: "tokenhash3",
			ClientIP:  "127.0.0.1",
			Purpose:   UserSession.PurposeSessionCreate,
			ExpireAt:  null.TimeFrom(time.Now().UTC().Add(-time.Hour)),
		},
		{
			UserID:    users[1].ID,
			Token:     "token4",
			TokenHash: "tokenhash4",
			ClientIP:  "127.0.0.1",
			Purpose:   UserSession.PurposePasswordReset,
			ExpireAt:  null.TimeFrom(time.Now().UTC().Add(time.Hour)),
		},
	}
	for _, userSession := range userSessions {
		if err := userSessionRepo.Create(context.Background(), userSession); err != nil {
			return nil, nil, err
		}
	}

	return users, userSessions, nil
}

func TestCreate(t *testing.T) {
	dbClean()

	_, userSessions, err := populate()
	if err != nil {
		t.Fatalf("populate error: %v", err)
	}

	for _, u := range userSessions {
		var uGot UserSession.UserSession
		err = dbTest.DB.GORM.
			Raw(`SELECT * FROM public.user_session WHERE id = ?`, u.ID).
			Scan(&uGot).
			Error
		if err != nil {
			t.Fatalf("select error: %v", err)
		}

		if u.ID != uGot.ID {
			t.Fatalf("want: id match; got: id mismatch: %d != %d", u.ID, uGot.ID)
		}

		if uGot.CreatedAt.IsZero() {
			t.Fatalf("want: created at to be set; got: created at zero")
		}

		if u.UserID != uGot.UserID {
			t.Fatalf("want: user id match; got: user id mismatch: %d != %d", u.UserID, uGot.UserID)
		}

		if uGot.Token != "" {
			t.Fatalf("want: token nil; got: token not nill: %s", uGot.Token)
		}

		if u.TokenHash != uGot.TokenHash {
			t.Fatalf("want: token hash match; got: token hash mismatch: %s != %s", u.TokenHash, uGot.TokenHash)
		}

		if u.ClientIP != uGot.ClientIP {
			t.Fatalf("want: client ip match; got: client ip mismatch: %s != %s", u.ClientIP, uGot.ClientIP)
		}

		if u.Purpose != uGot.Purpose {
			t.Fatalf("want: purpose match; got: purpose mismatch: %s != %s", u.Purpose, uGot.Purpose)
		}
	}
}

func TestGetByID(t *testing.T) {
	dbClean()

	_, userSessions, err := populate()
	if err != nil {
		t.Fatalf("populate error: %v", err)
	}

	userSessionRepo := New(dbTest.DB.GORM)

	for _, u := range userSessions {
		uGot, err := userSessionRepo.GetByID(context.Background(), u.ID)
		if err != nil {
			t.Fatalf("want: get by id error nil; got: %v", err)
		}

		if uGot == nil {
			t.Fatalf("want: user session get by id successful not nil; got: nil")
		}

		if u.ID != uGot.ID {
			t.Fatalf("want: id match; got: id mismatch: %d != %d", u.ID, uGot.ID)
		}

		if uGot.CreatedAt.IsZero() {
			t.Fatalf("want: created at to be set; got: created at zero")
		}

		if u.UserID != uGot.UserID {
			t.Fatalf("want: user id match; got: user id mismatch: %d != %d", u.UserID, uGot.UserID)
		}

		if uGot.Token != "" {
			t.Fatalf("want: token nil; got: token not nill: %s", uGot.Token)
		}

		if u.TokenHash != uGot.TokenHash {
			t.Fatalf("want: token hash match; got: token hash mismatch: %s != %s", u.TokenHash, uGot.TokenHash)
		}

		if u.ClientIP != uGot.ClientIP {
			t.Fatalf("want: client ip match; got: client ip mismatch: %s != %s", u.ClientIP, uGot.ClientIP)
		}

		if u.Purpose != uGot.Purpose {
			t.Fatalf("want: purpose match; got: purpose mismatch: %s != %s", u.Purpose, uGot.Purpose)
		}
	}
}

func TestListActive(t *testing.T) {
	dbClean()

	users, userSessions, err := populate()
	if err != nil {
		t.Fatalf("populate error: %v", err)
	}

	userSessionRepo := New(dbTest.DB.GORM)

	cases := []struct {
		userID int64
		want   []*UserSession.UserSession
	}{
		{
			userID: users[0].ID,
			want: []*UserSession.UserSession{
				userSessions[1],
			},
		},
		{
			userID: users[1].ID,
			want: []*UserSession.UserSession{
				userSessions[3],
			},
		},
		{
			userID: 999,
			want:   []*UserSession.UserSession{},
		},
	}

	for _, c := range cases {
		uGot, err := userSessionRepo.ListActive(context.Background(), c.userID)
		if err != nil {
			t.Fatalf("want: list error nil; got: %v", err)
		}

		if len(uGot) != len(c.want) {
			t.Fatalf("want: user sessions list length %d; got: %d", len(c.want), len(uGot))
		}

		for i, u := range uGot {
			if u.ID != c.want[i].ID {
				t.Fatalf("want: id match; got: id mismatch: %d != %d", u.ID, c.want[i].ID)
			}
		}
	}
}
