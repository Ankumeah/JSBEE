package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"context"
	"database/sql"
	"errors"
	"uuid"
)

// Add a new user
//
// May return the following errors:
//   - `database.ErrExistUser`
//   - Errors by the underlying DB
func (s *SqlxDBController) AddUser(
	ctx context.Context,
	user database.User,
) error {
	query := s.db.Rebind(`
    INSERT INTO users (uuid, name, email, subscribed, role)
    VALUES (?, ?, ?, ?, ?);
  `)

	_, err := s.db.ExecContext(
		ctx, query,
		user.UUID, user.Name, user.Email, user.Subscribed, user.Role.Role,
	)

	if isUniqueViolation(err) {
		return database.ErrExistUser
	}

	return err
}

// Delete a user
//
// May return the following errors:
//   - `database.ErrInvalidUser`
//   - Errors by the underlying DB
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

// Update a user
// This should not be directly exposed as a user
// user accessible API as it has the power to
// change a user's role. This function does not update
// a user's UUID
//
// May return the following errors:
//   - `database.ErrInvalidUser`
//   - Errors by the underlying DB
func (s *SqlxDBController) UpdateUser(
	ctx context.Context,
	uuid uuid.UUID,
	newUser database.User,
) error {
	query := s.db.Rebind(`
    UPDATE users
  "github.com/Ankumeah/JSBEE/backend/internal/database"
    SET name = ?, email = ?, role = ?, subscribed = ?
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(
		ctx, query,
		newUser.Name, newUser.Email, newUser.Role, newUser.Subscribed, uuid,
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

// Get the details of a user
//
// May return the following errors:
//   - `database.ErrInvalidUser`
//   - Errors by the underlying DB
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
