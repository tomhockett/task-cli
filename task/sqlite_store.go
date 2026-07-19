package task

import (
	"database/sql"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

// Compile-time check that SQLiteStore implements TaskStore.
// If SQLiteStore is missing any interface methods, this line will fail to compile.
var _ TaskStore = (*SQLiteStore)(nil)

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
    							id INTEGER PRIMARY KEY AUTOINCREMENT,
    							title TEXT,
    							status INTEGER,
    							priority INTEGER,
    							tags TEXT,
    							created_at DATETIME,
    							completed_at DATETIME
					)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	// Lightweight migration: ensure the "tags" column exists on existing databases.
	// If the column already exists, SQLite will return a "duplicate column name" error,
	// which we safely ignore.
	if _, err := db.Exec(`ALTER TABLE tasks ADD COLUMN tags TEXT DEFAULT ''`); err != nil {
		if !strings.Contains(err.Error(), "duplicate column name") {
			db.Close()
			return nil, err
		}
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Add(title string, opts AddOptions) (Task, error) {
	now := time.Now()
	var priority Priority
	if opts.Priority != nil {
		priority = *opts.Priority
	} else {
		priority = PriorityMedium
	}
	var tags string
	if opts.Tags != nil {
		tags = strings.Join(opts.Tags, ",")
	}
	result, err := s.db.Exec("INSERT INTO tasks (title, status, priority, tags, created_at) VALUES (?, ?, ?, ?, ?)", title, StatusTodo, priority, tags, now)
	if err != nil {
		return Task{}, err
	}
	id, err := result.LastInsertId() // this gets the auto-generated ID
	if err != nil {
		return Task{}, err
	}

	t := Task{
		ID:        int(id), // from LastInsertId()
		Title:     title,
		Status:    StatusTodo, // no prefix needed, same package
		Priority:  priority,
		Tags:      opts.Tags,
		CreatedAt: now,
	}
	return t, nil
}

func (s *SQLiteStore) Complete(id int) error {
	result, err := s.db.Exec("UPDATE tasks SET status = ?, completed_at = ? WHERE id = ?", StatusDone, time.Now(), id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (s *SQLiteStore) Delete(id int) error {
	result, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (s *SQLiteStore) List(opts ListOptions) ([]Task, error) {
	// Build the query dynamically — like an AR scope chain adding WHERE clauses.
	query := "SELECT id, title, status, priority, tags, created_at, completed_at FROM tasks"
	var args []any
	if opts.Status != nil {
		query += " WHERE status = ?"
		args = append(args, *opts.Status)
	}
	query += " ORDER BY id"
	results, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	var tasks []Task
	for results.Next() {
		var t Task
		var tagsStr sql.NullString
		err := results.Scan(&t.ID, &t.Title, &t.Status, &t.Priority, &tagsStr, &t.CreatedAt, &t.CompletedAt)
		if err != nil {
			return nil, err
		}
		if tagsStr.Valid && tagsStr.String != "" {
			t.Tags = strings.Split(tagsStr.String, ",")
		}
		// Tags are stored comma-joined, so SQL LIKE can't do an exact match
		// ("work" would also match "workout"). Filter after scanning instead.
		if !opts.matches(t) {
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
