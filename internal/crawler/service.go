
package crawler

import (
	"context"
	"log"
	"time"

	"github.com/Chamroosh98/YT-Crawler/internal/config"
	"github.com/Chamroosh98/YT-Crawler/internal/models"
	"github.com/Chamroosh98/YT-Crawler/internal/notifier"
	"github.com/Chamroosh98/YT-Crawler/internal/storage"
	"github.com/Chamroosh98/YT-Crawler/internal/yt"
)

type Service struct {
	cfg      *config.Config
	store    *storage.Storage
	ytClient *yt.Client
	telegram *notifier.Telegram
}

func New(cfg *config.Config, store *storage.Storage, ytClient *yt.Client, telegram *notifier.Telegram) *Service {
	return &Service{
		cfg:      cfg,
		store:    store,
		ytClient: ytClient,
		telegram: telegram,
	}
}

func (s *Service) Run(ctx context.Context) error {
	log.Println("  🚀 Starting crawler...")

	dbQueries, err := s.store.GetEnabledQueries()
	if err != nil {
		return err
	}

	if len(dbQueries) == 0 {
		log.Println("  ⚠️  No enabled queries found in database")
		return nil
	}

	log.Printf("  📋 Loaded %d enabled queries from database", len(dbQueries))

	var totalFound int
	var newVideos []models.Video

	for i, dbq := range dbQueries {
		q := config.Query{
			Query:      dbq.Query,
			Language:   dbq.Language,
			MaxResults: dbq.MaxResults,
		}

		log.Printf("(%d/%d) Searching : [%s] %s", i+1, len(dbQueries), q.Language, q.Query)

		videos, err := s.ytClient.Search(ctx, q, s.cfg.Search.DaysLookback)
		if err != nil {
			log.Printf("  ❌ Search failed : %v", err)
			continue
		}

		totalFound += len(videos)
		newCount := 0

		for _, video := range videos {
			isNew, err := s.store.SaveVideo(video)
			if err != nil {
				log.Printf("  ⚠️  Failed to save video %s : %v", video.ID, err)
				continue
			}

			if isNew {
				newCount++
				newVideos = append(newVideos, video)
			}
		}

		log.Printf("   → Found : %d | New: %d", len(videos), newCount)
		time.Sleep(300 * time.Millisecond)
	}

	// Send one summary instead of individual messages
	if len(newVideos) > 0 {
		if err := s.telegram.SendSummary(newVideos); err != nil {
			log.Printf("  ⚠️  Failed to send summary: %v", err)
		} else {
			log.Printf("  📨 Summary sent (%d new videos)", len(newVideos))
		}
	} else {
		log.Println("  ℹ️  No new videos to notify")
	}

	log.Println("  =====================================")
	log.Printf("  ✅ Crawler finished")
	log.Printf("    Total videos found : %d", totalFound)
	log.Printf("    New videos saved   : %d", len(newVideos))
	log.Println("  =====================================")

	return nil
}

func (s *Service) RunOnce(ctx context.Context, query config.Query) ([]models.Video, error) {
	return s.ytClient.Search(ctx, query, s.cfg.Search.DaysLookback)
}