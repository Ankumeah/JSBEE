package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"context"
	"uuid"
)

func (s *SqlxDBController) IncrementVolume(ctx context.Context) error {
	query := `
    UPDATE state
    SET volume = volume + 1,
    WHERE (id = 1);
  `

	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SqlxDBController) IncrementIssue(ctx context.Context) error {
	query := `
    UPDATE state
    SET issue = issue + 1,
    WHERE (id = 1);
  `

	_, err := s.db.ExecContext(ctx, query)
	return err
}

func (s *SqlxDBController) GetUnapprovedPapers(
	ctx context.Context,
) ([]database.Paper, error) {
	query := `
    SELECT uuid, title, number, filename, owner_uuid
    FROM papers
    WHERE (number IS NULL)
  `

	var papers []database.Paper
	err := s.db.SelectContext(ctx, &papers, query)

	return papers, err
}

func (s *SqlxDBController) ApprovePaper(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	query := s.db.Rebind(`
    UPDATE papers
    SET
      number = COALESCE(number, (SELECT COALESCE(MAX(number), 0) FROM papers) + 1),
      volume = COALESCE(volume, (SELECT volume FROM state)),
      issue = COALESCE(issue, (SELECT issue FROM state))
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(ctx, query, uuid)
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
