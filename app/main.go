package main

import (
	"net/http"
	"os"

	"github.com/Tanakaryuki/brachio-backend/config"
	"github.com/Tanakaryuki/brachio-backend/db"
	"github.com/Tanakaryuki/brachio-backend/firebase"
	"github.com/Tanakaryuki/brachio-backend/migrate"
	"github.com/Tanakaryuki/brachio-backend/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.LoadEnv()
	db := db.Init()
	migrate.AutoMigrate(db)

	err := firebase.InitFirebase("./serviceAccountKey.json")
	if err != nil {
		panic(err)
	}

	server := echo.New()
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
		Output: os.Stdout,
	}))

	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HelloWorld")
	})
	routes.RegisterRoutes(server)

	server.Logger.Fatal(server.Start(":5050"))
}
