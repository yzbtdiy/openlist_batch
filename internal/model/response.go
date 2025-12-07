// Package model 定义数据模型
package model

import "time"

// APIResponse 通用 API 响应
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// AuthResponse 登录响应数据
type AuthResponse struct {
	Token string `json:"token"`
}

// StorageListResponse 存储列表响应
type StorageListResponse struct {
	Content []StorageItem `json:"content"`
	Total   int           `json:"total"`
}

// StorageItem 存储项详细信息
type StorageItem struct {
	Id               int       `json:"id"`
	MountPath        string    `json:"mount_path"`
	Order            int       `json:"order"`
	Driver           string    `json:"driver"`
	CacheExpiration  int       `json:"cache_expiration"`
	Status           string    `json:"status"`
	Addition         string    `json:"addition"`
	Remark           string    `json:"remark"`
	Modified         time.Time `json:"modified"`
	Disabled         bool      `json:"disabled"`
	DisableIndex     bool      `json:"disable_index"`
	EnableSign       bool      `json:"enable_sign"`
	OrderBy          string    `json:"order_by"`
	OrderDirection   string    `json:"order_direction"`
	ExtractFolder    string    `json:"extract_folder"`
	WebProxy         bool      `json:"web_proxy"`
	WebdavPolicy     string    `json:"webdav_policy"`
	ProxyRange       bool      `json:"proxy_range"`
	DownProxyURL     string    `json:"down_proxy_url"`
	DisableProxySign bool      `json:"disable_proxy_sign"`
}
