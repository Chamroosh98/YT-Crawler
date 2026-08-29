
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