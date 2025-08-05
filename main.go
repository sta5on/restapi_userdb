package main

import (
	"github.com/labstack/echo/v4"
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

var users = make(map[int]User)
var nextID int = 1

func getHandler(c echo.Context) error {
	var usersSlice []User

	for _, msg := range users {
		usersSlice = append(usersSlice, msg)
	}
	return c.JSON(http.StatusOK, &usersSlice)
}

func postHandler(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt create user",
		})
	}

	user.Id = nextID
	nextID++

	users[user.Id] = user
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
	if err := c.Bind(&updatedUser); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, Response{
			Status:  "Error",
			Message: "Couldnt update message",
		})
	}

	if _, exist := users[id]; !exist {
		return c.JSON(http.StatusNotFound, Response{
			Status:  "Error",
			Message: "Couldnt find the message",
		})
	}

	updatedUser.Id = id
	users[id] = updatedUser

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

	if _, exist := users[id]; !exist {
		return c.JSON(http.StatusNotFound, Response{
			Status:  "Error",
			Message: "Couldnt find the user",
		})
	}

	delete(users, id)
	return c.JSON(http.StatusOK, Response{
		Status:  "Success",
		Message: "User was successfully deleted",
	})
}

func main() {
	e := echo.New()

	log.Println("app is running")

	e.GET("/users", getHandler)
	e.POST("/users", postHandler)
	e.PATCH("/users/:id", patchHandler)
	e.DELETE("/users/:id", deleteHandler)

	e.Start(":8080")
}
