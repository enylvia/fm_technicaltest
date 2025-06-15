package main

import (
	"FM_techincaltest/app"
	"FM_techincaltest/app/database"
	"FM_techincaltest/middleware"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	app.LoadConfig()
	log.Printf("App running on port: %s", app.Config.AppPort)
	log.Printf("Connecting to DB...")

	dbClient, err := database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer dbClient.Close()

	e := echo.New()
	userGroup := e.Group("/user")
	{
		// TODO REGISTER AND LOGIN USER WITH LOG
		userGroup.POST("/register")
		userGroup.POST("/login")
	}
	employeeGroup := e.Group("/employee")
	employeeGroup.Use(middleware.AuthenticateMiddleware())
	{
		// TODO ABSENCE WITH LOG
		employeeGroup.POST("/absence")
	}

	log.Fatal(e.Start(":" + app.Config.AppPort))
}
