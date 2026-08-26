package database

import (
	"github.com/Ankumeah/JSBEE/backend/internal/roles"

	"github.com/jmoiron/sqlx"

	"context"
)

type DBController interface {
	DB() *sqlx.DB
	AddUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, email string) error
	UpdateUser(ctx context.Context, email string, newUser User) error
	GetUser(ctx context.Context, name string) (User, error)

	GetVolumes(ctx context.Context) ([]Volume, error)
}

type User struct {
	Name  string     `json:"name" binding:"required"`
	Email string     `json:"email" binding:"required"`
	Role  roles.Role `json:"role" binding:"required"`
}

type Paper struct {
	Title    string  `json:"title" binding:"required"`
	Number   *uint64 `json:"number" binding:"required"`
	Filename string  `json:"filename" binding:"required"`
	Owner    string  `json:"owner" binding:"required"`
}

type Issue struct {
	Number uint64  `json:"number" binding:"required"`
	Papers []Paper `json:"papers" binding:"required"`
}

type Volume struct {
	Number uint64  `json:"number" binding:"required"`
	Issues []Issue `json:"issues" binding:"required"`
}
