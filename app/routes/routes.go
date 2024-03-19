package routes

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(server *echo.Echo) {
	server.POST("/users", createUser)
	server.GET("/users", getAllUsers)
	server.GET("/users/:userId", getUserById)
}
