package controllers

import (
	"CV-Maker/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	//"github.com/go-redis/redis/v8"
)

func CreateProject(c *gin.Context) {
	// var education models.Education
	// claims := jwt.ExtractClaims(c)
	// user_email, _ := claims["email"]
	//var User models.User
	// user_email, err := Rdb.HGet("user", "email").Result()
	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)
	roleId, _ := models.Rdb.HGet("user", "RoleID").Result()

	if roleId == "" {
		fmt.Println("Redis empty....checking Database for user...")
		err := FillRedis(c)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "something went wrong with redis",
			})
			return
		}
	}

	roleId, _ = models.Rdb.HGet("user", "RoleID").Result()

	if roleId != "2" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Details can only be added by user"})
		return
	}

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be added by admin user"})
	// }

	c.Request.ParseForm()
	var flag bool
	if c.PostForm("title") == "" {
		ReturnParameterMissingError(c, "title")
		flag = true
	}
	if c.PostForm("role") == "" {
		ReturnParameterMissingError(c, "role")
		flag = true
	}
	if c.PostForm("language") == "" {
		ReturnParameterMissingError(c, "language")
		flag = true
	}
	if c.PostForm("client") == "" {
		ReturnParameterMissingError(c, "client")
		flag = true
	}
	if c.PostForm("description") == "" {
		ReturnParameterMissingError(c, "description")
		flag = true
	}
	title := template.HTMLEscapeString(c.PostForm("title"))
	role := template.HTMLEscapeString(c.PostForm("role"))
	language := template.HTMLEscapeString(c.PostForm("language"))
	client := template.HTMLEscapeString(c.PostForm("client"))
	client_name := template.HTMLEscapeString(c.PostForm("client_name"))
	description := template.HTMLEscapeString(c.PostForm("description"))
	// fmt.Println(category_title)
	// fmt.Println("category printed")
	// Check if the category already exists.
	if flag {
		return
	}
	if client == "yes" {
		if c.PostForm("client_name") == "" {
			ReturnParameterMissingError(c, "client_name")
			return
		}
		project := models.Project{
			Title:       title,
			Description: description,
			Role:        role,
			Language:    language,
			Client:      true,
			ClientName:  client_name,
			UserID:      ID,
		}
		err := models.DB.Create(&project).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Project Details added successfully",
		})
	} else {
		project := models.Project{
			Title:       title,
			Description: description,
			Role:        role,
			Language:    language,
			Client:      false,
			ClientName:  "",
			UserID:      ID,
		}
		err := models.DB.Create(&project).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Project Details added successfully",
		})
	}

}

func GetProject(c *gin.Context) {
	var project models.Project
	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)

	err := models.DB.Where("ID = ?", c.Param("id")).Find(&project).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "No data found"})
		return
	}

	if project.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
		return
	}

	c.JSON(200, project)
}

func GetAllProject(c *gin.Context) {
	var project []models.Project

	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)
	email, _ := models.Rdb.HGet("user", "username").Result()

	if email == "" {
		c.JSON(401, gin.H{"message": "unauthorized"})
	}

	err := models.DB.Where("user_id = ?", ID).Find(&project).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "no data found"})
		return
	}

	c.JSON(200, project)
}

func UpdateProject(c *gin.Context) {
	var project models.Project
	var existingProject models.Project

	email, _ := models.Rdb.HGet("user", "username").Result()
	// id, _ := models.Rdb.HGet("user", "ID").Result()
	// ID, _ := strconv.Atoi(id)

	if email == "" {
		fmt.Println("redis empty, checking database for details...")
		FillRedis(c)
		email, _ = models.Rdb.HGet("user", "username").Result()
	}

	//  email, _ = models.Rdb.HGet("user", "email").Result()
	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)

	if email == "" {
		c.JSON(401, gin.H{"message": "unauthorized"})
	}

	err := c.ShouldBindJSON(&project)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	err = models.DB.Where("ID = ?", c.Param("id")).First(&existingProject).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "error fetching from database"})
	}

	if existingProject.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
		return
	}

	models.DB.Model(&existingProject).Updates(project)
}

func DeleteProject(c *gin.Context) {
	var project models.Project

	email, _ := models.Rdb.HGet("user", "username").Result()
	// id, _ := models.Rdb.HGet("user", "ID").Result()
	// ID, _ := strconv.Atoi(id)

	if email == "" {
		fmt.Println("redis empty, checking database for details...")
		FillRedis(c)
	}

	email, _ = models.Rdb.HGet("user", "username").Result()
	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)

	if email == "" {
		c.JSON(401, gin.H{"message": "unauthorized"})
	}

	err := models.DB.Where("ID = ?", c.Param("id")).First(&project).Error
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
	}

	if project.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
	}

	models.DB.Where("ID = ?", c.Param("id")).Delete(&project)

}

// func createPagination(c *gin.Context, page int) {
// 	var
// }
