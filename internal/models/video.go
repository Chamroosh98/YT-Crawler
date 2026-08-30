
package models

import "time"

// Video represents a YouTube video
type Video struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Channel     string    `json:"channel"`
	PublishedAt time.Time `json:"published_at"`
	URL         string    `json:"url"`
	Thumbnail   string    `json:"thumbnail,omitempty"`
	Language    string    `json:"language,omitempty"`
	FoundAt     time.Time `json:"found_at"`
}