package main

import (
	"github.com/Tanakaryuki/brachio-backend/config"
	"github.com/Tanakaryuki/brachio-backend/db"
	"github.com/Tanakaryuki/brachio-backend/firebase"
	"github.com/Tanakaryuki/brachio-backend/migrate"
	"github.com/Tanakaryuki/brachio-backend/routes"
	"github.com/labstack/echo/v4"
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
	routes.RegisterRoutes(server)

	server.Logger.Fatal(server.Start(":5050"))
}
