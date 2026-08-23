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
	GetUser(ctx context.Context, email string) (User, error)
}

type User struct {
	Name  string     `json:"name,required"`
	Email string     `json:"email,required"`
	Role  roles.Role `json:"role,required"`
}
