
package yt

import (
	"context"
	"fmt"
	"time"

	"github.com/Chamroosh98/YT-Crawler/internal/config"
	"github.com/Chamroosh98/YT-Crawler/internal/models"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	service *youtube.Service
}

func New(apiKey string) (*Client, error) {
	ctx := context.Background()

	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube service: %w", err)
	}

	return &Client{service: service}, nil
}

// Search searches YouTube and returns a list of videos
func (c *Client) Search(ctx context.Context, q config.Query, daysLookback int) ([]models.Video, error) {
	publishedAfter := time.Now().UTC().AddDate(0, 0, -daysLookback).Format(time.RFC3339)

	call := c.service.Search.List([]string{"id", "snippet"}).
		Q(q.Query).
		Type("video").
		Order("date").
		MaxResults(q.MaxResults).
		PublishedAfter(publishedAfter)

	resp, err := call.Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("youtube search failed for query '%s': %w", q.Query, err)
	}

	var videos []models.Video

	for _, item := range resp.Items {
		if item.Id.Kind != "youtube#video" {
			continue
		}

		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			publishedAt = time.Now().UTC()
		}

		video := models.Video{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Channel:     item.Snippet.ChannelTitle,
			PublishedAt: publishedAt,
			URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId),
			Thumbnail:   getBestThumbnail(item.Snippet.Thumbnails),
			Language:    q.Language,
			FoundAt:     time.Now().UTC(),
		}

		videos = append(videos, video)
	}

	return videos, nil
}

func getBestThumbnail(thumbs *youtube.ThumbnailDetails) string {
	if thumbs == nil {
		return ""
	}
	switch {
	case thumbs.High != nil:
		return thumbs.High.Url
	case thumbs.Medium != nil:
		return thumbs.Medium.Url
	case thumbs.Default != nil:
		return thumbs.Default.Url
	default:
		return ""
	}
}