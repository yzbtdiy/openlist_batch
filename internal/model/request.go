// Package model 定义数据模型
package model

// AuthRequest 登录认证请求
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// StorageRequest 存储挂载请求
type StorageRequest struct {
	MountPath       string `json:"mount_path"`
	Order           int    `json:"order"`
	Remark          string `json:"remark"`
	CacheExpiration int    `json:"cache_expiration"`
	WebProxy        bool   `json:"web_proxy"`
	WebdavPolicy    string `json:"webdav_policy"`
	DownProxyUrl    string `json:"down_proxy_url"`
	OrderBy         string `json:"order_by"`
	OrderDirection  string `json:"order_direction"`
	ExtractFolder   string `json:"extract_folder"`
	EnableSign      bool   `json:"enable_sign"`
	Driver          string `json:"driver"`
	Addition        string `json:"addition"`
}

// AliyunAddition 阿里云盘挂载附加信息
type AliyunAddition struct {
	RefreshToken   string `json:"refresh_token"`
	ShareId        string `json:"share_id"`
	SharePwd       string `json:"share_pwd"`
	RootFolderId   string `json:"root_folder_id"`
	OrderBy        string `json:"order_by"`
	OrderDirection string `json:"order_direction"`
}

// PikPakAddition PikPak 挂载附加信息
type PikPakAddition struct {
	RootFolderId          string `json:"root_folder_id"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	ShareId               string `json:"share_id"`
	SharePwd              string `json:"share_pwd"`
	OrderBy               string `json:"order_by"`
	OrderDirection        string `json:"order_direction"`
	Platform              string `json:"platform"`
	UseTranscodingAddress bool   `json:"use_transcoding_address"`
}

// OneDriveAddition OneDrive APP 挂载附加信息
type OneDriveAddition struct {
	RootFolderPath string `json:"root_folder_path"`
	Region         string `json:"region"`
	ClientId       string `json:"client_id"`
	ClientSecret   string `json:"client_secret"`
	TenantId       string `json:"tenant_id"`
	Email          string `json:"email"`
	ChunkSize      int    `json:"chunk_size"`
}
