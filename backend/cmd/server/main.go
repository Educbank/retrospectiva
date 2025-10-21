package main

import (
	"log"
	"os"

	"educ-retro/internal/database"
	"educ-retro/internal/handlers"
	"educ-retro/internal/repositories"
	"educ-retro/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	teamRepo := repositories.NewTeamRepository(database.DB)
	retroRepo := repositories.NewRetrospectiveRepository(database.DB)

	// Initialize services
	userService := services.NewUserService(userRepo)
	teamService := services.NewTeamService(teamRepo, userRepo)
	templateService := services.NewTemplateService()
	retrospectiveService := services.NewRetrospectiveService(retroRepo, teamRepo)

	// Initialize Realtime service
	realtimeService := services.NewRealtimeService()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	templateHandler := handlers.NewTemplateHandler(templateService)
	retrospectiveHandler := handlers.NewRetrospectiveHandler(retrospectiveService, realtimeService)
	sseHandler := handlers.NewSSEHandler(realtimeService)

	// Setup router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	v1 := r.Group("/api/v1")
	{
		userHandler.SetupRoutes(v1)
		teamHandler.SetupRoutes(v1)
		templateHandler.SetupRoutes(v1)
		retrospectiveHandler.SetupRoutes(v1)
		sseHandler.SetupRoutes(v1)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "educ-retro",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}
