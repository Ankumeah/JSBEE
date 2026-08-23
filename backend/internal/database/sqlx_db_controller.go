package database

import (
	"github.com/jmoiron/sqlx"

	"context"
	"database/sql"
	"errors"
)

type SqlxDBController struct{ db *sqlx.DB }

func GetSqlxDBController(db *sqlx.DB) DBController {
	return &SqlxDBController{db}
}

func (s *SqlxDBController) DB() *sqlx.DB {
	return s.db
}

func (s *SqlxDBController) AddUser(
	ctx context.Context,
	user User,
) error {
	query := s.db.Rebind(`
    INSERT INTO users (name, email)
    VALUES (?, ?);
  `)

	_, err := s.db.ExecContext(ctx, query, user.Name, user.Email)

	if isUniqueViolation(err) {
		return ErrExistUser
	}

	return err
}

func (s *SqlxDBController) DeleteUser(
	ctx context.Context,
	email string,
) error {
	query := s.db.Rebind(`
    DELETE FROM users
    WHERE (email = ?);
  `)

	res, err := s.db.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return ErrInvalidUser
	}

	return nil
}

func (s *SqlxDBController) UpdateUser(
	ctx context.Context,
	email string,
	newUser User,
) error {
	query := s.db.Rebind(`
    UPDATE users
    SET name = ?, email = ?, role = ?
    WHERE (email = ?);
  `)

	res, err := s.db.ExecContext(
		ctx, query, newUser.Name, newUser.Email, newUser.Role, email,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return ErrInvalidUser
	}

	return nil
}

func (s *SqlxDBController) GetUser(
	ctx context.Context,
	email string,
) (User, error) {
	query := s.db.Rebind(`
    SELECT name, email, role
    FROM users
    WHERE email = ?;
  `)

	var user User
	err := s.db.QueryRowContext(ctx, query, email).Scan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrInvalidUser
	}

	return user, err
}
