package routes

import (
	"student-management-system/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(router *gin.Engine) {
	studentController := &controllers.StudentController{}
	
	// API路由组
	api := router.Group("/api/v1")
	{
		// 学生相关路由
		students := api.Group("/students")
		{
			// GET /api/v1/students - 获取所有学生
			students.GET("", studentController.GetAllStudents)
			
			// GET /api/v1/students/search - 搜索学生
			students.GET("/search", studentController.SearchStudents)
			
			// GET /api/v1/students/:id - 根据ID获取学生
			students.GET("/:id", studentController.GetStudentByID)
			
			// POST /api/v1/students - 创建新学生
			students.POST("", studentController.CreateStudent)
			
			// PUT /api/v1/students/:id - 更新学生信息
			students.PUT("/:id", studentController.UpdateStudent)
			
			// DELETE /api/v1/students/:id - 删除学生
			students.DELETE("/:id", studentController.DeleteStudent)
		}
	}
	
	// 静态文件服务
	router.Static("/static", "./views/static")
	
	// HTML模板路由
	router.LoadHTMLGlob("views/templates/*")
	
	// 前端页面路由
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "学生管理系统",
		})
	})
	
	// 学生管理页面
	router.GET("/students", func(c *gin.Context) {
		c.HTML(200, "students.html", gin.H{
			"title": "学生管理",
		})
	})
}