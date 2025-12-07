// Package service 提供核心业务逻辑
package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/yzbtdiy/openlist_batch/internal/client"
	"github.com/yzbtdiy/openlist_batch/internal/config"
	"github.com/yzbtdiy/openlist_batch/internal/model"
	"github.com/yzbtdiy/openlist_batch/internal/provider"
)

// API 端点常量
const (
	LoginEndpoint         = "/api/auth/login"
	StorageListEndpoint   = "/api/admin/storage/list"
	StorageCreateEndpoint = "/api/admin/storage/create"
	StorageDeleteEndpoint = "/api/admin/storage/delete"
	StorageUpdateEndpoint = "/api/admin/storage/update"
)

// BatchService 批处理服务
type BatchService struct {
	cfg    *config.Config
	client *client.Client
	loader *config.Loader
}

// NewBatchService 创建批处理服务
func NewBatchService(cfg *config.Config, loader *config.Loader) *BatchService {
	return &BatchService{
		cfg:    cfg,
		client: client.NewClient(cfg.URL, cfg.Token),
		loader: loader,
	}
}

// ValidateToken 验证 token 是否有效
func (s *BatchService) ValidateToken() bool {
	resp, err := s.client.Get(StorageListEndpoint)
	if err != nil {
		return false
	}
	return resp.Code == 200
}

// RefreshToken 刷新 token
func (s *BatchService) RefreshToken() error {
	authReq := model.AuthRequest{
		Username: s.cfg.Auth.Username,
		Password: s.cfg.Auth.Password,
	}

	data, err := json.Marshal(authReq)
	if err != nil {
		return fmt.Errorf("序列化登录请求失败: %w", err)
	}

	resp, err := s.client.Post(LoginEndpoint, data)
	if err != nil {
		return fmt.Errorf("登录请求失败: %w", err)
	}

	if resp.Code != 200 {
		return fmt.Errorf("登录失败: %s", resp.Message)
	}

	// 解析 token
	tokenData, err := json.Marshal(resp.Data)
	if err != nil {
		return fmt.Errorf("解析登录响应失败: %w", err)
	}

	var authResp model.AuthResponse
	if err := json.Unmarshal(tokenData, &authResp); err != nil {
		return fmt.Errorf("解析 token 失败: %w", err)
	}

	// 更新配置并保存
	s.cfg.Token = authResp.Token
	s.client.SetToken(authResp.Token)

	if err := s.loader.SaveConfig(s.cfg); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	log.Println("Token 已更新")
	return nil
}

// GetStorageList 获取存储列表
func (s *BatchService) GetStorageList() (*model.StorageListResponse, error) {
	resp, err := s.client.Get(StorageListEndpoint)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("序列化存储列表失败: %w", err)
	}

	var list model.StorageListResponse
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("解析存储列表失败: %w", err)
	}

	return &list, nil
}

// AddStorage 添加单个存储
func (s *BatchService) AddStorage(req *model.StorageRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	resp, err := s.client.Post(StorageCreateEndpoint, data)
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return fmt.Errorf("%s", resp.Message)
	}

	return nil
}

// DeleteStorage 删除存储
func (s *BatchService) DeleteStorage(id int) error {
	endpoint := fmt.Sprintf("%s?id=%d", StorageDeleteEndpoint, id)
	resp, err := s.client.Post(endpoint, []byte{})
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return fmt.Errorf("%s", resp.Message)
	}

	return nil
}

// UpdateStorage 更新存储
func (s *BatchService) UpdateStorage(req *model.StorageRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	resp, err := s.client.Post(StorageUpdateEndpoint, data)
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return fmt.Errorf("%s", resp.Message)
	}

	return nil
}

// BatchAddShares 批量添加分享链接
func (s *BatchService) BatchAddShares(p provider.Provider, shares config.ShareList) {
	var wg sync.WaitGroup

	for category, shareMap := range shares {
		for name, url := range shareMap {
			wg.Add(1)
			go func(category, name, url string) {
				defer wg.Done()

				mountPath := "/" + category + "/" + name
				req, err := p.BuildRequest(mountPath, url)
				if err != nil {
					log.Printf("[%s] %s/%s 构建请求失败: %v", p.Name(), category, name, err)
					return
				}

				if err := s.AddStorage(req); err != nil {
					log.Printf("[%s] %s/%s 添加失败: %v", p.Name(), category, name, err)
					return
				}

				log.Printf("[%s] %s/%s 添加成功", p.Name(), category, name)
			}(category, name, url)
		}
	}

	wg.Wait()
}

// BatchAddOneDriveApp 批量添加 OneDrive 应用
func (s *BatchService) BatchAddOneDriveApp(p *provider.OneDriveApp, appList config.ShareList) {
	var wg sync.WaitGroup

	for category, appMap := range appList {
		for name, emailInfo := range appMap {
			wg.Add(1)
			go func(category, name, emailInfo string) {
				defer wg.Done()

				mountPath := "/" + category + "/" + name
				req, err := p.BuildRequest(mountPath, emailInfo)
				if err != nil {
					log.Printf("[OneDrive] %s/%s 构建请求失败: %v", category, name, err)
					return
				}

				if err := s.AddStorage(req); err != nil {
					log.Printf("[OneDrive] %s/%s 添加失败: %v", category, name, err)
					return
				}

				log.Printf("[OneDrive] %s/%s 添加成功", category, name)
			}(category, name, emailInfo)
		}
	}
	wg.Wait()
}

