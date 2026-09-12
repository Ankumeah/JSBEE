package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"context"
	"uuid"
)

func (s *SqlxDBController) IncrementVolume(ctx context.Context) error {
	query := `
    UPDATE state
    SET volume = volume + 1, issue = 1
    WHERE (id = 1);
  `

	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SqlxDBController) IncrementIssue(ctx context.Context) error {
	query := `
    UPDATE state
    SET issue = issue + 1
    WHERE (id = 1);
  `

	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SqlxDBController) GetUnapprovedPapers(
	ctx context.Context,
) ([]database.Paper, error) {
	query := `
    SELECT uuid, title, number, filename, owner_uuid, reviewed
    FROM papers
    WHERE (number IS NULL AND reviewed = 0)
  `

	var papers []database.Paper = []database.Paper{}
	err := s.db.SelectContext(ctx, &papers, query)

	return papers, err
}

func (s *SqlxDBController) GetReviewedPapers(
	ctx context.Context,
) ([]database.Paper, error) {
	query := `
    SELECT uuid, title, number, filename, owner_uuid, reviewed
    FROM papers
    WHERE (number IS NULL AND reviewed = 1)
  `

	var papers []database.Paper = []database.Paper{}
	err := s.db.SelectContext(ctx, &papers, query)

	return papers, err
}

func (s *SqlxDBController) ApprovePaper(
	ctx context.Context,
	paperUUID uuid.UUID,
) error {
	query := s.db.Rebind(`
    UPDATE papers
    SET reviewed = 1
    WHERE (uuid = ? AND number IS NULL);
  `)

	res, err := s.db.ExecContext(ctx, query, paperUUID)
	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected < 1 {
		return database.ErrInvalidPaper
	}

	return nil
}

func (s *SqlxDBController) PublishPapers(
	ctx context.Context,
) (int64, error) {
	query := `
    WITH cur_state AS (
      SELECT volume, issue
      FROM state
      WHERE id = 1
    ),
    reviewed AS (
      SELECT
        id,
        ROW_NUMBER() OVER (ORDER BY id) - 1 AS rn
      FROM papers
      WHERE
        reviewed = 1
        AND number IS NULL
    ),
    current_max AS (
      SELECT COALESCE(MAX(p.number), 0) AS max_number
      FROM papers AS p
      CROSS JOIN cur_state AS cs
      WHERE
        p.volume = cs.volume
        AND p.issue = cs.issue
    )
    UPDATE papers AS p
    SET
      number = current_max.max_number + reviewed.rn + 1,
      volume = cur_state.volume,
      issue = cur_state.issue
    FROM reviewed
    CROSS JOIN cur_state
    CROSS JOIN current_max
    WHERE p.id = reviewed.id;
  `

	res, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (s *SqlxDBController) RejectPaper(
	ctx context.Context,
	paperUUID uuid.UUID,
) error {
	query := s.db.Rebind(`
    DELETE FROM papers
    WHERE (
      number IS NULL AND
      uuid = ?
    );
  `)

	res, err := s.db.ExecContext(ctx, query, paperUUID)
	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected < 1 {
		return database.ErrInvalidPaper
	}

	return nil
}
