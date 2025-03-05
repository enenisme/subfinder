package scrape

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/corpix/uarand"
)

// NewSession
func NewSession(domain string, timeout int) (*Session, error) {
	transport := &http.Transport{
		MaxIdleConns:        100, // 最大空闲连接数
		MaxIdleConnsPerHost: 100, // 每个主机的最大空闲连接数
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableCompression: true, // 禁用压缩
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	extractor, err := NewSubdomainExtractor(domain)

	return &Session{
		HttpClient: client,
		Extractor:  extractor,
	}, err

}

// Get 用于发起 GET 请求
func (s *Session) Get(ctx context.Context, url, cookies string, headers map[string]string) (*http.Response, error) {
	return s.HTTPRequest(ctx, http.MethodGet, url, cookies, headers, nil)
}

// Post 用于发起 POST 请求
func (s *Session) Post(ctx context.Context, url, cookies string, headers map[string]string, body io.Reader) (*http.Response, error) {
	return s.HTTPRequest(ctx, http.MethodPost, url, cookies, headers, body)
}

// HTTPRequest 用于发起 HTTP 请求
func (s *Session) HTTPRequest(ctx context.Context, method, url, cookies string,
	headers map[string]string, body io.Reader) (*http.Response, error) {
	// 创建 http 请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %v", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Connection", "close")

	// 如果提供了 Cookie，设置 Cookie 请求头
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}

	// 设置自定义的请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return s.HttpClient.Do(req)
}

// Close 用于关闭 HTTP 客户端
func (s *Session) Close() {
	s.HttpClient.CloseIdleConnections()
}
