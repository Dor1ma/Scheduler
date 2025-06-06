package postgres

import (
	"database/sql"
	"github.com/Dor1ma/Scheduler/internal/timer"
	"time"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Save(task timer.Task) error {
	_, err := s.db.Exec(`
        INSERT INTO tasks (id, execute_at, payload, url, method)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE
        SET execute_at = EXCLUDED.execute_at,
            payload = EXCLUDED.payload,
            url = EXCLUDED.url,
            method = EXCLUDED.method
    `, task.ID, task.ExecuteAt, task.Payload, task.URL, task.Method)
	return err
}

func (s *PostgresStore) Delete(id string) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	return err
}

func (s *PostgresStore) LoadPending(until time.Time) ([]timer.Task, error) {
	rows, err := s.db.Query(`
        SELECT id, execute_at, payload, url, method 
        FROM tasks 
        WHERE execute_at <= $1
    `, until)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []timer.Task
	for rows.Next() {
		var t timer.Task
		if err := rows.Scan(&t.ID, &t.ExecuteAt, &t.Payload, &t.URL, &t.Method); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
