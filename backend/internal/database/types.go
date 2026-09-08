package database

import (
	"github.com/Ankumeah/JSBEE/backend/internal/roles"

	"github.com/jmoiron/sqlx"

	"context"
	"uuid"
)

// This interface is reposonsible for executing all
// DB queries and returning type safe results
type DBController interface {
	// This exposes the raw underlying DB object.
	// This is only to be used to execute migrations
	DB() *sqlx.DB

	// Add a new user
	//
	// May return the following errors:
	//   - `database.ErrExistUser`
	//   - Errors by the underlying DB
	AddUser(ctx context.Context, user User) error

	// Delete a user
	//
	// May return the following errors:
	//   - `database.ErrInvalidUser`
	//   - Errors by the underlying DB
	DeleteUser(ctx context.Context, uuid uuid.UUID) error

	// Sets a user's subscription
	//
	// May return the following errors:
	//   - `database.ErrInvalidUser`
	//   - Errors by the underlying DB
	SetSubscription(ctx context.Context, subscribed bool, uuid uuid.UUID) error

	// Get the details of a user
	//
	// May return the following errors:
	//   - `database.ErrInvalidUser`
	//   - Errors by the underlying DB
	GetUser(ctx context.Context, uuid uuid.UUID) (User, error)

	// Get the details of a user by email
	//
	// May return the following errors:
	//   - `database.ErrInvalidUser`
	//   - Errors by the underlying DB
	GetUserByEmail(ctx context.Context, email string) (User, error)

	// Adds a new unapproved paper
	//
	// May return the following errors:
	//   - `database.ErrInvalidUser`
	//   - `database.ErrExistPaper`
	//   - Errors by the underlying DB
	AddPaper(ctx context.Context, paper Paper) error

	// Gets the details of a paper
	//
	// If the owner of a paper has been deleted
	// `paper.OwnerUUID` will be `nil`
	//
	// May return the following errors:
	//   - `database.ErrInvalidPaper`
	//   - Errors by the underlying DB
	GetPaper(ctx context.Context, uuid uuid.UUID) (Paper, error)

	// Get all papers by a user
	//
	// May return the following errors:
	//   - `database.ErrInvalidUser`
	//   - Errors by the underlying DB
	GetUserPapers(ctx context.Context, userUUID uuid.UUID) ([]Paper, error)

	// Updates a paper to max + 1 number and max volume and issue.
	// This is safe to run on an already approved paper
	//
	// May return the following errors:
	//   - `database.ErrInvalidPaper`
	//   - Errors by the underlying DB
	ApprovePaper(ctx context.Context, uuid uuid.UUID) error

	// Deletes a paper if it is unapproved.
	// In case the paper trying to be deleted is approved,
	// `database.ErrInvalidPaper` is returned
	//
	// May return the following errors:
	//   - `database.ErrInvalidPaper`
	//   - Errors by the underlying DB
	RejectPaper(ctx context.Context, paperUUID uuid.UUID) error

	// Get all volumes
	//
	// In case any paper's author has been deleted the
	// `Paper.OwnerUUID` feild will be `nil`
	//
	// May return the following errors:
	//   - Errors by the underlying DB
	GetVolumes(ctx context.Context) ([]Volume, error)

	// Get all unapproved papers
	//
	// May return the following errors:
	//   - Errors by the underlying DB
	GetUnapprovedPapers(ctx context.Context) ([]Paper, error)

	// Increment state.volume by 1
	//
	// May return the following errors:
	//   - Errors by the underlying DB
	IncrementVolume(ctx context.Context) error

	// Increment state.issue by 1
	//
	// May return the following errors:
	//   - Errors by the underlying DB
	IncrementIssue(ctx context.Context) error
}

type User struct {
	UUID       uuid.UUID  `json:"uuid" binding:"required" db:"uuid"`
	Name       string     `json:"name" binding:"required" db:"name"`
	Email      string     `json:"email" binding:"required" db:"email"`
	Role       roles.Role `json:"role" binding:"required" db:"role"`
	Subscribed bool       `json:"subscribed" binding:"required" db:"subscribed"`
}

type Paper struct {
	UUID      uuid.UUID  `json:"uuid" binding:"required" db:"uuid"`
	Title     string     `json:"title" binding:"required" db:"title"`
	Number    *uint64    `json:"number" binding:"required" db:"number"`
	Filename  string     `json:"filename" binding:"required" db:"filename"`
	OwnerUUID *uuid.UUID `json:"owner_uuid" binding:"required" db:"owner_uuid"`
}

type Issue struct {
	Number uint64  `json:"number" binding:"required"`
	Papers []Paper `json:"papers" binding:"required"`
}

type Volume struct {
	Number uint64  `json:"number" binding:"required"`
	Issues []Issue `json:"issues" binding:"required"`
}
