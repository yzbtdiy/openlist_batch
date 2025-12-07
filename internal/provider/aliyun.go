// Package provider 提供阿里云盘存储支持
package provider

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/yzbtdiy/openlist_batch/internal/model"
)

// AliyunShare 阿里云盘提供商
type AliyunShare struct {
	RefreshToken string
}

// NewAliyunShare 创建阿里云盘提供商
func NewAliyunShare(refreshToken string) *AliyunShare {
	return &AliyunShare{RefreshToken: refreshToken}
}

// Name 返回提供商名称
func (a *AliyunShare) Name() string {
	return "阿里云盘"
}

// Driver 返回 OpenList 驱动名称
func (a *AliyunShare) Driver() string {
	return "AliyundriveShare"
}

// BuildRequest 构建存储挂载请求
func (a *AliyunShare) BuildRequest(mountPath string, shareURL string) (*model.StorageRequest, error) {
	parsed, err := url.Parse(shareURL)
	if err != nil {
		return nil, fmt.Errorf("解析分享链接失败: %w", err)
	}

	// 提取分享密码
	sharePwd := parsed.Query().Get("pwd")

	// 解析路径: /s/shareId/folder/folderId
	pathParts := strings.Split(parsed.Path, "/")
	if len(pathParts) < 5 {
		return nil, fmt.Errorf("无效的阿里云盘分享链接格式")
	}

	shareID := pathParts[2]
	folderID := pathParts[4]

	addition := model.AliyunShareAddition{
		RefreshToken:   a.RefreshToken,
		ShareId:        shareID,
		SharePwd:       sharePwd,
		RootFolderId:   folderID,
		OrderBy:        "",
		OrderDirection: "",
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
		Driver:          a.Driver(),
		Addition:        string(additionJSON),
	}, nil
}

// BuildUpdateRequest 构建更新请求 (更新 RefreshToken)
func (a *AliyunShare) BuildUpdateRequest(item model.StorageItem, newToken string) (*model.StorageRequest, error) {
	var oldAddition model.AliyunShareAddition
	if err := json.Unmarshal([]byte(item.Addition), &oldAddition); err != nil {
		return nil, fmt.Errorf("解析原有附加信息失败: %w", err)
	}

	oldAddition.RefreshToken = newToken
	newAdditionJSON, err := json.Marshal(oldAddition)
	if err != nil {
		return nil, fmt.Errorf("序列化附加信息失败: %w", err)
	}

	return &model.StorageRequest{
		MountPath:       item.MountPath,
		Order:           item.Order,
		Remark:          item.Remark,
		CacheExpiration: item.CacheExpiration,
		WebProxy:        item.WebProxy,
		WebdavPolicy:    item.WebdavPolicy,
		DownProxyUrl:    item.DownProxyURL,
		OrderBy:         item.OrderBy,
		OrderDirection:  item.OrderDirection,
		ExtractFolder:   item.ExtractFolder,
		EnableSign:      false,
		Driver:          item.Driver,
		Addition:        string(newAdditionJSON),
	}, nil
}
