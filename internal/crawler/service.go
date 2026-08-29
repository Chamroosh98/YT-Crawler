
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
	log.Println("🚀 Starting crawler ...")

	var totalFound int
	var totalNew int

	for i, query := range s.cfg.Queries {
		log.Printf("(%d/%d) Searching: [%s] %s", i+1, len(s.cfg.Queries), query.Language, query.Query)

		videos, err := s.ytClient.Search(ctx, query, s.cfg.Search.DaysLookback)
		if err != nil {
			log.Printf("❌ Search failed : %v", err)
			continue
		}

		totalFound += len(videos)
		newCount := 0

		for _, video := range videos {
			isNew, err := s.store.SaveVideo(video)
			if err != nil {
				log.Printf("⚠️  Failed to save video %s : %v", video.ID, err)
				continue
			}

			if isNew {
				newCount++
				totalNew++

				if err := s.telegram.SendVideo(video); err != nil {
					log.Printf("⚠️  Failed to send telegram notification : %v", err)
				}

				time.Sleep(400 * time.Millisecond)
			}
		}

		log.Printf("   → Found: %d | New: %d", len(videos), newCount)

		time.Sleep(300 * time.Millisecond)
	}

	log.Println("=====================================")
	log.Printf("✅ Crawler finished")
	log.Printf("   Total videos found : %d", totalFound)
	log.Printf("   New videos saved   : %d", totalNew)
	log.Println("=====================================")

	return nil
}

func (s *Service) RunOnce(ctx context.Context, query config.Query) ([]models.Video, error) {
	return s.ytClient.Search(ctx, query, s.cfg.Search.DaysLookback)
}