package sqlite

import (
	"context"
	"os"
)

func (s *SqlxDBController) BackupDB(ctx context.Context) (string, error) {
	// VACUUM INTO fails if the target file exists,
	// so clear any leftover from a previous run first.
	_ = os.Remove(backupFile)

	_, err := s.db.ExecContext(
		ctx, "VACUUM INTO '"+backupFile+"'",
	)

	return backupFile, err
}
