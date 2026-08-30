
package notifier

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Chamroosh98/YT-Crawler/internal/models"
)

type Telegram struct {
	botToken string
	chatID   string
	client   *http.Client
}

func NewTelegram(botToken, chatID string) *Telegram {
	return &Telegram{
		botToken: botToken,
		chatID:   chatID,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (t *Telegram) Enabled() bool {
	return t.botToken != "" && t.chatID != ""
}

// SendVideo sends a video notification to Telegram
func (t *Telegram) SendVideo(v models.Video) error {
	if !t.Enabled() {
		return nil
	}

	text := formatVideoMessage(v)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	data := url.Values{}
	data.Set("chat_id", t.chatID)
	data.Set("text", text)
	data.Set("parse_mode", "HTML")
	data.Set("disable_web_page_preview", "false")

	resp, err := t.client.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status: %s", resp.Status)
	}

	return nil
}

func formatVideoMessage(v models.Video) string {
	var sb strings.Builder

	sb.WriteString("🎬 <b>" + escapeHTML(v.Title) + "</b>\n\n")
	sb.WriteString("👤 Channel: " + escapeHTML(v.Channel) + "\n")
	sb.WriteString("🗓 Published: " + v.PublishedAt.Format("2006-01-02") + "\n")

	if v.Language != "" {
		sb.WriteString("🌐 Language: " + strings.ToUpper(v.Language) + "\n")
	}

	sb.WriteString("\n🔗 " + v.URL)

	return sb.String()
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}


// SendSummary sends a summary of new videos
func (t *Telegram) SendSummary(videos []models.Video) error {
	if !t.Enabled() || len(videos) == 0 {
		return nil
	}

	// Telegram message limit is 4096 characters
	const maxLen = 4000

	var sb strings.Builder
	sb.WriteString("📊 <b>Crawler Summary</b>\n")
	sb.WriteString(fmt.Sprintf("🆕 New videos: <b>%d</b>\n\n", len(videos)))

	for i, v := range videos {
		line := fmt.Sprintf("%d. <b>%s</b>\n👤 %s\n🔗 %s\n\n",
			i+1,
			escapeHTML(v.Title),
			escapeHTML(v.Channel),
			v.URL,
		)

		// Prevent exceeding Telegram limit
		if sb.Len()+len(line) > maxLen {
			sb.WriteString(fmt.Sprintf("... and %d more videos\n", len(videos)-i))
			break
		}
		sb.WriteString(line)
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	data := url.Values{}
	data.Set("chat_id", t.chatID)
	data.Set("text", sb.String())
	data.Set("parse_mode", "HTML")
	data.Set("disable_web_page_preview", "true")

	resp, err := t.client.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("failed to send summary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned status: %s", resp.Status)
	}

	return nil
}