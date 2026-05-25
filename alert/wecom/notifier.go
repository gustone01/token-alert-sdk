package wecom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gustone01/token-alert-sdk/alert/core"
)

type Notifier struct {
	webhookURL    string
	mentionedList []string
	client        *http.Client
}

func New(webhookKeyOrURL string, mentionedList []string) *Notifier {
	return &Notifier{
		webhookURL:    BuildWebhookURL(webhookKeyOrURL),
		mentionedList: mentionedList,
		client:        http.DefaultClient,
	}
}

func BuildWebhookURL(webhookKeyOrURL string) string {
	v := strings.TrimSpace(webhookKeyOrURL)
	if v == "" {
		return ""
	}
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return v
	}
	return "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + url.QueryEscape(v)
}

func (n *Notifier) Send(ctx context.Context, event core.Event) error {
	if n == nil || n.webhookURL == "" {
		return fmt.Errorf("wecom webhook URL is empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	content := renderContent(event)
	body, err := json.Marshal(map[string]any{
		"msgtype": "text",
		"text": map[string]any{
			"content":        content,
			"mentioned_list": n.mentionedList,
		},
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := n.client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("wecom webhook http status %d body %s", resp.StatusCode, string(respBody))
	}
	var jr struct {
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(respBody, &jr); err == nil && jr.ErrCode != 0 {
		return fmt.Errorf("wecom errcode=%d errmsg=%s", jr.ErrCode, jr.ErrMsg)
	}
	return nil
}

func renderContent(event core.Event) string {
	when := event.OccurredAt
	if when.IsZero() {
		when = time.Now()
	}
	lines := []string{
		"【媒体Token失效告警】",
		fmt.Sprintf("来源：%s", event.Service),
		fmt.Sprintf("平台：%s", event.Platform.DisplayName()),
		fmt.Sprintf("接口：%s", event.APIPath),
		fmt.Sprintf("错误码：%d", event.Code),
		fmt.Sprintf("错误信息：%s", event.Message),
	}
	if event.Enrich.AccountID != "" || event.Enrich.Name != "" {
		lines = append(lines, fmt.Sprintf("账号：%s / %s", event.Enrich.AccountID, event.Enrich.Name))
	}
	lines = append(lines, fmt.Sprintf("时间：%s", when.Format("2006-01-02 15:04:05")))
	return strings.Join(lines, "\n")
}
