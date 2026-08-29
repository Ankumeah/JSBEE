package sqlite_test

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/database/migrations"
	"github.com/Ankumeah/JSBEE/backend/internal/database/sqlite"
	"github.com/Ankumeah/JSBEE/backend/internal/roles"

	"context"
	"errors"
	"testing"
	"uuid"
)

var Ctx = context.Background()

func TestMigrations(t *testing.T) {
	db, err := database.GetDBConnection(Ctx, "file::memory:", sqlite.DriverName)
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
		db, err := database.GetDBConnection(
			Ctx,
			"file::memory:",
			sqlite.DriverName,
		)
		if err != nil {
			t.Fatalf("Error while getting DB connection: %v\n", err.Error())
		}
		defer db.Close()

		if err := migrations.ApplyMigrations(Ctx, db); err != nil {
			t.Fatalf("Error while applying migrations: %v\n", err.Error())
		}

		sqlxDB := sqlite.GetSqlxDBController(db)

		if err := sqlxDB.AddUser(Ctx, user); err != nil {
			t.Fatalf("Error while adding user: %v\n", err.Error())
		}

		gotUser, err := sqlxDB.GetUser(Ctx, user.UUID)
		if err != nil {
			t.Fatalf("Error while getting user: %v\n", err.Error())
		}

		t.Logf("\nCreated user: %v\nGot user    : %v",
			user, gotUser,
		)

		if user.UUID != gotUser.UUID ||
			user.Name != gotUser.Name ||
			user.Email != gotUser.Email ||
			user.Role != gotUser.Role ||
			user.Subscribed != gotUser.Subscribed {
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

func TestPapers(t *testing.T) {
	paper := database.Paper{
		UUID:      uuid.New(),
		Title:     "Test",
		Number:    nil,
		Filename:  "test.pdf",
		OwnerUUID: nil,
	}
	expectedResult := database.Volume{
		Number: 1,
		Issues: []database.Issue{
			database.Issue{
				Number: 1,
				Papers: []database.Paper{paper},
			},
		},
	}

	db, err := database.GetDBConnection(Ctx, "file::memory:", sqlite.DriverName)
	if err != nil {
		t.Fatalf("Error while getting DB connection: %v\n", err.Error())
	}
	defer db.Close()

	if err := migrations.ApplyMigrations(Ctx, db); err != nil {
		t.Fatalf("Error while applying migrations: %v\n", err.Error())
	}

	sqlxDB := sqlite.GetSqlxDBController(db)

	testUser := database.User{
		UUID:       uuid.New(),
		Name:       "test",
		Email:      "test@test.test",
		Role:       roles.Viewer,
		Subscribed: true,
	}
	paper.OwnerUUID = &testUser.UUID

	if err := sqlxDB.AddUser(Ctx, testUser); err != nil {
		t.Fatalf("Error while adding user: %v\n", err.Error())
	}

	if err := sqlxDB.AddPaper(Ctx, paper); err != nil {
		t.Fatalf("Error while adding paper: %v\n", err.Error())
	}

	volumes, err := sqlxDB.GetVolumes(Ctx)
	if err != nil {
		t.Fatalf("Error while getting volumes: %v\n", err.Error())
	}
	t.Logf("Before approval\nExpected volumes: %v\nGot volumes     : %v",
		expectedResult, volumes,
	)

	if len(volumes) != 0 {
		t.Fatal("Got unapproved papers")
	}

	if err := sqlxDB.ApprovePaper(Ctx, paper.UUID); err != nil {
		t.Fatalf("Error while approveing paper: %v\n", err.Error())
	}

	volumes, err = sqlxDB.GetVolumes(Ctx)
	if err != nil {
		t.Fatalf("Error while getting volumes: %v\n", err.Error())
	}
	t.Logf("After approval\nExpected volumes: %v\nGot volumes     : %v",
		expectedResult, volumes,
	)

	if len(volumes) != 1 {
		t.Fatal("Incorrect number of volumes")
	} else if len(volumes[0].Issues) != 1 {
		t.Fatal("Incorrect number of issues")
	} else if len(volumes[0].Issues[0].Papers) != 1 {
		t.Fatal("Incorrect number of papers")
	}

	volume := volumes[0]
	issue := volume.Issues[0]
	gotPaper := issue.Papers[0]

	if volume.Number != expectedResult.Number ||
		issue.Number != expectedResult.Issues[0].Number ||
		gotPaper.Number != expectedResult.Issues[0].Papers[0].Number {
		t.Log("Incorrent numbering")
	}

	if gotPaper.UUID != paper.UUID ||
		gotPaper.Title != paper.Title ||
		gotPaper.Filename != paper.Filename ||
		gotPaper.OwnerUUID != paper.OwnerUUID {
		t.Log("Papers dont match")
	}
}
