package controllers

import (
	"net/http"
	"strconv"
	"student-management-system/config"
	"student-management-system/models"

	"github.com/gin-gonic/gin"
)

type StudentController struct{}

// GetAllStudents 获取所有学生
func (sc *StudentController) GetAllStudents(c *gin.Context) {
	var students []models.Student
	db := config.GetDB()
	
	// 确保只获取未删除的记录
	result := db.Where("deleted_at IS NULL").Find(&students)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch students",
			"message": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": students,
		"count": len(students),
	})
}

// GetStudentByID 根据ID获取学生
func (sc *StudentController) GetStudentByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid student ID",
		})
		return
	}
	
	var student models.Student
	db := config.GetDB()
	
	result := db.First(&student, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": student,
	})
}

// CreateStudent 创建新学生
func (sc *StudentController) CreateStudent(c *gin.Context) {
	var student models.Student
	
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
			"message": err.Error(),
		})
		return
	}
	
	db := config.GetDB()
	result := db.Create(&student)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create student",
			"message": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Student created successfully",
		"data": student,
	})
}

// UpdateStudent 更新学生信息
func (sc *StudentController) UpdateStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid student ID",
		})
		return
	}
	
	var student models.Student
	db := config.GetDB()
	
	// 检查学生是否存在
	result := db.First(&student, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}
	
	// 绑定更新数据
	var updateData models.Student
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input data",
			"message": err.Error(),
		})
		return
	}
	
	// 更新学生信息
	result = db.Model(&student).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update student",
			"message": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Student updated successfully",
		"data": student,
	})
}

// DeleteStudent 删除学生
func (sc *StudentController) DeleteStudent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid student ID",
		})
		return
	}
	
	var student models.Student
	db := config.GetDB()
	
	// 检查学生是否存在
	result := db.First(&student, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Student not found",
		})
		return
	}
	
	// 软删除学生
	result = db.Delete(&student)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete student",
			"message": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Student deleted successfully",
	})
}

// SearchStudents 搜索学生
func (sc *StudentController) SearchStudents(c *gin.Context) {
	name := c.Query("name")
	major := c.Query("major")
	grade := c.Query("grade")
	
	var students []models.Student
	db := config.GetDB()
	// 确保只搜索未删除的记录
	query := db.Where("deleted_at IS NULL")
	
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if major != "" {
		query = query.Where("major LIKE ?", "%"+major+"%")
	}
	if grade != "" {
		query = query.Where("grade = ?", grade)
	}
	
	result := query.Find(&students)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search students",
			"message": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": students,
		"count": len(students),
	})
}