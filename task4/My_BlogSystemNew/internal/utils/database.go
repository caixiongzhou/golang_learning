package utils

import (
	"blog-system/internal/config"
	"blog-system/internal/model"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %v", err)
	}

	log.Println("数据库表结构迁移成功")
	return db, nil
}

// 说明：
// - 函数使用 GORM 连接 MySQL 并执行自动迁移。
// - DSN 字符串使用 utf8mb4 编码并启用 parseTime，以确保时间类型字段的正确处理。
// - 在生产环境中应考虑连接池配置、超时和错误重试策略。