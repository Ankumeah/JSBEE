package main

import (
	"context"
	"database/sql"
	"uuid"

	_ "modernc.org/sqlite"
)

const queueSchema = `
  PRAGMA foreign_keys = ON;

  CREATE TABLE IF NOT EXITS batches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT NOT NULL UNIQUE CHECK (uuid != ''),
    creation INTEGER NOT NULL
  );

  CREATE TABLE IF NOT EXISTS jobs (
    uuid TEXT NOT NULL PRIMARY KEY CHECK (uuid != ''),
    type TEXT NOT NULL CHECK (uuid != ''),
    user TEXT NOT NULL CHECK (uuid != ''),
    email TEXT NOT NULL CHECK (uuid != ''),
    paper_title TEXT NOT NULL CHECK (uuid != ''),
    creation INTEGER NOT NULL,
    feedback TEXT,
    batch TEXT REFERENCES users(uuid) ON DELETE CASCADE
  );
`

type Job struct {
	UUID       uuid.UUID
	Type       string
	User       string
	Email      string
	PaperTitle string
	Creation   int64
	Feedback   *string
	Batch      *string
}

type Queue struct {
	db    *sql.DB
	queue chan Job
}

func loadQueue(ctx context.Context) (*Queue, error) {
	db, err := sql.Open("sqlite", envVars["DB_URL"])
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, queueSchema); err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, `
    SELECT
      uuid, type, user,
      email, paper_title,
      creation, feedback, batch
    FROM jobs
    SORT BY creation;
  `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make(chan Job)
	for rows.Next() {
		var job Job
		if err := rows.Scan(
			&job.UUID, &job.Type, &job.User,
			&job.Email, &job.PaperTitle,
			&job.Creation, &job.Feedback, &job.Batch,
		); err != nil {
			return nil, err
		}

		jobs <- job
	}

	return &Queue{db, jobs}, err
}

func (q *Queue) Enqueue(ctx context.Context, job Job) error {
	if _, err := q.db.ExecContext(ctx,
		`INSERT INTO jobs (
      uuid, type, user,
      email, paper_title,
      creation, feedback, batch
    )
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,

		job.UUID, job.Type, job.User,
		job.Email, job.PaperTitle,
		job.Creation, job.Feedback, job.Batch,
	); err != nil {
		return err
	}
	q.queue <- job

	return nil
}
