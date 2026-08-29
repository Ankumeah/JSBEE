package database

import (
	"github.com/jmoiron/sqlx"

	"context"
	"database/sql"
	"errors"
	"uuid"
)

type SqlxDBController struct{ db *sqlx.DB }

func GetSqlxDBController(db *sqlx.DB) DBController {
	return &SqlxDBController{db}
}

// This exposes the raw underlying DB object.
// This is only to be used to execute migrations
func (s *SqlxDBController) DB() *sqlx.DB {
	return s.db
}

// Add a new user
//
// May return the following errors:
//   - `ErrExistUser`
//   - Errors by the underlying DB
func (s *SqlxDBController) AddUser(
	ctx context.Context,
	user User,
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
		return ErrExistUser
	}

	return err
}

// Delete a user
//
// May return the following errors:
//   - `ErrInvalidUser`
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
	// return `ErrInvalidUser` if no user was deleted
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return ErrInvalidUser
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
//   - `ErrInvalidUser`
//   - Errors by the underlying DB
func (s *SqlxDBController) UpdateUser(
	ctx context.Context,
	uuid uuid.UUID,
	newUser User,
) error {
	query := s.db.Rebind(`
    UPDATE users
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
	// return `ErrInvalidUser` if no user was updated
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return ErrInvalidUser
	}

	return nil
}

// Get the details of a user
//
// May return the following errors:
//   - `ErrInvalidUser`
//   - Errors by the underlying DB
func (s *SqlxDBController) GetUser(
	ctx context.Context,
	uuid uuid.UUID,
) (User, error) {
	query := s.db.Rebind(`
    SELECT uuid, name, email, role, subscribed
    FROM users
    WHERE (uuid = ?);
  `)

	var user User
	err := s.db.QueryRowxContext(ctx, query, uuid).StructScan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrInvalidUser
	}

	return user, err
}

// Adds a new unapproved paper
//
// May return the following errors:
//   - `ErrInvalidUser`
//   - `ErrExistPaper`
//   - Errors by the underlying DB
func (s *SqlxDBController) AddPaper(
	ctx context.Context,
	paper Paper,
) error {
	query := s.db.Rebind(`
    INSERT INTO papers (title, filename, owner_uuid)
    VALUES (?, ?, ?);
  `)

	if _, err := s.db.ExecContext(
		ctx, query, paper.Title, paper.Filename, paper.OwnerUUID,
	); isUniqueViolation(err) {
		return ErrExistPaper
	} else if isForeignKeyViolation(err) {
		return ErrInvalidUser
	} else if err != nil {
		return err
	}

	return nil
}

// Updates a paper to approved and gives
// it max + 1 number and max volume and issue
//
// May return the following errors:
//   - `ErrInvalidPaper`
//   - Errors by the underlying DB
func (s *SqlxDBController) ApprovePaper(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	query := s.db.Rebind(`
    UPDATE papers
    SET
      approved = 1,
      number = (SELECT COALESCE(MAX(number), 0) + 1
      volume = (SELECT volume FROM state)
      issue = (SELECT volume FROM state)
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return nil
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return ErrInvalidPaper
	}

	return nil
}

// Delete a paper
//
// May return the following errors:
//   - `ErrInvalidPaper`
//   - Errors by the underlying DB
func (s *SqlxDBController) DeletePaper(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	query := s.db.Rebind(`
    DELETE FROM papers
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return nil
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil
	} else if affected < 1 {
		return ErrInvalidPaper
	}

	return nil
}

// Get all volumes
//
// In case any paper's author has been deleted the
// User.Name feild will be `nil`
//
// May return the following errors:
//   - Errors by the underlying DB
func (s *SqlxDBController) GetVolumes(
	ctx context.Context,
) ([]Volume, error) {
	query := s.db.Rebind(`
    SELECT p.title, p.number, p.filename,
      p.volume, p.issue, u.uuid
    FROM papers AS p
    LEFT JOIN users AS u ON p.owner_uuid = u.uuid
    ORDER BY p.volume, p.issue, p.number;
  `)

	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumes []Volume
	for rows.Next() {
		var paper Paper
		var volume uint64
		var issue uint64

		if err := rows.Scan(
			&paper.Title, &paper.Number, &paper.Filename,
			&volume, &issue, &paper.OwnerUUID,
		); err != nil {
			return nil, err
		}

		// If the current paper's volume is greater then the last then create a new volume
		if len(volumes) == 0 || volumes[len(volumes)-1].Number != volume {
			volumes = append(volumes, Volume{Number: volume})
		}
		vol := &volumes[len(volumes)-1]

		// If the current paper's issue is greater then the last then create a new issue
		if len(vol.Issues) == 0 || vol.Issues[len(vol.Issues)-1].Number != issue {
			vol.Issues = append(vol.Issues, Issue{Number: issue})
		}
		iss := &vol.Issues[len(vol.Issues)-1]

		iss.Papers = append(iss.Papers, paper)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return volumes, nil
}
