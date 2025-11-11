// @title KAI Backend
// @version 1.0
// @description This is a comprehensive REST API for Indonesian Railway (KAI) ticket booking system.
// @description The API provides endpoints for user authentication, station management, train management, schedule search, and booking operations.

// @contact.name Louis Fernando
// @contact.email fernandolouis55@gmail.com

// @license.name MIT
// @license.url http://opensource.org/licenses/MIT

// @host 103.161.185.247:8080
// @schemes http https
// @BasePath /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Authorization header using the Bearer scheme. Example: "Authorization: Bearer {token}"

package main

import (
	"log"

	_ "github.com/LouisFernando1204/kai-backend.git/docs"
	"github.com/LouisFernando1204/kai-backend.git/internal/config"
	"github.com/LouisFernando1204/kai-backend.git/internal/connection"
	"github.com/LouisFernando1204/kai-backend.git/internal/repository"
	"github.com/LouisFernando1204/kai-backend.git/internal/router"
	"github.com/LouisFernando1204/kai-backend.git/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cnf, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	connection.ConnectDatabase(cnf.Database)
	defer connection.DisconnectDatabase()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: false,
	}))
	app.Use(logger.New())

	userRepo := repository.NewUserRepository(cnf.Database)
	stationRepo := repository.NewStationRepository(cnf.Database)
	trainRepo := repository.NewTrainRepository(cnf.Database)
	scheduleRepo := repository.NewScheduleRepository(cnf.Database)
	bookingRepo := repository.NewBookingRepository(cnf.Database)

	authService := service.NewAuthService(userRepo, cnf)
	stationService := service.NewStationService(stationRepo)
	trainService := service.NewTrainService(trainRepo)
	scheduleService := service.NewScheduleService(scheduleRepo, trainRepo, stationRepo, bookingRepo)
	bookingService := service.NewBookingService(bookingRepo, scheduleRepo, scheduleService)

	router.SetupRoutes(app, cnf, authService, stationService, trainService, scheduleService, bookingService)

	serverAddr := cnf.Server.Host + ":" + cnf.Server.Port
	log.Printf("Server is running on http://%s", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}
