# OpenList Batch

OpenList 批量存储管理工具，目前支持批量添加阿里云盘分享链接、PikPak分享链接、OneDriveApp挂载。

添加存储API所用API在OpenList v4.1.8抓包测试。

有需求请提 Issues, 最好说明使用的版本。

OpenList Github仓库地址: https://github.com/OpenListTeam/OpenList

## 功能特性

- 🚀 批量添加阿里云盘分享链接
- 🚀 批量添加 PikPak 分享链接
- 🚀 批量添加 OneDrive APP
- 🔄 自动获取并保存 Token
- 🗑️ 批量删除存储（支持删除禁用/全部）
- 🔧 批量更新阿里云盘 RefreshToken

## 项目结构

```
openlist_batch/
├── cmd/
│   └── openlist_batch/
│       └── main.go           # 程序入口
├── internal/
│   ├── client/
│   │   └── http.go           # HTTP 客户端封装
│   ├── config/
│   │   ├── config.go         # 配置结构定义
│   │   ├── loader.go         # 配置加载器
│   │   └── templates/        # 配置模板
│   │       ├── config.yaml
│   │       ├── aliyun_share.yaml
│   │       ├── pikpak_share.yaml
│   │       └── onedrive_app.yaml
│   ├── model/
│   │   ├── request.go        # 请求模型
│   │   └── response.go       # 响应模型
│   ├── provider/
│   │   ├── provider.go       # 提供商接口
│   │   ├── aliyun.go         # 阿里云盘
│   │   ├── pikpak.go         # PikPak
│   │   └── onedrive.go       # OneDrive
│   └── service/
│       └── batch.go          # 批处理服务
├── go.mod
└── README.md
```

## 安装

### 从源码编译

```bash
git clone https://github.com/yzbtdiy/openlist_batch.git
cd openlist_batch
go mod tidy
go build -o openlist_batch ./cmd/openlist_batch
```

### 使用 go install

```bash
go install github.com/yzbtdiy/openlist_batch/cmd/openlist_batch@VERSION
```

## 使用方法

### 1. 初始化配置

首次运行会自动生成配置模板：

```bash
./openlist_batch
```

### 2. 编辑配置文件

编辑 `config.yaml`，填写 OpenList 地址和认证信息：

```yaml
url: http://localhost:5244  # OpenList 地址
auth:
  username: admin           # 用户名
  password: password        # 密码
token: ""                   # Token（可选，会自动获取）

aliyun:
  enable: true              # 是否启用阿里云盘
  refresh_token: xxx        # 阿里云盘 RefreshToken

pikpak:
  enable: false             # 使用启用 PikPak
  use_transcoding_address: true
  username: xxx
  password: xxx

onedrive_app:
  enable: false             # 是否启用 OneDrive
  region: global
  tenants:
    - id: 1
      client_id: xxx
      client_secret: xxx
      tenant_id: xxx
```

### 3. 添加分享链接

根据启用的存储类型，编辑对应的分享文件：

**aliyun_share.yaml** (阿里云盘):
```yaml
电视剧:
  西游记86版: https://www.aliyundrive.com/s/MmMR3zaoXLf/folder/61d259418d27bae8656f47aca23ee03b40275432

电影:
  新海诚&宫崎骏合集: https://www.aliyundrive.com/s/FzcMCgG8YwC/folder/61ffb364be026f8c1b764182922eaeb2d3950ef4
  林正英合集: https://www.aliyundrive.com/s/PrcaqZ2XPxM/folder/621c950a633c7c7ab8de4db1a86a1232dea530d1
```

**pikpak_share.yaml** (PikPak):
```yaml
电影:
  太空之城: https://mypikpak.com/s/VNP2_7OhUCdC2aI3JSSnD--eo1
  阿飞正传: https://mypikpak.com/s/VNP2d8tHvt4TVPKPacCUYRaXo1/VNP2G0YUcYmtVw025fNVqgDdo1
```

**onedrive_app.yaml** (OneDrive):
```yaml
个人网盘:
  工作文件: 1:user@example.com:/Work
  游戏娱乐: 1:user@xxx.onmicrosoft.com:/Games
```

### 4. 运行

```bash
# 批量添加
./openlist_batch

# 删除禁用的存储
./openlist_batch -delete dis

# 删除所有存储（慎用）
./openlist_batch -delete all

# 更新阿里云盘 RefreshToken
./openlist_batch -update ali
```

## 分享链接格式

### 阿里云盘
```
https://www.alipan.com/s/shareId/folder/folderId?pwd=提取码
```
注意：链接必须包含 `folder/folderId` 部分

### PikPak
```
https://mypikpak.com/s/shareId
https://mypikpak.com/s/shareId/folderId?pwd=提取码
```

### OneDrive
```
tid:email:path
```
- `tid`: 租户ID索引（对应 config.yaml 中 tenants 的序号，从1开始）
- `email`: 账户邮箱
- `path`: 文件夹路径（可选，默认为 /）

## 注意事项

- OpenList URL 结尾不要加 `/`
- 首次运行需要配置用户名密码或有效 Token
- 此工具仅用于批量挂载，遇到问题请参考 [OpenList 官方文档](https://doc.oplist.org/) 或 [Github Issues](https://github.com/OpenListTeam/OpenList/issues)

## 许可证

MIT License
