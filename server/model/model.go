package model

import (
	"funtabs-server/config"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 是全局数据库实例
var DB *gorm.DB

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Password  string    `gorm:"size:128;not null" json:"-"` // bcrypt 哈希，不对外暴露
	Avatar    string    `gorm:"size:512" json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserData 存储用户的 localStorage 快照
type UserData struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Data      string    `gorm:"type:text" json:"data"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Init 初始化数据库连接并自动迁移表结构
func Init() {
	dbPath := config.Cfg.DB.Path

	// 确保数据库文件所在目录存在
	if dir := filepath.Dir(dbPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("创建数据库目录失败: %v", err)
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	if err = DB.AutoMigrate(&User{}, &UserData{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	log.Printf("数据库已就绪: %s", dbPath)
}
