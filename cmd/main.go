
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chamroosh98/YT-Crawler/internal/config"
	"github.com/Chamroosh98/YT-Crawler/internal/crawler"
	"github.com/Chamroosh98/YT-Crawler/internal/notifier"
	"github.com/Chamroosh98/YT-Crawler/internal/storage"
	"github.com/Chamroosh98/YT-Crawler/internal/yt"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// 1. Load config (YAML)
	cfg, err := config.Load("config")
	if err != nil {
		log.Fatalf("  ❌ failed to load config : %v", err)
	}
	log.Printf("  ✅ Config loaded | YAML queries : %d", len(cfg.Queries))

	// 2. Storage (Turso or local SQLite)
	store, err := storage.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("  ❌ failed to init storage : %v", err)
	}
	defer store.Close()

	if os.Getenv("TURSO_DATABASE_URL") != "" {
		log.Println("  ✅ Connected to Turso!")
	} else {
		log.Println("  ℹ️ Using local SQLite!")
	}

	// 3. Sync YAML → Database
	if err := store.SyncFromYAML(cfg); err != nil {
		log.Fatalf("  ❌ sync failed : %v", err)
	}

	// 4. YouTube client
	ytClient, err := yt.New(cfg.YouTubeAPIKey)
	if err != nil {
		log.Fatalf("  ❌ failed to create youtube client : %v", err)
	}
	log.Println("  ✅ YouTube client ready")

	// 5. Telegram notifier
	telegram := notifier.NewTelegram(cfg.TelegramBotToken, cfg.TelegramChatID)
	if telegram.Enabled() {
		log.Println("  ✅ Telegram notifier enabled")
	} else {
		log.Println("  ℹ️ Telegram notifier disabled")
	}

	// 6. Crawler service
	service := crawler.New(cfg, store, ytClient, telegram)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := service.Run(ctx); err != nil {
		log.Fatalf("  ❌ crawler failed : %v", err)
	}
}