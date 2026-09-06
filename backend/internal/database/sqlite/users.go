package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"context"
	"database/sql"
	"errors"
	"uuid"
)

func (s *SqlxDBController) AddUser(
	ctx context.Context,
	user database.User,
) error {
	query := s.db.Rebind(`
    INSERT INTO users (uuid, name, email, subscribed)
    VALUES (?, ?, ?, ?);
  `)

	_, err := s.db.ExecContext(
		ctx, query,
		user.UUID, user.Name, user.Email, user.Subscribed,
	)

	if isUniqueViolation(err) {
		return database.ErrExistUser
	}

	return err
}

func (s *SqlxDBController) DeleteUser(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	query := s.db.Rebind(`
    DELETE FROM users
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}
	// return `database.ErrInvalidUser` if no user was deleted
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return database.ErrInvalidUser
	}

	return nil
}

func (s *SqlxDBController) SetSubscription(
	ctx context.Context,
	subscribed bool,
	uuid uuid.UUID,
) error {
	query := s.db.Rebind(`
    UPDATE users
    SET subscribed = ?
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(
		ctx, query, subscribed, uuid,
	)
	if err != nil {
		return err
	}
	// return `database.ErrInvalidUser` if no user was updated
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return database.ErrInvalidUser
	}

	return nil
}

func (s *SqlxDBController) GetUser(
	ctx context.Context,
	uuid uuid.UUID,
) (database.User, error) {
	query := s.db.Rebind(`
    SELECT uuid, name, email, role, subscribed
    FROM users
    WHERE (uuid = ?);
  `)

	var user database.User
	err := s.db.QueryRowxContext(ctx, query, uuid).StructScan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return user, database.ErrInvalidUser
	}

	return user, err
}

func (s *SqlxDBController) GetUserByEmail(
	ctx context.Context,
	email string,
) (database.User, error) {
	query := s.db.Rebind(`
    SELECT uuid, name, email, role, subscribed
    FROM users
    WHERE (email = ?);
  `)

	var user database.User
	err := s.db.QueryRowxContext(ctx, query, email).StructScan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return user, database.ErrInvalidUser
	}

	return user, err
}
