package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"context"
	"database/sql"
	"errors"
	"uuid"
)

// Adds a new unapproved paper
//
// May return the following errors:
//   - `database.ErrInvalidUser`
//   - `database.ErrExistPaper`
//   - Errors by the underlying DB
func (s *SqlxDBController) AddPaper(
	ctx context.Context,
	paper database.Paper,
) error {
	query := s.db.Rebind(`
    INSERT INTO papers (uuid, title, filename, owner_uuid)
    VALUES (?, ?, ?, ?);
  `)

	if _, err := s.db.ExecContext(
		ctx, query, paper.UUID, paper.Title, paper.Filename, paper.OwnerUUID,
	); isUniqueViolation(err) {
		return database.ErrExistPaper
	} else if isForeignKeyViolation(err) {
		return database.ErrInvalidUser
	} else if err != nil {
		return err
	}

	return nil
}

// Updates a paper to max + 1 number and max volume and issue
//
// May return the following errors:
//   - `database.ErrInvalidPaper`
//   - Errors by the underlying DB
func (s *SqlxDBController) ApprovePaper(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	query := s.db.Rebind(`
    UPDATE papers
    SET
      number = (SELECT COALESCE(MAX(number), 0) FROM papers) + 1,
      volume = (SELECT volume FROM state),
      issue = (SELECT volume FROM state)
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return database.ErrInvalidPaper
	}

	return nil
}

// Delete a paper
//
// May return the following errors:
//   - `database.ErrInvalidPaper`
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
		return database.ErrInvalidPaper
	}

	return nil
}

func (s *SqlxDBController) GetPaper(
	ctx context.Context,
	uuid uuid.UUID,
) (database.Paper, error) {
	query := s.db.Rebind(`
    SELECT uuid, title, number, filename, owner_uuid
    FROM papers
    WHERE (uuid = ?);
  `)

	var paper database.Paper
	err := s.db.QueryRowxContext(ctx, query, uuid).Scan(&paper)
	if errors.Is(err, sql.ErrNoRows) {
		return paper, database.ErrInvalidPaper
	}

	return paper, err
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
) ([]database.Volume, error) {
	query := `
    SELECT title, number, filename, volume, issue, uuid
    FROM papers
    WHERE (number IS NOT NULL)
    ORDER BY volume, issue, number;
  `

	rows, err := s.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumes []database.Volume
	for rows.Next() {
		var paper database.Paper
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
			volumes = append(volumes, database.Volume{Number: volume})
		}
		vol := &volumes[len(volumes)-1]

		// If the current paper's issue is greater then the last then create a new issue
		if len(vol.Issues) == 0 || vol.Issues[len(vol.Issues)-1].Number != issue {
			vol.Issues = append(vol.Issues, database.Issue{Number: issue})
		}
		iss := &vol.Issues[len(vol.Issues)-1]

		iss.Papers = append(iss.Papers, paper)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return volumes, nil
}

func (s *SqlxDBController) GetUnapprovedPapers(
	ctx context.Context,
) ([]database.Paper, error) {
	query := `
    SELECT uuid, title, number, filename, owner_id
    FROM papers
    WHERE (number IS NULL)
  `

	var papers []database.Paper
	err := s.db.SelectContext(ctx, papers, query)

	return papers, err
}
