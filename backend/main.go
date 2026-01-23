package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/silaeder-labs/bank/backend/config"
	"github.com/silaeder-labs/bank/backend/handlers"
	"github.com/silaeder-labs/bank/backend/routes"

	echoMw "github.com/labstack/echo/v4/middleware"
	echokitMw "github.com/nrf24l01/go-web-utils/echokit/middleware"
	echokitSchemas "github.com/nrf24l01/go-web-utils/echokit/schemas"
	pgKit "github.com/nrf24l01/go-web-utils/pg_kit"

	gologger "github.com/nrf24l01/go-logger"
)

func main() {
	logger := slog.Default()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Warn("Failed to load .env file: %v", slog.Any("error", err))
	}

	// Configuration initialization
	config, err := config.BuildConfigFromEnv()
	if err != nil {
		logger.Error("failed to build config", slog.Any("error", err))
	}

	// Data sources initialization
	db, err := pgKit.RegisterPostgres(config.PGConfig, false)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// Create echo object
	e := echo.New()

	// Register custom validator
	v := validator.New()
	e.Validator = &echokitMw.CustomValidator{Validator: v}

	// Echo Configs
	e.Use(echoMw.Recover())
	e.Use(echoMw.RemoveTrailingSlash())
	e.Use(echokitMw.TraceMiddleware())
	// create and pass colored console logger into middleware
	// register custom LogType "HTTP" with cyan background
	consoleLogger := gologger.NewLogger(os.Stdout, "bank-backend", gologger.WithTypeColors(map[gologger.LogType]string{
		gologger.LogType("HTTP"): gologger.BgCyan,
		gologger.LogType("AUTH"): gologger.BgGreen,
	}),
	)
	e.Use(echokitMw.RequestLogger(consoleLogger))

	// Cors
	log.Printf("Setting allowed origin to: %s", config.WebAppConfig.AllowOrigin)
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins:     []string{config.WebAppConfig.AllowOrigin},
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS, echo.DELETE, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Api group
	api := e.Group("")

	// Health check endpoint
	api.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, echokitSchemas.Message{Status: "Sl-eco/bank backend is OK"})
	})

	// Register routes
	handler := &handlers.Handler{DB: db, Config: config}
	routes.RegisterRoutes(api, handler)

	// Start server
	e.Logger.Fatal(e.Start(config.WebAppConfig.AppHost))
}
