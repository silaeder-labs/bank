package main

import (
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
	pgKit "github.com/nrf24l01/go-web-utils/pg_kit"

	gologger "github.com/nrf24l01/go-logger"
)

func main() {
	// Logger create
	logger := gologger.NewLogger(os.Stdout, "bank", gologger.WithTypeColors(map[gologger.LogType]string{
		gologger.LogType("HTTP"): gologger.BgCyan,
		gologger.LogType("AUTH"): gologger.BgGreen,
		gologger.LogType("SETUP"): gologger.BgRed,
	}),
	)

	err := godotenv.Load(".env")
	if err != nil {
		logger.Log(gologger.LevelWarn, gologger.LogType("SETUP"), fmt.Sprintf("Failed to load .env file: %v", err), "")
	}

	// Configuration initialization
	config, err := config.BuildConfigFromEnv()
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to build config: %v", err), "")
		return
	}

	// Data sources initialization
	db, err := pgKit.RegisterPostgres(config.PGConfig, false)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to connect to postgres: %v", err), "")
		return
	}

	// Keycloak key verifier init
	jwks, err := auth.RegisterJwks(config.KeyCloakConfig, logger)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to register JWKS: %v", err), "")
		return
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
