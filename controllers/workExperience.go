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

func CreateWorkExperience(c *gin.Context) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Work Experiences can only be added by user"})
		return
	}

	// // Check if the current user had admin role.
	// if err := models.DB.Where("email = ? AND user_role_id=1", user_email).First(&User).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Category can only be added by admin user"})
	// }

	c.Request.ParseForm()
	var flag bool
	if c.PostForm("company_name") == "" {
		ReturnParameterMissingError(c, "company_name")
		flag = true
	}
	if c.PostForm("degree") == "" {
		ReturnParameterMissingError(c, "employment_type")
		flag = true
	}
	if c.PostForm("start_date") == "" {
		ReturnParameterMissingError(c, "start_date")
		flag = true
	}
	if c.PostForm("end_date") == "" {
		ReturnParameterMissingError(c, "end_date")
		flag = true
	}
	if c.PostForm("grade") == "" {
		ReturnParameterMissingError(c, "job_role")
		flag = true
	}
	if c.PostForm("grade") == "" {
		ReturnParameterMissingError(c, "job_location")
		flag = true
	}
	if c.PostForm("description") == "" {
		ReturnParameterMissingError(c, "description")
		flag = true
	}
	company_name := template.HTMLEscapeString(c.PostForm("company_name"))
	employment_type := template.HTMLEscapeString(c.PostForm("employment_type"))
	start_month := template.HTMLEscapeString(c.PostForm("start_month"))
	start_year := template.HTMLEscapeString(c.PostForm("start_year"))
	end_month := template.HTMLEscapeString(c.PostForm("end_month"))
	end_year := template.HTMLEscapeString(c.PostForm("end_year"))
	job_role := template.HTMLEscapeString(c.PostForm("job_role"))
	job_location := template.HTMLEscapeString(c.PostForm("job_location"))
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
	// err := models.DB.Where("start_date = ? AND end_date = ?", start_date, end_date).First(&education).Error
	// if err == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "date."})
	// 	return
	// }

	work := &models.WorkExperiences{
		CompanyName:    company_name,
		EmploymentType: employment_type,
		From:           start_date,
		To:             end_date,
		JobRole:        job_role,
		JobLocation:    job_location,
		Description:    description,
		UserID:         ID,
		User:           models.User{},
	}

	err := models.DB.Create(&work).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Work Experience added successfully",
	})

}

func GetWorkExperience(c *gin.Context) {
	var work models.WorkExperiences

	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)

	err := models.DB.Where("ID = ?", c.Param(id)).Find(&work).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "No data found"})
		return
	}

	if work.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
		return
	}

	c.JSON(200, work)
}

func GetAllWorkExperience(c *gin.Context) {
	var works []models.Education

	id, _ := models.Rdb.HGet("user", "ID").Result()
	ID, _ := strconv.Atoi(id)
	email, _ := models.Rdb.HGet("user", "username").Result()

	if email == "" {
		c.JSON(401, gin.H{"message": "unauthorized"})
	}

	p := c.Query("page")
	page, _ := strconv.Atoi(p)
	order := c.Query("order")
	err := models.DB.Where("user_id = ?", ID).Order(order).Limit(2).Offset((page - 1) * 2).Find(&works).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "no data found"})
		return
	}

	c.JSON(200, works)
}

func UpdateWorkExperience(c *gin.Context) {
	var work models.WorkExperiences
	var existingWork models.WorkExperiences

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

	err := c.ShouldBindJSON(&work)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	err = models.DB.Where("ID = ?", c.Param("id")).First(&existingWork).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "error fetching from database"})
	}

	if existingWork.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
		return
	}

	models.DB.Model(&existingWork).Updates(work)
}

func DeleteWorkExperience(c *gin.Context) {
	var work models.WorkExperiences

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

	err := models.DB.Where("ID = ?", c.Param("id")).First(&work).Error
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
	}

	if work.UserID != ID {
		c.JSON(400, gin.H{"message": "invalid requested ID"})
	}

	models.DB.Where("ID = ?", c.Param("id")).Delete(&work)

}

func GetWorkExperienceByYear(c *gin.Context) {
	// var allWork []models.WorkExperiences
	var work []models.WorkExperiences
	getYear, _ := strconv.Atoi(template.HTMLEscapeString(c.PostForm("year")))
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

	err := models.DB.Where("user_id = ? AND from <= ?", ID, getYear).Or("user_id = ? AND to >= ?", ID, getYear).Find(&work).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "error fetching from database"})
		return
	}

	// err := models.DB.Where("belongs_to_id = ?", ID).Find(&allWork).Error
	// if err != nil {
	// 	c.JSON(404, gin.H{"message": "error fetching from database"})
	// 	return
	// }

	// for _, edu := range allWork {
	// 	start_year := edu.From.Year()
	// 	end_year := edu.To.Year()
	// 	fmt.Println(start_year, end_year, getYear)

	// 	if start_year == getYear || end_year == getYear || (start_year < getYear && getYear < end_year) {
	// 		work = append(work, edu)
	// 	}
	// }
	c.JSON(200, work)
}
