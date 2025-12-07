// Package config 处理应用程序配置
package config

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//go:embed templates/*.yaml
var templates embed.FS

const (
	ConfigFile      = "config.yaml"
	AliyunShareFile = "aliyun_share.yaml"
	PikPakShareFile = "pikpak_share.yaml"
	OneDriveAppFile = "onedrive_app.yaml"
)

// Loader 配置加载器
type Loader struct {
	workDir string
}

// NewLoader 创建配置加载器
func NewLoader(workDir string) *Loader {
	if workDir == "" {
		workDir = "."
	}
	return &Loader{workDir: workDir}
}

// LoadConfig 加载主配置文件
func (l *Loader) LoadConfig() (*Config, error) {
	path := l.filePath(ConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	return &cfg, nil
}

// LoadShareList 加载分享链接列表
func (l *Loader) LoadShareList(filename string) (ShareList, error) {
	path := l.filePath(filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取分享列表失败: %w", err)
	}

	var list ShareList
	if err := yaml.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("解析分享列表失败: %w", err)
	}
	return list, nil
}

// LoadOneDriveList 加载 OneDrive 应用列表
func (l *Loader) LoadOneDriveList(filename string) (OneDriveList, error) {
	path := l.filePath(filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 OneDrive 列表失败: %w", err)
	}

	var list OneDriveList
	if err := yaml.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("解析 OneDrive 列表失败: %w", err)
	}
	return list, nil
}

// SaveConfig 保存配置文件
func (l *Loader) SaveConfig(cfg *Config) error {
	path := l.filePath(ConfigFile)
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// SaveShareList 保存分享链接列表到文件
func (l *Loader) SaveShareList(filename string, list ShareList) error {
	path := l.filePath(filename)
	data, err := yaml.Marshal(list)
	if err != nil {
		return fmt.Errorf("序列化分享列表失败: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// FileExists 检查文件是否存在
func (l *Loader) FileExists(filename string) bool {
	path := l.filePath(filename)
	_, err := os.Stat(path)
	return err == nil
}

// GenerateTemplate 生成模板文件
func (l *Loader) GenerateTemplate(filename string) error {
	templatePath := "templates/" + filename
	data, err := templates.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取模板失败: %w", err)
	}

	path := l.filePath(filename)
	return os.WriteFile(path, data, 0644)
}

// filePath 获取完整文件路径
func (l *Loader) filePath(filename string) string {
	return filepath.Join(l.workDir, filename)
}

// Validate 验证配置有效性
func (cfg *Config) Validate() error {
	if cfg.URL == "" || cfg.URL == "OPENLIST_URL" {
		return fmt.Errorf("URL 未配置")
	}

	hasAuth := cfg.Auth.Username != "" && cfg.Auth.Username != "USERNAME" &&
		cfg.Auth.Password != "" && cfg.Auth.Password != "PASSWORD"
	hasToken := cfg.Token != "" && cfg.Token != "OPENLIST_TOKEN"

	if !hasAuth && !hasToken {
		return fmt.Errorf("token 和用户密码至少需要配置一项")
	}

	if cfg.AliyunShare.Enable {
		if cfg.AliyunShare.RefreshToken == "" || cfg.AliyunShare.RefreshToken == "ALI_YUNPAN_REFRESH_TOKEN" {
			return fmt.Errorf("阿里云盘分享需要配置 refresh_token")
		}
	}

	if cfg.OneDriveApp.Enable {
		if len(cfg.OneDriveApp.Tenants) == 0 {
			return fmt.Errorf("OneDrive 需要配置租户信息")
		}
		for _, t := range cfg.OneDriveApp.Tenants {
			if t.ClientID == "" || t.ClientID == "CLIENT_ID" ||
				t.ClientSecret == "" || t.ClientSecret == "CLIENT_SECRET" ||
				t.TenantID == "" || t.TenantID == "TENANT_ID" {
				return fmt.Errorf("OneDrive 租户配置不完整")
			}
		}
	}

	return nil
}
