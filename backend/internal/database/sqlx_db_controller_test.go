package database_test

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/database/migrations"
	"github.com/Ankumeah/JSBEE/backend/internal/roles"

	"context"
	"errors"
	"testing"
	"uuid"
)

var Ctx = context.Background()

func TestMigrations(t *testing.T) {
	dbURL := "file::memory:?cache=shared"
	db, err := database.GetDBConnection(Ctx, dbURL)
	if err != nil {
		t.Fatalf("Error while getting DB connection: %v\n", err.Error())
	}
	defer db.Close()

	if err := migrations.ApplyMigrations(Ctx, db); err != nil {
		t.Fatalf("Error while applying migrations: %v\n", err.Error())
	}
}

// This is a minimal test that just checks the
// happy path. It only proves that user related
// functions work, not that they work well
//
// TODO: Check non happy paths
func TestUser(t *testing.T) {
	cases := []database.User{
		database.User{
			UUID:       uuid.New(),
			Name:       "test",
			Email:      "test@test.test",
			Role:       roles.Viewer,
			Subscribed: true,
		},
	}

	for _, user := range cases {
		dbURL := "file::memory:?cache=shared"
		db, err := database.GetDBConnection(Ctx, dbURL)
		if err != nil {
			t.Fatalf("Error while getting DB connection: %v\n", err.Error())
		}
		defer db.Close()

		if err := migrations.ApplyMigrations(Ctx, db); err != nil {
			t.Fatalf("Error while applying migrations: %v\n", err.Error())
		}

		sqlxDB := database.GetSqlxDBController(db)

		if err := sqlxDB.AddUser(Ctx, user); err != nil {
			t.Fatalf("Error while adding user: %v\n", err.Error())
		}

		gotUser, err := sqlxDB.GetUser(Ctx, user.UUID)
		if err != nil {
			t.Fatalf("Error while getting user: %v\n", err.Error())
		}

		t.Logf(`
      Created user: %v
      Got user    : %v `,
			user, gotUser,
		)

		if !(user.UUID == gotUser.UUID) ||
			!(user.Name == gotUser.Name) ||
			!(user.Email == gotUser.Email) ||
			!(user.Role == gotUser.Role) ||
			!(user.Subscribed == gotUser.Subscribed) {
			t.FailNow()
		}

		if err := sqlxDB.DeleteUser(Ctx, user.UUID); err != nil {
			t.Fatalf("Error while deleting user: %v\n", err.Error())
		}

		if _, err := sqlxDB.GetUser(
			Ctx, user.UUID,
		); !errors.Is(err, database.ErrInvalidUser) && err != nil {
			t.Fatalf("Error while getting user: %v\n", err.Error())
		} else if !errors.Is(err, database.ErrInvalidUser) {
			t.Fatal("User wasnt deleted")
		}
	}
}
