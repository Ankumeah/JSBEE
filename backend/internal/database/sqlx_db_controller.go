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
	name string,
) (User, error) {
	query := s.db.Rebind(`
    SELECT name, email, role
    FROM users
    WHERE name = ?;
  `)

	var user User
	err := s.db.QueryRowxContext(ctx, query, name).Scan(&user)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrInvalidUser
	}

	return user, err
}

func (s *SqlxDBController) GetVolumes(
	ctx context.Context,
) ([]Volume, error) {
	query := s.db.Rebind(`
    SELECT p.title, p.number, p.filename,
      p.volume, p.issue, u.name
    FROM papers AS p
    LEFT JOIN users AS u ON p.owner_id = u.id
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
		var owner sql.NullString

		if err := rows.Scan(
			&paper.Title, &paper.Number, &paper.Filename,
			&volume, &issue, &owner,
		); err != nil {
			return nil, err
		}

		if owner.Valid {
			paper.Owner = owner.String
		}

		if len(volumes) == 0 || volumes[len(volumes)-1].Number != volume {
			volumes = append(volumes, Volume{Number: volume})
		}
		vol := &volumes[len(volumes)-1]

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
