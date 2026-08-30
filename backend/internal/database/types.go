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
	DB() *sqlx.DB
	AddUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, uuid uuid.UUID) error
	UpdateUser(ctx context.Context, uuid uuid.UUID, newUser User) error
	GetUser(ctx context.Context, uuid uuid.UUID) (User, error)

	AddPaper(ctx context.Context, paper Paper) error
	ApprovePaper(ctx context.Context, uuid uuid.UUID) error
	DeletePaper(ctx context.Context, uuid uuid.UUID) error

	GetVolumes(ctx context.Context) ([]Volume, error)
	GetUnapprovedPapers(ctx context.Context) ([]Paper, error)
}

type User struct {
	UUID       uuid.UUID  `json:"uuid" binding:"required"`
	Name       string     `json:"name" binding:"required"`
	Email      string     `json:"email" binding:"required"`
	Role       roles.Role `json:"role" binding:"required"`
	Subscribed bool       `json:"subscribed" binding:"required"`
}

type Paper struct {
	UUID      uuid.UUID  `json:"uuid" binding:"required"`
	Title     string     `json:"title" binding:"required"`
	Number    *uint64    `json:"number" binding:"required"`
	Filename  string     `json:"filename" binding:"required"`
	OwnerUUID *uuid.UUID `json:"owner" binding:"required"`
}

type Issue struct {
	Number uint64  `json:"number" binding:"required"`
	Papers []Paper `json:"papers" binding:"required"`
}

type Volume struct {
	Number uint64  `json:"number" binding:"required"`
	Issues []Issue `json:"issues" binding:"required"`
}
