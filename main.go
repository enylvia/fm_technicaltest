package main

import (
	"FM_techincaltest/app"
	"FM_techincaltest/app/database"
	_ "FM_techincaltest/docs"
	"FM_techincaltest/handler"
	"FM_techincaltest/middleware"
	"FM_techincaltest/repository"
	"FM_techincaltest/service"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
)

// @title FM Technical Test API
// @version 1.0
// @description Dokumentasi API untuk FM Technical Test
// @host localhost:50001
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token JWT kamu dengan format: Bearer <token>
func main() {
	uploadDir := "./upload/image/"
	app.LoadConfig()
	log.Printf("App running on port: %s", app.Config.AppPort)
	log.Printf("Connecting to DB...")

	dbClient, err := database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer dbClient.Close()
	// Depedency Injection
	userRepository := repository.NewUserRepository(dbClient)
	userService := service.NewUserServiceImplement(userRepository)
	userHandler := handler.NewUserHandlerImplement(userService)

	employeeRepository := repository.NewEmployeeAbsenceImplement(dbClient)
	employeeService := service.NewEmployeeServiceImplement(userRepository, employeeRepository)
	employeeHandler := handler.NewEmployeeHandlerImplement(employeeService)

	imageRepo := repository.NewImageRepository(uploadDir)
	imageService := service.NewImageService(imageRepo)
	imageHandler := handler.NewImageHandler(imageService)

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1")
	userGroup := v1.Group("/user")
	{
		userGroup.POST("/register", userHandler.RegisterUserAndEmployee)
		userGroup.POST("/login", userHandler.LoginUser)
	}
	employeeGroup := v1.Group("/employee")
	employeeGroup.Use(middleware.AuthenticateMiddleware())
	{
		employeeGroup.POST("/clock_in", employeeHandler.ClockInRequest)
		employeeGroup.POST("/clock_out", employeeHandler.ClockOutRequest)
		employeeGroup.GET("/absence/log", employeeHandler.AbsenceHistory)
	}
	imageGroup := v1.Group("/image")
	imageGroup.Use(middleware.AuthenticateMiddleware())
	{
		imageGroup.POST("/save", imageHandler.UploadImage)

	}
	e.Static("/uploads", uploadDir)
	log.Fatal(e.Start(":" + app.Config.AppPort))
}
