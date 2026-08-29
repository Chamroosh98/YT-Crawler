
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
)

func main() {
	
	cfg, err := config.Load("config")
	if err != nil {
		log.Fatalf("❌ failed to load config: %v", err)
	}
	log.Printf("✅ Config loaded | Queries: %d", len(cfg.Queries))

	// Storage
	store, err := storage.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("❌ failed to init storage: %v", err)
	}
	defer store.Close()
	log.Println("✅ Storage initialized")

	// YouTube Client
	ytClient, err := yt.New(cfg.YouTubeAPIKey)
	if err != nil {
		log.Fatalf("❌ failed to create youtube client: %v", err)
	}
	log.Println("✅ YouTube client ready")

	// Telegram Notifier
	telegram := notifier.NewTelegram(cfg.TelegramBotToken, cfg.TelegramChatID)
	if telegram.Enabled() {
		log.Println("✅ Telegram notifier enabled")
	} else {
		log.Println("ℹ️  Telegram notifier disabled (missing token or chat_id)")
	}

	service := crawler.New(cfg, store, ytClient, telegram)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := service.Run(ctx); err != nil {
		log.Fatalf("❌ crawler failed: %v", err)
	}
}