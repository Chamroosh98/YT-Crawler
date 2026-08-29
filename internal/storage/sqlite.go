
package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/Chamroosh98/YT-Crawler/internal/models"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	s := &Storage{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) migrate() error {
	query := `
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

	CREATE INDEX IF NOT EXISTS idx_videos_published_at ON videos(published_at);
	CREATE INDEX IF NOT EXISTS idx_videos_found_at ON videos(found_at);
	`
	_, err := s.db.Exec(query)
	return err
}

func (s *Storage) SaveVideo(v models.Video) (bool, error) {
	if v.FoundAt.IsZero() {
		v.FoundAt = time.Now().UTC()
	}

	res, err := s.db.Exec(`
		INSERT OR IGNORE INTO videos 
		(id, title, channel, published_at, url, thumbnail, language, found_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		v.ID,
		v.Title,
		v.Channel,
		v.PublishedAt.Format(time.RFC3339),
		v.URL,
		v.Thumbnail,
		v.Language,
		v.FoundAt.Format(time.RFC3339),
	)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (s *Storage) Exists(videoID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM videos WHERE id = ?)`, videoID).Scan(&exists)
	return exists, err
}

func (s *Storage) GetRecentVideos(limit int) ([]models.Video, error) {
	rows, err := s.db.Query(`
		SELECT id, title, channel, published_at, url, thumbnail, found_at
		FROM videos
		ORDER BY found_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []models.Video
	for rows.Next() {
		var v models.Video
		var publishedAt, foundAt string

		err := rows.Scan(
			&v.ID, &v.Title, &v.Channel, &publishedAt,
			&v.URL, &v.Thumbnail, &foundAt,
		)
		if err != nil {
			return nil, err
		}

		v.PublishedAt, _ = time.Parse(time.RFC3339, publishedAt)
		v.FoundAt, _ = time.Parse(time.RFC3339, foundAt)
		videos = append(videos, v)
	}

	return videos, rows.Err()
}

func (s *Storage) Count() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM videos`).Scan(&count)
	return count, err
}