// DeleteDisabledStorages 删除禁用的存储
func (s *BatchService) DeleteDisabledStorages() error {
	list, err := s.GetStorageList()
	if err != nil {
		return err
	}

	for _, item := range list.Content {
		if item.Disabled {
			if err := s.DeleteStorage(item.Id); err != nil {
				log.Printf("删除存储 %d (%s) 失败: %v", item.Id, item.MountPath, err)
			} else {
				log.Printf("已删除存储 %d (%s)", item.Id, item.MountPath)
			}
		}
	}

	return nil
}

// DeleteAllStorages 删除所有存储
func (s *BatchService) DeleteAllStorages() error {
	list, err := s.GetStorageList()
	if err != nil {
		return err
	}

	for _, item := range list.Content {
		if err := s.DeleteStorage(item.Id); err != nil {
			log.Printf("删除存储 %d (%s) 失败: %v", item.Id, item.MountPath, err)
		} else {
			log.Printf("已删除存储 %d (%s)", item.Id, item.MountPath)
		}
	}

	return nil
}

// UpdateAliyunRefreshToken 更新阿里云盘 RefreshToken
func (s *BatchService) UpdateAliyunRefreshToken(newToken string) error {
	list, err := s.GetStorageList()
	if err != nil {
		return err
	}

	aliyunProvider := provider.NewAliyunShare(newToken)

	for _, item := range list.Content {
		if item.Driver != aliyunProvider.Driver() {
			continue
		}

		req, err := aliyunProvider.BuildUpdateRequest(item, newToken)
		if err != nil {
			log.Printf("构建更新请求失败 (%s): %v", item.MountPath, err)
			continue
		}

		// 需要设置 ID 用于更新
		data, _ := json.Marshal(req)
		var reqMap map[string]any
		json.Unmarshal(data, &reqMap)
		reqMap["id"] = item.Id
		reqMap["status"] = "work"

		updateData, _ := json.Marshal(reqMap)
		resp, err := s.client.Post(StorageUpdateEndpoint, updateData)
		if err != nil {
			log.Printf("更新存储失败 (%s): %v", item.MountPath, err)
			continue
		}

		if resp.Code == 200 {
			log.Printf("已更新 %s", item.MountPath)
		} else {
			log.Printf("更新失败 (%s): %s", item.MountPath, resp.Message)
		}
	}

	return nil
}

// DeleteStorageByID 根据 ID 删除存储
func (s *BatchService) DeleteStorageByID(ids []string) error {
	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("无效的存储 ID: %s", idStr)
			continue
		}

		if err := s.DeleteStorage(id); err != nil {
			log.Printf("删除存储 %d 失败: %v", id, err)
		} else {
			log.Printf("已删除存储 %d", id)
		}
	}

	return nil
}

// ExportPikPakShare 导出 PikPakShare 存储到 ShareList 格式
func (s *BatchService) ExportPikPakShare() (config.ShareList, error) {
	list, err := s.GetStorageList()
	if err != nil {
		return nil, fmt.Errorf("获取存储列表失败: %w", err)
	}

	result := make(config.ShareList)

	for _, item := range list.Content {
		if item.Driver != "PikPakShare" {
			continue
		}

		// 解析 mount_path: /分类/资源名
		parts := splitMountPath(item.MountPath)
		if len(parts) < 2 {
			log.Printf("跳过无效的挂载路径: %s", item.MountPath)
			continue
		}

		category := parts[0]
		name := parts[1]

		// 解析 addition 获取分享信息
		var addition model.PikPakShareAddition
		if err := json.Unmarshal([]byte(item.Addition), &addition); err != nil {
			log.Printf("解析存储附加信息失败 (%s): %v", item.MountPath, err)
			continue
		}

		// 构建分享链接
		shareURL := buildPikPakShareURL(addition)

		// 添加到结果
		if result[category] == nil {
			result[category] = make(map[string]string)
		}
		result[category][name] = shareURL
	}

	return result, nil
}

// splitMountPath 分割挂载路径
func splitMountPath(mountPath string) []string {
	// 去除开头的 /
	if len(mountPath) > 0 && mountPath[0] == '/' {
		mountPath = mountPath[1:]
	}

	parts := make([]string, 0)
	for _, p := range splitPath(mountPath) {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

// splitPath 按 / 分割路径
func splitPath(path string) []string {
	result := make([]string, 0)
	start := 0
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if i > start {
				result = append(result, path[start:i])
			}
			start = i + 1
		}
	}
	if start < len(path) {
		result = append(result, path[start:])
	}
	return result
}

// buildPikPakShareURL 根据 addition 构建 PikPak 分享链接
func buildPikPakShareURL(addition model.PikPakShareAddition) string {
	baseURL := "https://mypikpak.com/s/" + addition.ShareId

	if addition.RootFolderId != "" {
		baseURL += "/" + addition.RootFolderId
	}

	if addition.SharePwd != "" {
		baseURL += "?pwd=" + addition.SharePwd
	}

	return baseURL
}

// Close 关闭服务
func (s *BatchService) Close() {
	s.client.Close()
}
