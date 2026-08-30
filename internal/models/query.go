package models

import "time"

// Query represents a search query stored in the database
type Query struct {
	ID         int64     `json:"id"`
	Query      string    `json:"query"`
	Language   string    `json:"language"`
	Topic      string    `json:"topic"` 			// for future multi-topic support
	MaxResults int64     `json:"max_results"`
	Enabled    bool      `json:"enabled"`
	Source     string    `json:"source"` 		// "yaml" or "bot"
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}