package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"
	"github.com/Ankumeah/JSBEE/backend/internal/roles"

	"context"
	"uuid"
)

func (s *SqlxDBController) ChangeRole(
	ctx context.Context,
	userUUID uuid.UUID,
	newRole roles.Role,
) error {
	query := s.db.Rebind(`
    UPDATE users
    SET role = ?
    WHERE uuid = ?;
  `)

	res, err := s.db.ExecContext(ctx, query, newRole.Role, userUUID)
	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected < 1 {
		return database.ErrInvalidUser
	}

	return nil
}
