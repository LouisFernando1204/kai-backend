package router

import (
	"github.com/LouisFernando1204/kai-backend.git/domain"
	"github.com/LouisFernando1204/kai-backend.git/internal/api"
	"github.com/LouisFernando1204/kai-backend.git/internal/config"
	"github.com/LouisFernando1204/kai-backend.git/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(
	app *fiber.App,
	conf *config.Config,
	authService domain.AuthService,
	stationService domain.StationService,
	trainService domain.TrainService,
	scheduleService domain.ScheduleService,
	bookingService domain.BookingService,
) {
	authHandler := api.NewAuthHandler(authService)
	stationHandler := api.NewStationHandler(stationService)
	trainHandler := api.NewTrainHandler(trainService)
	scheduleHandler := api.NewScheduleHandler(scheduleService)
	bookingHandler := api.NewBookingHandler(bookingService)
	statusHandler := api.NewStatusHandler()

	app.Get("/docs/*", swagger.HandlerDefault)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to KAI Backend!"})
	})

	// Health check endpoint for Docker
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"service": "kai-backend",
		})
	})

	// Detailed status endpoint for monitoring
	app.Get("/status", statusHandler.GetStatus)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	api.Get("/stations", stationHandler.GetAll)
	api.Get("/stations/:id", stationHandler.GetByID)
	adminStations := api.Group("/stations", middleware.Protected(conf))
	adminStations.Post("/", stationHandler.Create)
	adminStations.Put("/:id", stationHandler.Update)
	adminStations.Delete("/:id", stationHandler.Delete)

	api.Get("/trains", trainHandler.GetAll)
	api.Get("/trains/:id", trainHandler.GetByID)
	adminTrains := api.Group("/trains", middleware.Protected(conf))
	adminTrains.Post("/", trainHandler.Create)
	adminTrains.Put("/:id", trainHandler.Update)
	adminTrains.Delete("/:id", trainHandler.Delete)

	schedules := api.Group("/schedules")
	schedules.Post("/search", scheduleHandler.Search)
	schedules.Get("/:id", scheduleHandler.GetByID)
	schedules.Get("/:id/seats", scheduleHandler.GetSeatLayout)
	adminSchedules := api.Group("/schedules", middleware.Protected(conf))
	adminSchedules.Post("/", scheduleHandler.Create)
	adminSchedules.Get("/", scheduleHandler.GetAll)
	adminSchedules.Put("/:id", scheduleHandler.Update)
	adminSchedules.Delete("/:id", scheduleHandler.Delete)

	bookings := api.Group("/bookings", middleware.Protected(conf))
	bookings.Post("/", bookingHandler.CreateBooking)
	bookings.Get("/", bookingHandler.GetMyBookings)
	bookings.Get("/:id", bookingHandler.GetBookingByID)
	bookings.Patch("/:id/status", bookingHandler.UpdateStatus)
	bookings.Post("/:id/cancel", bookingHandler.CancelBooking)
}
