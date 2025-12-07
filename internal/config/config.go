// Package config 处理应用程序配置
package config

// Config 主配置结构
type Config struct {
	URL      string   `yaml:"url"`
	Auth     Auth     `yaml:"auth"`
	Token    string   `yaml:"token"`
	Aliyun   Aliyun   `yaml:"aliyun"`
	PikPak   PikPak   `yaml:"pikpak"`
	OneDrive OneDrive `yaml:"onedrive_app"`
}

// Auth OpenList 登录认证信息
type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Aliyun 阿里云盘配置
type Aliyun struct {
	Enable       bool   `yaml:"enable"`
	RefreshToken string `yaml:"refresh_token"`
}

// PikPak PikPak 配置
type PikPak struct {
	Enable                bool   `yaml:"enable"`
	UseTranscodingAddress bool   `yaml:"use_transcoding_address"`
	Username              string `yaml:"username"`
	Password              string `yaml:"password"`
}

// OneDrive OneDrive APP 配置
type OneDrive struct {
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
