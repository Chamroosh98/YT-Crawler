
package storage

import (
	"log"

	"github.com/Chamroosh98/YT-Crawler/internal/config"
	"github.com/Chamroosh98/YT-Crawler/internal/models"
)

// SyncFromYAML reads queries from config and upserts them into the database.
// Existing bot-managed queries are not deleted.
func (s *Storage) SyncFromYAML(cfg *config.Config) error {
	log.Println("🔄 Syncing queries from YAML to database...")

	count := 0
	for _, q := range cfg.Queries {
		mq := models.Query{
			Query:      q.Query,
			Language:   q.Language,
			Topic:      "default", // later we can make this dynamic
			MaxResults: q.MaxResults,
			Enabled:    true,
			Source:     "yaml",
		}

		if err := s.UpsertQuery(mq); err != nil {
			log.Printf("⚠️  failed to upsert query '%s': %v", q.Query, err)
			continue
		}
		count++
	}

	log.Printf("✅ Synced %d queries from YAML", count)
	return nil
}