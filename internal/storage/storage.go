package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Chamroosh98/YT-Crawler/internal/models"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

// New creates a new storage connection.
// If TURSO_DATABASE_URL is set, it connects to Turso.
// Otherwise falls back to local SQLite.
func New(localDBPath string) (*Storage, error) {
	var (
		db  *sql.DB
		err error
	)

	tursoURL := os.Getenv("TURSO_DATABASE_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")

	if tursoURL != "" {
		// Connect to Turso
		dsn := tursoURL
		if tursoToken != "" {
			dsn = fmt.Sprintf("%s?authToken=%s", tursoURL, tursoToken)
		}
		db, err = sql.Open("libsql", dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to Turso: %w", err)
		}
	} else {
		// Fallback to local SQLite
		if err := os.MkdirAll("data", 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
		db, err = sql.Open("sqlite", localDBPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
		if err != nil {
			return nil, fmt.Errorf("failed to open local database: %w", err)
		}
	}

	s := &Storage{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return s, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) migrate() error {
	queries := `
	CREATE TABLE IF NOT EXISTS videos (
		id           TEXT PRIMARY KEY,
		title        TEXT NOT NULL,
		channel      TEXT NOT NULL,
		published_at TEXT NOT NULL,
		url          TEXT NOT NULL,
		thumbnail    TEXT,
		language     TEXT,
		found_at     TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS queries (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		query        TEXT NOT NULL,
		language     TEXT NOT NULL DEFAULT 'en',
		topic        TEXT NOT NULL DEFAULT 'default',
		max_results  INTEGER NOT NULL DEFAULT 15,
		enabled      INTEGER NOT NULL DEFAULT 1,
		source       TEXT NOT NULL DEFAULT 'yaml',
		created_at   TEXT NOT NULL,
		updated_at   TEXT NOT NULL,
		UNIQUE(query, language, topic)
	);

	CREATE INDEX IF NOT EXISTS idx_videos_found_at ON videos(found_at);
	CREATE INDEX IF NOT EXISTS idx_queries_enabled ON queries(enabled);
	`
	_, err := s.db.Exec(queries)
	return err
}

// ==================== Video methods ====================

func (s *Storage) SaveVideo(v models.Video) (bool, error) {
	if v.FoundAt.IsZero() {
		v.FoundAt = time.Now().UTC()
	}

	res, err := s.db.Exec(`
		INSERT OR IGNORE INTO videos
		(id, title, channel, published_at, url, thumbnail, language, found_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		v.ID, v.Title, v.Channel,
		v.PublishedAt.Format(time.RFC3339),
		v.URL, v.Thumbnail, v.Language,
		v.FoundAt.Format(time.RFC3339),
	)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	return affected > 0, err
}

func (s *Storage) Exists(videoID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM videos WHERE id = ?)`, videoID).Scan(&exists)
	return exists, err
}

func (s *Storage) Count() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM videos`).Scan(&count)
	return count, err
}

// ==================== Query methods ====================

// UpsertQuery inserts or updates a query (used for YAML sync)
func (s *Storage) UpsertQuery(q models.Query) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := s.db.Exec(`
		INSERT INTO queries (query, language, topic, max_results, enabled, source, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(query, language, topic) DO UPDATE SET
			max_results = excluded.max_results,
			source = excluded.source,
			updated_at = excluded.updated_at
	`, q.Query, q.Language, q.Topic, q.MaxResults, boolToInt(q.Enabled), q.Source, now, now)

	return err
}

// GetEnabledQueries returns all enabled queries
func (s *Storage) GetEnabledQueries() ([]models.Query, error) {
	rows, err := s.db.Query(`
		SELECT id, query, language, topic, max_results, enabled, source, created_at, updated_at
		FROM queries
		WHERE enabled = 1
		ORDER BY topic, language, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanQueries(rows)
}

// GetAllQueries returns all queries (for bot management)
func (s *Storage) GetAllQueries() ([]models.Query, error) {
	rows, err := s.db.Query(`
		SELECT id, query, language, topic, max_results, enabled, source, created_at, updated_at
		FROM queries
		ORDER BY topic, language, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanQueries(rows)
}

// AddQuery adds a new query (used by Telegram bot)
func (s *Storage) AddQuery(query, language, topic string, maxResults int64) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	if topic == "" {
		topic = "default"
	}
	if maxResults <= 0 {
		maxResults = 15
	}

	res, err := s.db.Exec(`
		INSERT INTO queries (query, language, topic, max_results, enabled, source, created_at, updated_at)
		VALUES (?, ?, ?, ?, 1, 'bot', ?, ?)
	`, query, language, topic, maxResults, now, now)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// SetQueryEnabled enables or disables a query
func (s *Storage) SetQueryEnabled(id int64, enabled bool) error {
	_, err := s.db.Exec(`
		UPDATE queries SET enabled = ?, updated_at = ? WHERE id = ?
	`, boolToInt(enabled), time.Now().UTC().Format(time.RFC3339), id)
	return err
}

// DeleteQuery removes a query by ID
func (s *Storage) DeleteQuery(id int64) error {
	_, err := s.db.Exec(`DELETE FROM queries WHERE id = ?`, id)
	return err
}

// ==================== Helpers ====================

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func scanQueries(rows *sql.Rows) ([]models.Query, error) {
	var result []models.Query

	for rows.Next() {
		var q models.Query
		var enabledInt int
		var createdAt, updatedAt string

		err := rows.Scan(
			&q.ID, &q.Query, &q.Language, &q.Topic,
			&q.MaxResults, &enabledInt, &q.Source,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, err
		}

		q.Enabled = enabledInt == 1
		q.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		q.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		result = append(result, q)
	}

	return result, rows.Err()
}