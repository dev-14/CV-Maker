package controllers

import (
	"CV-Maker/models"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	//"github.com/go-redis/redis/v8"
)

func CreateEducation(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Educational Qualifications can only be added by user"})
		return
	}

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be added by admin user"})
	// }

	c.Request.ParseForm()
	var flag bool
	if c.PostForm("institution") == "" {
		ReturnParameterMissingError(c, "institution")
		flag = true
	}
	if c.PostForm("degree") == "" {
		ReturnParameterMissingError(c, "degree")
		flag = true
	}
	if c.PostForm("start_month") == "" {
		ReturnParameterMissingError(c, "start_month")
		flag = true
	}
	if c.PostForm("start_year") == "" {
		ReturnParameterMissingError(c, "start_year")
		flag = true
	}
	if c.PostForm("end_month") == "" {
		ReturnParameterMissingError(c, "end_month")
		flag = true
	}
	if c.PostForm("end_year") == "" {
		ReturnParameterMissingError(c, "end_year")
		flag = true
	}
	if c.PostForm("grade") == "" {
		ReturnParameterMissingError(c, "grade")
		flag = true
	}
	if c.PostForm("description") == "" {
		ReturnParameterMissingError(c, "description")
		flag = true
	}
	institution := template.HTMLEscapeString(c.PostForm("institution"))
	degree := template.HTMLEscapeString(c.PostForm("degree"))
	start_month := template.HTMLEscapeString(c.PostForm("start_month"))
	start_year := template.HTMLEscapeString(c.PostForm("start_year"))
	end_month := template.HTMLEscapeString(c.PostForm("end_month"))
	end_year := template.HTMLEscapeString(c.PostForm("end_year"))
	grade := template.HTMLEscapeString(c.PostForm("grade"))
	description := template.HTMLEscapeString(c.PostForm("description"))
	// fmt.Println(category_title)
	// fmt.Println("category printed")
	// Check if the category already exists.
	if flag {
		return
	}

	startYear, _ := strconv.Atoi(start_year)
	startMonth, _ := strconv.Atoi(start_month)
	start_date := time.Date(startYear, time.Month(startMonth), 1, 0, 0, 0, 0, time.Local)

	endYear, _ := strconv.Atoi(end_year)
	endMonth, _ := strconv.Atoi(end_month)
	end_date := time.Date(endYear, time.Month(endMonth), 1, 0, 0, 0, 0, time.Local)

	education := models.Education{
		Institution: institution,
		Degree:      degree,
		StartDate:   start_date,
		EndDate:     end_date,
		Grade:       grade,
		Description: description,
		UserID:      ID,
		User:        models.User{},
	}

	err := models.DB.Create(&education).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Educational Details added successfully",
	})

}
func GetEducation(c *gin.Context) {
	var education models.Education

	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)

	err := models.DB.Where("ID = ?", c.Param("id")).Find(&education).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "No data found"})
		return
	}

	if education.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
		return
	}

	c.JSON(200, education)
}

func GetAllEducation(c *gin.Context) {
	var e []models.Education

	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)
	email, _ := models.Rdb.HGet("user", "username").Result()

	if email == "" {
		c.JSON(401, gin.H{"message": "unauthorized"})
	}

	p := c.Query("page")
	page, _ := strconv.Atoi(p)
	order := c.Query("order")
	err := models.DB.Where("user_id = ?", ID).Order(order).Limit(2).Offset((page - 1) * 2).Find(&e).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "no data found"})
		return
	}

	c.JSON(200, e)
}

func UpdateEducation(c *gin.Context) {
	var education models.Education
	var existingEducation models.Education

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

	err := c.ShouldBindJSON(&education)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	err = models.DB.Where("ID = ?", c.Param("id")).First(&existingEducation).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "error fetching from database"})
	}

	if existingEducation.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
		return
	}

	models.DB.Model(&existingEducation).Updates(education)
}

func DeleteEducation(c *gin.Context) {
	var education models.Education

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

	err := models.DB.Where("ID = ?", c.Param("id")).First(&education).Error
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
	}

	if education.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
	}

	models.DB.Where("ID = ?", c.Param("id")).Delete(&education)

}

func GetEducationByYear(c *gin.Context) {
	// var allEducation []models.Education
	var education []models.Education
	getYear := c.Query("year")
	// getYear, _ := strconv.Atoi(template.HTMLEscapeString(c.PostForm("year")))
	email, _ := models.Rdb.HGet("user", "username").Result()
	fmt.Println(email)
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
		return
	}

	err := models.DB.Where("user_id=? AND start_date <= ?", ID, getYear).Or("user_id=? AND end_date <= ?", ID, getYear).Find(&education).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "error fetching from database"})
		return
	}
	// err := models.DB.Where("belongs_to_id = ?", ID).Find(&allEducation).Error
	// if err != nil {
	// 	c.JSON(404, gin.H{"message": "error fetching from database"})
	// 	return
	// }

	// for _, edu := range allEducation {
	// 	start_year := edu.StartDate.Year()
	// 	end_year := edu.EndDate.Year()
	// 	fmt.Println(start_year, end_year, getYear)
	// 	if start_year == getYear || end_year == getYear || (start_year < getYear && getYear < end_year) {
	// 		education = append(education, edu)
	// 	}
	// }
	c.JSON(200, education)
}
