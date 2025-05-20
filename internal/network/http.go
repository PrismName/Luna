package network

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
)

// HTTPClient 是HTTP客户端的接口定义
type HTTPClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
	Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error)
	Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HTTPResponse, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error)
	Do(req *http.Request) (*HTTPResponse, error)
}

// HTTPResponse 封装HTTP响应
type HTTPResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Request    *http.Request
}

// HTTPClientConfig 配置HTTP客户端
type HTTPClientConfig struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryInterval  time.Duration
	BaseURL        string
	DefaultHeaders map[string]string
}

// DefaultHTTPClientConfig 返回默认的HTTP客户端配置
func DefaultHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		RetryInterval: 1 * time.Second,
		DefaultHeaders: map[string]string{
			"User-Agent": "Luna/1.0",
		},
	}
}

// Client 实现HTTPClient接口
type Client struct {
	client *http.Client
	config HTTPClientConfig
}

// NewHTTPClient 创建一个新的HTTP客户端
func NewHTTPClient(config HTTPClientConfig) *Client {
	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &Client{
		client: client,
		config: config,
	}
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, urlStr string, headers map[string]string) (*HTTPResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, urlStr, nil, headers)
	if err != nil {
		return nil, err
	}

	return c.doWithRetry(req)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, urlStr string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPost, urlStr, body, headers)
	if err != nil {
		return nil, err
	}

	return c.doWithRetry(req)
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, urlStr string, body interface{}, headers map[string]string) (*HTTPResponse, error) {
	req, err := c.newRequest(ctx, http.MethodPut, urlStr, body, headers)
	if err != nil {
		return nil, err
	}

	return c.doWithRetry(req)
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, urlStr string, headers map[string]string) (*HTTPResponse, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, urlStr, nil, headers)
	if err != nil {
		return nil, err
	}

	return c.doWithRetry(req)
}

// Do 执行HTTP请求
func (c *Client) Do(req *http.Request) (*HTTPResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
		Request:    req,
	}, nil
}

// newRequest 创建一个新的HTTP请求
func (c *Client) newRequest(ctx context.Context, method, urlStr string, body interface{}, headers map[string]string) (*http.Request, error) {
	// 处理基础URL
	if c.config.BaseURL != "" && !strings.HasPrefix(urlStr, "http") {
		urlStr = fmt.Sprintf("%s/%s", strings.TrimRight(c.config.BaseURL, "/"), strings.TrimLeft(urlStr, "/"))
	}

	var bodyReader io.Reader
	if body != nil {
		switch v := body.(type) {
		case string:
			bodyReader = strings.NewReader(v)
		case []byte:
			bodyReader = bytes.NewReader(v)
		case io.Reader:
			bodyReader = v
		default:
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(b)
			// 如果没有指定Content-Type，则默认为JSON
			if headers == nil {
				headers = make(map[string]string)
			}
			if _, ok := headers["Content-Type"]; !ok {
				headers["Content-Type"] = "application/json"
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, bodyReader)
	if err != nil {
		return nil, err
	}

	// 添加默认请求头
	for k, v := range c.config.DefaultHeaders {
		req.Header.Set(k, v)
	}

	// 添加自定义请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// doWithRetry 执行HTTP请求并支持重试
func (c *Client) doWithRetry(req *http.Request) (*HTTPResponse, error) {
	var (
		resp *HTTPResponse
		err  error
		try  int
	)

	for try = 0; try <= c.config.MaxRetries; try++ {
		// 如果不是第一次尝试，则等待重试间隔
		if try > 0 {
			time.Sleep(c.config.RetryInterval)
		}

		// 创建请求的副本，因为原始请求的Body可能已经被消费
		reqCopy := req.Clone(req.Context())
		resp, err = c.Do(reqCopy)

		// 如果请求成功或者是非临时性错误，则不再重试
		if err == nil || !isTemporaryError(err) {
			break
		}
	}

	return resp, err
}

// isTemporaryError 判断错误是否为临时性错误
func isTemporaryError(err error) bool {
	return true
}

// ParseJSON 解析响应体为JSON
func (r *HTTPResponse) ParseJSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
}

// String 返回响应体的字符串表示
func (r *HTTPResponse) String() string {
	return string(r.Body)
}

// IsSuccess 判断响应是否成功
func (r *HTTPResponse) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// BuildURL 构建URL，添加查询参数
func BuildURL(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
