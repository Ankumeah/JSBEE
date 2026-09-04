package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"context"
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
