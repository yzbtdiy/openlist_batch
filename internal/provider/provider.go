// Package provider 定义存储提供商接口
package provider

import "github.com/yzbtdiy/openlist_batch/internal/model"

// Provider 存储提供商接口
type Provider interface {
	// Name 返回提供商名称
	Name() string
	// Driver 返回 OpenList 驱动名称
	Driver() string
	// BuildRequest 构建存储挂载请求
	BuildRequest(mountPath string, shareURL string) (*model.StorageRequest, error)
}

// UpdateableProvider 支持更新的提供商接口
type UpdateableProvider interface {
	Provider
	// BuildUpdateRequest 构建更新请求
	BuildUpdateRequest(item model.StorageItem, newToken string) (*model.StorageRequest, error)
}
