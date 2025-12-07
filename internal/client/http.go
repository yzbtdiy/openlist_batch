// Package client 提供 HTTP 客户端封装
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yzbtdiy/openlist_batch/internal/model"
)

// Client HTTP 客户端
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// NewClient 创建新的 HTTP 客户端
func NewClient(baseURL, token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		token:      token,
	}
}

// SetToken 设置认证 token
func (c *Client) SetToken(token string) {
	c.token = token
}

// Post 发送 POST 请求
func (c *Client) Post(endpoint string, data []byte) (*model.APIResponse, error) {
	url := c.baseURL + endpoint

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result model.APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// Get 发送 GET 请求
func (c *Client) Get(endpoint string) (*model.APIResponse, error) {
	url := c.baseURL + endpoint

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result model.APIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// Close 关闭客户端
func (c *Client) Close() {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}
