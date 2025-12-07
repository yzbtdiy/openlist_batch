// Package provider 提供 OneDrive 存储支持
package provider

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/yzbtdiy/openlist_batch/internal/config"
	"github.com/yzbtdiy/openlist_batch/internal/model"
)

// OneDrive OneDrive APP 提供商
type OneDrive struct {
	Region  string
	Tenants []config.Tenant
}

// NewOneDrive 创建 OneDrive 提供商
func NewOneDrive(region string, tenants []config.Tenant) *OneDrive {
	return &OneDrive{
		Region:  region,
		Tenants: tenants,
	}
}

// Name 返回提供商名称
func (o *OneDrive) Name() string {
	return "OneDrive APP"
}

// Driver 返回 OpenList 驱动名称
func (o *OneDrive) Driver() string {
	return "OnedriveAPP"
}

// BuildRequest 构建存储挂载请求
// emailInfo 格式: "tid:email:path" 或 "tid:email"
func (o *OneDrive) BuildRequest(mountPath string, emailInfo string) (*model.StorageRequest, error) {
	parts := strings.Split(emailInfo, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return nil, fmt.Errorf("无效的 OneDrive 配置格式: %s, 应为 tid:email[:path]", emailInfo)
	}

	// 解析租户 ID 索引
	tid, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("无效的租户 ID: %s", parts[0])
	}

	if tid < 1 || tid > len(o.Tenants) {
		return nil, fmt.Errorf("租户 ID 超出范围: %d, 有效范围: 1-%d", tid, len(o.Tenants))
	}

	email := parts[1]
	folderPath := "/"
	if len(parts) == 3 {
		folderPath = parts[2]
	}

	tenant := o.Tenants[tid-1]
	addition := model.OneDriveAddition{
		RootFolderPath: folderPath,
		Region:         o.Region,
		ClientId:       tenant.ClientID,
		ClientSecret:   tenant.ClientSecret,
		TenantId:       tenant.TenantID,
		Email:          email,
		ChunkSize:      5,
	}

	additionJSON, err := json.Marshal(addition)
	if err != nil {
		return nil, fmt.Errorf("序列化附加信息失败: %w", err)
	}

	return &model.StorageRequest{
		MountPath:       mountPath,
		Order:           0,
		Remark:          "",
		CacheExpiration: 30,
		WebProxy:        false,
		WebdavPolicy:    "302_redirect",
		DownProxyUrl:    "",
		OrderBy:         "",
		OrderDirection:  "",
		ExtractFolder:   "",
		EnableSign:      false,
		Driver:          o.Driver(),
		Addition:        string(additionJSON),
	}, nil
}
