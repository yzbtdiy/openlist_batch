// OpenList Batch - OpenList 批量存储管理工具
//
// 支持批量添加阿里云盘、PikPak、OneDrive 分享链接到 OpenList
package main

import (
	"flag"
	"log"
	"os"

	"github.com/yzbtdiy/openlist_batch/internal/config"
	"github.com/yzbtdiy/openlist_batch/internal/provider"
	"github.com/yzbtdiy/openlist_batch/internal/service"
)

var (
	deleteFlag = flag.String("delete", "", `删除存储:
  dis    删除已禁用存储
  all    删除所有存储(慎用)`)

	updateFlag = flag.String("update", "", `更新存储:
  ali    更新阿里云盘 refresh_token`)
)

func main() {
	flag.Parse()

	loader := config.NewLoader(".")

	// 检查并生成配置文件
	if !loader.FileExists(config.ConfigFile) {
		log.Println("配置文件不存在，正在生成...")
		if err := loader.GenerateTemplate(config.ConfigFile); err != nil {
			log.Fatalf("生成配置文件失败: %v", err)
		}
		log.Println("已生成 config.yaml，请配置后重新运行")
		return
	}

	// 加载配置
	cfg, err := loader.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 检查分享文件
	if cfg.Aliyun.Enable && !loader.FileExists(config.AliyunFile) {
		log.Println("阿里云盘分享文件不存在，正在生成...")
		if err := loader.GenerateTemplate(config.AliyunFile); err != nil {
			log.Fatalf("生成分享文件失败: %v", err)
		}
		log.Println("已生成 aliyun.yaml，请添加分享链接后重新运行")
		return
	}

	if cfg.PikPak.Enable && !loader.FileExists(config.PikPakFile) {
		log.Println("PikPak 分享文件不存在，正在生成...")
		if err := loader.GenerateTemplate(config.PikPakFile); err != nil {
			log.Fatalf("生成分享文件失败: %v", err)
		}
		log.Println("已生成 pikpak.yaml，请添加分享链接后重新运行")
		return
	}

	if cfg.OneDrive.Enable && !loader.FileExists(config.OneDriveFile) {
		log.Println("OneDrive 配置文件不存在，正在生成...")
		if err := loader.GenerateTemplate(config.OneDriveFile); err != nil {
			log.Fatalf("生成配置文件失败: %v", err)
		}
		log.Println("已生成 onedrive.yaml，请配置后重新运行")
		return
	}

	// 创建批处理服务
	svc := service.NewBatchService(cfg, loader)
	defer svc.Close()

	// 验证或刷新 Token
	if !svc.ValidateToken() {
		log.Println("Token 无效，正在刷新...")
		if err := svc.RefreshToken(); err != nil {
			log.Fatalf("刷新 Token 失败: %v", err)
		}
		log.Println("请重新运行程序")
		return
	}

	// 处理删除命令
	if *deleteFlag != "" {
		handleDelete(svc, *deleteFlag)
		return
	}

	// 处理更新命令
	if *updateFlag != "" {
		handleUpdate(svc, cfg, *updateFlag)
		return
	}

	// 批量添加存储
	addStorages(svc, cfg, loader)
}

func handleDelete(svc *service.BatchService, mode string) {
	switch mode {
	case "dis":
		log.Println("正在删除禁用的存储...")
		if err := svc.DeleteDisabledStorages(); err != nil {
			log.Fatalf("删除失败: %v", err)
		}
	case "all":
		log.Println("警告: 正在删除所有存储!")
		if err := svc.DeleteAllStorages(); err != nil {
			log.Fatalf("删除失败: %v", err)
		}
	default:
		log.Printf("未知的删除模式: %s", mode)
		os.Exit(1)
	}
}

func handleUpdate(svc *service.BatchService, cfg *config.Config, mode string) {
	switch mode {
	case "ali":
		if !cfg.Aliyun.Enable {
			log.Fatal("阿里云盘未启用")
		}
		log.Println("正在更新阿里云盘 RefreshToken...")
		if err := svc.UpdateAliyunRefreshToken(cfg.Aliyun.RefreshToken); err != nil {
			log.Fatalf("更新失败: %v", err)
		}
	default:
		log.Printf("未知的更新模式: %s", mode)
		os.Exit(1)
	}
}

func addStorages(svc *service.BatchService, cfg *config.Config, loader *config.Loader) {
	// 添加阿里云盘分享
	if cfg.Aliyun.Enable {
		shares, err := loader.LoadShareList(config.AliyunFile)
		if err != nil {
			log.Printf("加载阿里云盘分享失败: %v", err)
		} else {
			log.Println("正在添加阿里云盘分享...")
			aliyun := provider.NewAliyun(cfg.Aliyun.RefreshToken)
			svc.BatchAddShares(aliyun, shares)
		}
	}

	// 添加 PikPak 分享
	if cfg.PikPak.Enable {
		shares, err := loader.LoadShareList(config.PikPakFile)
		if err != nil {
			log.Printf("加载 PikPak 分享失败: %v", err)
		} else {
			log.Println("正在添加 PikPak 分享...")
			pikpak := provider.NewPikPak(cfg.PikPak.Username, cfg.PikPak.Password, cfg.PikPak.UseTranscodingAddress)
			svc.BatchAddShares(pikpak, shares)
		}
	}

	// 添加 OneDrive 应用
	if cfg.OneDrive.Enable {
		shares, err := loader.LoadShareList(config.OneDriveFile)
		if err != nil {
			log.Printf("加载 OneDrive 配置失败: %v", err)
		} else {
			log.Println("正在添加 OneDrive 应用...")
			onedrive := provider.NewOneDrive(cfg.OneDrive.Region, cfg.OneDrive.Tenants)
			svc.BatchAddOneDrive(onedrive, shares)
		}
	}

	log.Println("批量操作完成")
}
