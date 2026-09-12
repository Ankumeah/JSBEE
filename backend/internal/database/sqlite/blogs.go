package sqlite

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"context"
	"database/sql"
	"errors"
	"uuid"
)

func (s *SqlxDBController) AddBlog(
	ctx context.Context,
	blog database.Blog,
) error {
	query := s.db.Rebind(`
    INSERT INTO blogs (uuid, title, filename, created_at, updated_at)
    VALUES (?, ?, ?, ?, ?);
  `)

	if _, err := s.db.ExecContext(
		ctx, query,
		blog.UUID, blog.Title, blog.Filename, blog.CreatedAt, blog.UpdatedAt,
	); isUniqueViolation(err) {
		return database.ErrExistBlog
	} else if err != nil {
		return err
	}

	return nil
}

func (s *SqlxDBController) GetBlog(
	ctx context.Context,
	uuid uuid.UUID,
) (database.Blog, error) {
	query := s.db.Rebind(`
    SELECT uuid, title, filename, created_at, updated_at
    FROM blogs
    WHERE (uuid = ?);
  `)

	var blog database.Blog
	err := s.db.QueryRowxContext(ctx, query, uuid).StructScan(&blog)
	if errors.Is(err, sql.ErrNoRows) {
		return blog, database.ErrInvalidBlog
	}

	return blog, err
}

func (s *SqlxDBController) GetBlogs(
	ctx context.Context,
) ([]database.Blog, error) {
	query := s.db.Rebind(`
    SELECT uuid, title, filename, created_at, updated_at
    FROM blogs
    ORDER BY created_at DESC, uuid;
  `)

	var blogs []database.Blog
	err := s.db.SelectContext(ctx, &blogs, query)

	return blogs, err
}

func (s *SqlxDBController) UpdateBlog(
	ctx context.Context,
	blog database.Blog,
) error {
	query := s.db.Rebind(`
    UPDATE blogs
    SET title = ?, filename = ?, updated_at = ?
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(
		ctx, query, blog.Title, blog.Filename, blog.UpdatedAt, blog.UUID,
	)
	if err != nil {
		return err
	}
	// return `database.ErrInvalidBlog` if no blog was updated
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return database.ErrInvalidBlog
	}

	return nil
}

func (s *SqlxDBController) DeleteBlog(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	query := s.db.Rebind(`
    DELETE FROM blogs
    WHERE (uuid = ?);
  `)

	res, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return err
	}
	// return `database.ErrInvalidBlog` if no blog was deleted
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected < 1 {
		return database.ErrInvalidBlog
	}

	return nil
}
