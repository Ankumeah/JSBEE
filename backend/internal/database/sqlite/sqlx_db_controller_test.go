package sqlite_test

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/database/migrations"
	"github.com/Ankumeah/JSBEE/backend/internal/database/sqlite"
	"github.com/Ankumeah/JSBEE/backend/internal/roles"

	"context"
	"errors"
	"slices"
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

		gotByEmail, err := sqlxDB.GetUserByEmail(Ctx, user.Email)
		if err != nil {
			t.Fatalf("Error while getting user by email: %v\n", err.Error())
		}
		if gotByEmail.UUID != user.UUID {
			t.Fatal("User looked up by email did not match")
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

	// Approving moves the paper to the waiting area but does NOT publish it
	volumes, err = sqlxDB.GetVolumes(Ctx)
	if err != nil {
		t.Fatalf("Error while getting volumes: %v\n", err.Error())
	}
	if len(volumes) != 0 {
		t.Fatal("Reviewed paper was published before PublishPapers was called")
	}

	reviewed, err := sqlxDB.GetReviewedPapers(Ctx)
	if err != nil {
		t.Fatalf("Error while getting reviewed papers: %v\n", err.Error())
	}
	if len(reviewed) != 1 {
		t.Fatalf("Expected 1 reviewed paper, got %v\n", len(reviewed))
	}
	if !reviewed[0].Reviewed {
		t.Fatal("Reviewed paper has reviewed=false")
	}

	count, err := sqlxDB.PublishPapers(Ctx)
	if err != nil {
		t.Fatalf("Error while publishing papers: %v\n", err.Error())
	}
	if count != 1 {
		t.Fatalf("Expected to publish 1 paper, got %v\n", count)
	}

	volumes, err = sqlxDB.GetVolumes(Ctx)
	if err != nil {
		t.Fatalf("Error while getting volumes: %v\n", err.Error())
	}
	t.Logf("After publishing\nExpected volumes: %v\nGot volumes     : %v",
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

	if gotPaper.UUID != paper.UUID ||
		gotPaper.Title != paper.Title ||
		gotPaper.Filename != paper.Filename ||
		gotPaper.OwnerUUID != paper.OwnerUUID {
		t.Log("Papers dont match")
	}
}

func TestPublishPapersSingleQuery(t *testing.T) {
	db, err := database.GetDBConnection(Ctx, "file::memory:", sqlite.DriverName)
	if err != nil {
		t.Fatalf("Error while getting DB connection: %v\n", err.Error())
	}
	defer db.Close()

	if err := migrations.ApplyMigrations(Ctx, db); err != nil {
		t.Fatalf("Error while applying migrations: %v\n", err.Error())
	}

	sqlxDB := sqlite.GetSqlxDBController(db)

	user := database.User{
		UUID:       uuid.New(),
		Name:       "test",
		Email:      "test@test.test",
		Role:       roles.Reviewer,
		Subscribed: true,
	}
	if err := sqlxDB.AddUser(Ctx, user); err != nil {
		t.Fatalf("Error while adding user: %v\n", err.Error())
	}

	// Add 3 papers, approve all so they land in the waiting area
	for i := 0; i < 3; i++ {
		p := database.Paper{
			UUID:      uuid.New(),
			Title:     "Test" + string(rune('A'+i)),
			Number:    nil,
			Filename:  "test" + string(rune('0'+i)) + ".pdf",
			OwnerUUID: &user.UUID,
		}
		if err := sqlxDB.AddPaper(Ctx, p); err != nil {
			t.Fatalf("Error while adding paper: %v\n", err.Error())
		}
		if err := sqlxDB.ApprovePaper(Ctx, p.UUID); err != nil {
			t.Fatalf("Error while approving paper: %v\n", err.Error())
		}
	}

	// Publish the 3 waiting papers in ONE query, numbers should be 1,2,3
	count, err := sqlxDB.PublishPapers(Ctx)
	if err != nil {
		t.Fatalf("Error while publishing papers: %v\n", err.Error())
	}
	if count != 3 {
		t.Fatalf("Expected to publish 3 papers, got %v\n", count)
	}

	// Add a 4th paper, approve it, publish again — it should get number 4
	extra := database.Paper{
		UUID:      uuid.New(),
		Title:     "Extra",
		Number:    nil,
		Filename:  "extra.pdf",
		OwnerUUID: &user.UUID,
	}
	if err := sqlxDB.AddPaper(Ctx, extra); err != nil {
		t.Fatalf("Error while adding extra paper: %v\n", err.Error())
	}
	if err := sqlxDB.ApprovePaper(Ctx, extra.UUID); err != nil {
		t.Fatalf("Error while approving extra paper: %v\n", err.Error())
	}
	count, err = sqlxDB.PublishPapers(Ctx)
	if err != nil {
		t.Fatalf("Error while publishing extra paper: %v\n", err.Error())
	}
	if count != 1 {
		t.Fatalf("Expected to publish 1 extra paper, got %v\n", count)
	}

	// Publishing again should publish 0 (waiting area is empty)
	count, err = sqlxDB.PublishPapers(Ctx)
	if err != nil {
		t.Fatalf("Error while publishing papers again: %v\n", err.Error())
	}
	if count != 0 {
		t.Fatalf("Expected 0 second publish, got %v\n", count)
	}

	papers, err := sqlxDB.GetUserPapers(Ctx, user.UUID)
	if err != nil {
		t.Fatalf("Error while getting user papers: %v\n", err.Error())
	}
	var numbers []uint64
	for _, p := range papers {
		if p.Number != nil {
			numbers = append(numbers, *p.Number)
		} else {
			t.Fatal("Expected all papers to have numbers")
		}
	}

	// Papers sorted? GetUserPapers has no ORDER BY so sort manually
	slices.Sort(numbers)
	expected := []uint64{1, 2, 3, 4}
	if !slices.Equal(numbers, expected) {
		t.Fatalf("Expected numbers %v, got %v\n", expected, numbers)
	}
}

func TestBlog(t *testing.T) {
	db, err := database.GetDBConnection(Ctx, "file::memory:", sqlite.DriverName)
	if err != nil {
		t.Fatalf("Error while getting DB connection: %v\n", err.Error())
	}
	defer db.Close()

	if err := migrations.ApplyMigrations(Ctx, db); err != nil {
		t.Fatalf("Error while applying migrations: %v\n", err.Error())
	}

	sqlxDB := sqlite.GetSqlxDBController(db)

	blog := database.Blog{
		UUID:      uuid.New(),
		Title:     "First post",
		Filename:  "blog-123.md",
		CreatedAt: 123,
		UpdatedAt: 123,
	}

	if err := sqlxDB.AddBlog(Ctx, blog); err != nil {
		t.Fatalf("Error while adding blog: %v\n", err.Error())
	}

	if err := sqlxDB.AddBlog(Ctx, blog); !errors.Is(err, database.ErrExistBlog) {
		t.Fatalf("Expected ErrExistBlog on duplicate, got %v\n", err)
	}

	gotBlog, err := sqlxDB.GetBlog(Ctx, blog.UUID)
	if err != nil {
		t.Fatalf("Error while getting blog: %v\n", err.Error())
	}
	if gotBlog != blog {
		t.Fatalf("Expected blog %v, got %v\n", blog, gotBlog)
	}

	blogs, err := sqlxDB.GetBlogs(Ctx)
	if err != nil {
		t.Fatalf("Error while getting blogs: %v\n", err.Error())
	}
	if len(blogs) != 1 || blogs[0] != blog {
		t.Fatalf("Expected 1 blog %v, got %v\n", blog, blogs)
	}

	blog.Title = "Edited post"
	blog.UpdatedAt = 456
	if err := sqlxDB.UpdateBlog(Ctx, blog); err != nil {
		t.Fatalf("Error while updating blog: %v\n", err.Error())
	}
	gotBlog, err = sqlxDB.GetBlog(Ctx, blog.UUID)
	if err != nil {
		t.Fatalf("Error while getting updated blog: %v\n", err.Error())
	}
	if gotBlog != blog {
		t.Fatalf("Expected updated blog %v, got %v\n", blog, gotBlog)
	}

	if err := sqlxDB.DeleteBlog(Ctx, blog.UUID); err != nil {
		t.Fatalf("Error while deleting blog: %v\n", err.Error())
	}
	if _, err := sqlxDB.GetBlog(
		Ctx, blog.UUID,
	); !errors.Is(err, database.ErrInvalidBlog) {
		t.Fatalf("Expected ErrInvalidBlog after delete, got %v\n", err)
	}
}
