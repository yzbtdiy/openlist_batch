// Package config 处理应用程序配置
package config

// Config 主配置结构
type Config struct {
	URL         string      `yaml:"url"`
	Auth        Auth        `yaml:"auth"`
	Token       string      `yaml:"token"`
	AliyunShare AliyunShare `yaml:"aliyun_share"`
	PikPakShare PikPakShare `yaml:"pikpak_share"`
	OneDriveApp OneDriveApp `yaml:"onedrive_app"`
}

// Auth OpenList 登录认证信息
type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// 阿里云盘配置
type AliyunShare struct {
	Enable       bool   `yaml:"enable"`
	RefreshToken string `yaml:"refresh_token"`
}

// PikPak 配置
type PikPak struct {
	Enable       bool   `yaml:"enable"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Platform     string `yaml:"platform"`
	RefreshToken string `yaml:"refresh_token"`
}

type PikPakShare struct {
	Enable                bool   `yaml:"enable"`
	UseTranscodingAddress bool   `yaml:"use_transcoding_address"`
	Platform              string `yaml:"platform"`
}

// OneDrive APP 配置
type OneDriveApp struct {
	Enable  bool     `yaml:"enable"`
	Region  string   `yaml:"region"`
	Tenants []Tenant `yaml:"tenants"`
}

// Tenant OneDrive 租户信息
type Tenant struct {
	ID           int    `yaml:"id"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	TenantID     string `yaml:"tenant_id"`
}

// ShareList 分享链接列表 (用于 aliyun 和 pikpak)
type ShareList map[string]map[string]string

// OneDriveList OneDrive 应用列表
type OneDriveList map[string][]string
