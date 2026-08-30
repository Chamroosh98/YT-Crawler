
package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Chamroosh98/YT-Crawler/internal/config"
	"github.com/Chamroosh98/YT-Crawler/internal/storage"
	"github.com/joho/godotenv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load("config")
	if err != nil {
		log.Fatalf("failed to load config : %v", err)
	}

	store, err := storage.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to init storage : %v", err)
	}
	defer store.Close()

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	// Only allow commands from your chat
	allowedChatID := os.Getenv("TELEGRAM_CHAT_ID")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("failed to create bot : %v", err)
	}

	log.Printf("✅ Bot authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	// Graceful shutdown
	go func() {
		mod := make(chan os.Signal, 1)
		signal.Notify(mod, syscall.SIGINT, syscall.SIGTERM)
		<-mod
		log.Println("Shutting down bot ...")
		bot.StopReceivingUpdates()
		os.Exit(0)
	}()

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		// Security: only accept commands from allowed chat
		if allowedChatID != "" && strconv.FormatInt(update.Message.Chat.ID, 10) != allowedChatID {
			continue
		}

		chatID := update.Message.Chat.ID
		cmd := update.Message.Command()
		args := strings.TrimSpace(update.Message.CommandArguments())

		var reply string

		switch cmd {
		case "start", "help":
			reply = helpText()

		case "list":
			reply = handleList(store, args)

		case "add":
			reply = handleAdd(store, args)

		case "remove", "rm", "del":
			reply = handleRemove(store, args)

		case "enable":
			reply = handleSetEnabled(store, args, true)

		case "disable":
			reply = handleSetEnabled(store, args, false)

		default:
			reply = "Unknown command.\n\n" + helpText()
		}

		msg := tgbotapi.NewMessage(chatID, reply)
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Printf("failed to send message : %v", err)
		}
	}
}

func helpText() string {
	return `🤖 <b>Query Manager Bot</b>

Available commands:

/list [lang] — Show queries (optional language filter)
/add &lt;lang&gt; &lt;query&gt; — Add a new query
/remove &lt;id&gt; — Delete a query by ID
/enable &lt;id&gt; — Enable a query
/disable &lt;id&gt; — Disable a query
/help — Show this message

Examples:
<code>/add fa کلودفلر ورکرز</code>
<code>/add en cloudflare workers tutorial</code>
<code>/list fa</code>
<code>/disable 3</code>`
}

func handleList(store *storage.Storage, langFilter string) string {
	queries, err := store.GetAllQueries()
	if err != nil {
		return "❌ Failed to load queries : " + err.Error()
	}

	if len(queries) == 0 {
		return "No queries found!"
	}

	var b strings.Builder
	b.WriteString("📋 <b>Queries</b>\n\n")

	count := 0
	for _, q := range queries {
		if langFilter != "" && !strings.EqualFold(q.Language, langFilter) {
			continue
		}

		status := "✅"
		if !q.Enabled {
			status = "❌"
		}

		b.WriteString(strconv.FormatInt(q.ID, 10))
		b.WriteString(". ")
		b.WriteString(status)
		b.WriteString(" [")
		b.WriteString(q.Language)
		b.WriteString("] ")
		b.WriteString(htmlEscape(q.Query))
		b.WriteString(" <i>(")
		b.WriteString(q.Source)
		b.WriteString(")</i>\n")
		count++
	}

	if count == 0 {
		return "No queries matched the filter."
	}

	return b.String()
}

func handleAdd(store *storage.Storage, args string) string {
	parts := strings.SplitN(args, " ", 2)
	if len(parts) < 2 {
		return "Usage: /add &lt;lang&gt; &lt;query&gt;\nExample: /add fa کلودفلر ورکرز"
	}

	lang := strings.ToLower(strings.TrimSpace(parts[0]))
	query := strings.TrimSpace(parts[1])

	if lang == "" || query == "" {
		return "Language and query cannot be empty!"
	}

	id, err := store.AddQuery(query, lang, "default", 15)
	if err != nil {
		return "❌ Failed to add query : " + err.Error()
	}

	return "✅ Query added with ID <b>" + strconv.FormatInt(id, 10) + "</b>\n[" + lang + "] " + htmlEscape(query)
}

func handleRemove(store *storage.Storage, args string) string {
	id, err := strconv.ParseInt(strings.TrimSpace(args), 10, 64)
	if err != nil {
		return "Usage : /remove &lt;id&gt;"
	}

	if err := store.DeleteQuery(id); err != nil {
		return "❌ Failed to remove query : " + err.Error()
	}

	return "🗑 Query <b>" + strconv.FormatInt(id, 10) + "</b> removed."
}

func handleSetEnabled(store *storage.Storage, args string, enabled bool) string {
	id, err := strconv.ParseInt(strings.TrimSpace(args), 10, 64)
	if err != nil {
		if enabled {
			return "Usage : /enable &lt;id&gt;"
		}
		return "Usage : /disable &lt;id&gt;"
	}

	if err := store.SetQueryEnabled(id, enabled); err != nil {
		return "❌ Failed to update query : " + err.Error()
	}

	if enabled {
		return "✅ Query <b>" + strconv.FormatInt(id, 10) + "</b> enabled."
	}
	return "⏸ Query <b>" + strconv.FormatInt(id, 10) + "</b> disabled."
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}