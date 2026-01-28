package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/silaeder-labs/bank/backend/auth"
	"github.com/silaeder-labs/bank/backend/config"
	"github.com/silaeder-labs/bank/backend/handlers"
	"github.com/silaeder-labs/bank/backend/routes"

	echoMw "github.com/labstack/echo/v4/middleware"
	echokitMw "github.com/nrf24l01/go-web-utils/echokit/middleware"
	echokitSchemas "github.com/nrf24l01/go-web-utils/echokit/schemas"
	"github.com/nrf24l01/go-web-utils/pgkit"

	gologger "github.com/nrf24l01/go-logger"
)

func main() {
	ctx := context.Background()

	// Logger create
	logger := gologger.NewLogger(os.Stdout, "bank",
		gologger.WithTypeColors(map[gologger.LogType]string{
			gologger.LogType("HTTP"):  gologger.BgCyan,
			gologger.LogType("DB"):  gologger.BgGreen,
			gologger.LogType("SETUP"): gologger.BgRed,
		}),
	)
	log.Printf("Logger initialized")

	err := godotenv.Load(".env")
	if err != nil {
		logger.Log(gologger.LevelWarn, gologger.LogType("SETUP"), fmt.Sprintf("Failed to load .env file: %v", err), "")
	} else {
		logger.Log(gologger.LevelInfo, gologger.LogType("SETUP"), ".env file loaded", "")
	}

	// Configuration initialization
	config, err := config.BuildConfigFromEnv()
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to build config: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelInfo, gologger.LogType("SETUP"), "Configuration loaded", "")
	}

	// Data sources initialization
	db, err := pgkit.NewDB(ctx, config.PGConfig)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to connect to postgres: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelInfo, gologger.LogType("SETUP"), "Connected to Postgres database", "")
	}
	err = pgkit.RunMigrations(db.SQL, config.PGConfig)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to run migrations: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelInfo, gologger.LogType("SETUP"), "Migrations ran successfully", "")
	}

	// Keycloak key verifier init
	jwks, err := auth.RegisterJwks(config.KeyCloakConfig, logger, &ctx)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to register JWKS: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelInfo, gologger.LogType("SETUP"), "JWKS registered", "")
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

	e.Use(echokitMw.RequestLogger(logger))

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
	handler := &handlers.Handler{DB: db, Config: config, Logger: logger, Jwks: jwks}
	routes.RegisterRoutes(api, handler)

	// Start server
	e.Logger.Fatal(e.Start(config.WebAppConfig.AppHost))
}
