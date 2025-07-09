package config

import (
	"log"
	"student-management-system/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase() {
	var err error
	
	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open("students.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// 自动迁移数据库表
	err = DB.AutoMigrate(&models.Student{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("Database connected and migrated successfully")
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}