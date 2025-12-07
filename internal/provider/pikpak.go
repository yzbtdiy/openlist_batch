// Package provider 提供 PikPak 存储支持
package provider

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/yzbtdiy/openlist_batch/internal/model"
)

// PikPak PikPak 提供商
type PikPak struct {
	Username              string
	Password              string
	UseTranscodingAddress bool
}

// NewPikPak 创建 PikPak 提供商
func NewPikPak(username, password string, useTranscoding bool) *PikPak {
	return &PikPak{
		Username:              username,
		Password:              password,
		UseTranscodingAddress: useTranscoding,
	}
}

// Name 返回提供商名称
func (p *PikPak) Name() string {
	return "PikPak"
}

// Driver 返回 OpenList 驱动名称
func (p *PikPak) Driver() string {
	return "PikPakShare"
}

// BuildRequest 构建存储挂载请求
func (p *PikPak) BuildRequest(mountPath string, shareURL string) (*model.StorageRequest, error) {
	parsed, err := url.Parse(shareURL)
	if err != nil {
		return nil, fmt.Errorf("解析分享链接失败: %w", err)
	}

	// 提取分享密码
	sharePwd := parsed.Query().Get("pwd")

	// 解析路径: /s/shareId 或 /s/shareId/folderId
	pathParts := strings.Split(parsed.Path, "/")
	if len(pathParts) < 3 {
		return nil, fmt.Errorf("无效的 PikPak 分享链接格式")
	}

	shareID := pathParts[2]
	folderID := ""
	if len(pathParts) >= 4 {
		folderID = pathParts[3]
	}

	addition := model.PikPakAddition{
		RootFolderId:          folderID,
		Username:              p.Username,
		Password:              p.Password,
		ShareId:               shareID,
		SharePwd:              sharePwd,
		OrderBy:               "",
		OrderDirection:        "",
		Platform:              "android",
		UseTranscodingAddress: p.UseTranscodingAddress,
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
		Driver:          p.Driver(),
		Addition:        string(additionJSON),
	}, nil
}
