package main

import (
	"log"
	"student-management-system/config"
	"student-management-system/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	config.InitDatabase()
	
	// 创建Gin路由器
	router := gin.Default()
	
	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// 设置路由
	routes.SetupRoutes(router)
	
	// 启动服务器
	log.Println("Starting server on :8080...")
	log.Println("API Documentation:")
	log.Println("GET    /api/v1/students        - 获取所有学生")
	log.Println("GET    /api/v1/students/search - 搜索学生 (参数: name, major, grade)")
	log.Println("GET    /api/v1/students/:id    - 根据ID获取学生")
	log.Println("POST   /api/v1/students        - 创建新学生")
	log.Println("PUT    /api/v1/students/:id    - 更新学生信息")
	log.Println("DELETE /api/v1/students/:id    - 删除学生")
	log.Println("")
	log.Println("Web Interface:")
	log.Println("http://localhost:8080          - 首页")
	log.Println("http://localhost:8080/students - 学生管理页面")
	
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}