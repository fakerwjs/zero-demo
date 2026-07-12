package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SMSConfig 短信网关配置（HTTP 形式，兼容阿里云/腾讯云等自建代理）。
type SMSConfig struct {
	Endpoint  string // 短信网关 HTTP 接口
	AccessKey string
	SignName  string // 短信签名
	Template  string // 模板ID
}

// SMSSender 通过 HTTP 调用短信网关真实发送。
type SMSSender struct {
	cfg    SMSConfig
	client *http.Client
}

func NewSMSSender(cfg SMSConfig) *SMSSender {
	return &SMSSender{cfg: cfg, client: &http.Client{Timeout: 5 * time.Second}}
}

func (s *SMSSender) Send(ctx context.Context, to, subject, body string) error {
	payload, _ := json.Marshal(map[string]string{
		"accessKey": s.cfg.AccessKey,
		"signName":  s.cfg.SignName,
		"template":  s.cfg.Template,
		"mobile":    to,
		"content":   body,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sms gateway status=%d body=%s", resp.StatusCode, string(b))
	}
	return nil
}
