package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"context"
	"database/sql"
	"errors"
	"uuid"
)

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

func (s *SqlxDBController) GetUserPapers(
	ctx context.Context,
	userUUID uuid.UUID,
) ([]database.Paper, error) {
	query := s.db.Rebind(`
    SELECT uuid, title, number, filename, owner_uuid
    FROM papers
    WHERE (owner_uuid = ?);
  `)

	var papers []database.Paper
	err := s.db.SelectContext(ctx, &papers, query, userUUID)

	return papers, err
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
	err := s.db.QueryRowxContext(ctx, query, uuid).StructScan(&paper)
	if errors.Is(err, sql.ErrNoRows) {
		return paper, database.ErrInvalidPaper
	}

	return paper, err
}

func (s *SqlxDBController) GetVolumes(
	ctx context.Context,
) ([]database.Volume, error) {
	query := `
    SELECT title, number, filename, volume, issue, uuid, owner_uuid
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
			&volume, &issue, &paper.UUID, &paper.OwnerUUID,
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
