package api

import (
	"context"
	"runtime"
	"time"

	"github.com/LouisFernando1204/kai-backend.git/internal/connection"
	"github.com/gofiber/fiber/v2"
)

var startTime = time.Now()

// StatusHandler handles status endpoint
type statusHandler struct{}

// NewStatusHandler creates a new status handler
func NewStatusHandler() *statusHandler {
	return &statusHandler{}
}

// GetStatus returns comprehensive service status
// @Summary Get service status
// @Description Get detailed service status including database connectivity, uptime, and system information
// @Tags System
// @Produce json
// @Success 200 {object} map[string]interface{} "Service status"
// @Failure 503 {object} map[string]interface{} "Service unavailable"
// @Router /status [get]
func (h *statusHandler) GetStatus(c *fiber.Ctx) error {
	// Check database connection
	dbStatus := "disconnected"
	dbError := ""

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := connection.GetClient()
	if client != nil {
		err := client.Ping(ctx, nil)
		if err == nil {
			dbStatus = "connected"
		} else {
			dbError = err.Error()
		}
	}

	// Calculate uptime
	uptime := time.Since(startTime)

	// Get memory stats
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	status := fiber.Map{
		"status":    "healthy",
		"service":   "kai-backend",
		"version":   "1.0.0",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    uptime.String(),
		"database": fiber.Map{
			"status": dbStatus,
			"type":   "mongodb",
		},
		"system": fiber.Map{
			"go_version":    runtime.Version(),
			"num_goroutine": runtime.NumGoroutine(),
			"num_cpu":       runtime.NumCPU(),
			"memory": fiber.Map{
				"alloc_mb":       memStats.Alloc / 1024 / 1024,
				"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
				"sys_mb":         memStats.Sys / 1024 / 1024,
			},
		},
	}

	// If database is disconnected, return 503
	if dbStatus == "disconnected" {
		status["status"] = "unhealthy"
		if dbError != "" {
			status["database"].(fiber.Map)["error"] = dbError
		}
		return c.Status(fiber.StatusServiceUnavailable).JSON(status)
	}

	return c.Status(fiber.StatusOK).JSON(status)
}
