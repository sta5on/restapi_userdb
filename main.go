package main

import (
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}
	db.AutoMigrate(&User{})
}

// main gorm methods, find, create, update, delete

func getHandler(c echo.Context) error {
	var users []User

	if err := db.Find(&users).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt find users",
		})
	}

	return c.JSON(http.StatusOK, &users)
}

func postHandler(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt add the user",
		})
	}

	if err := db.Create(&user).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt create user",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "Succes",
		Message: "User was successfully created",
	})

}

func patchHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt parse id",
		})
	}

	var updatedUser User
	if err = c.Bind(&updatedUser); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Invalid input",
		})
	}

	if err = db.Model(&User{}).Where("id = ?", id).Update("name", updatedUser.Name).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt update user",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "Message was successfully updated",
	})
}

func deleteHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt parse id",
		})
	}

	if err = db.Delete(&User{}, id).Error; err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt delete user",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "User was successfully deleted",
	})
}

func main() {
	initDB()
	e := echo.New()

	log.Println("app is running")

	e.GET("/v1/users", getHandler)
	e.POST("/v1/users", postHandler)
	e.PATCH("/v1/users/:id", patchHandler)
	e.DELETE("/v1/users/:id", deleteHandler)

	e.Start(":8080")
}
