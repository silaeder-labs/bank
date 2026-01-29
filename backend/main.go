package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/silaeder-labs/bank/backend/auth"
	"github.com/silaeder-labs/bank/backend/config"
	"github.com/silaeder-labs/bank/backend/handlers"
	"github.com/silaeder-labs/bank/backend/postgres"
	"github.com/silaeder-labs/bank/backend/routes"

	echoMw "github.com/labstack/echo/v4/middleware"
	echokitMw "github.com/nrf24l01/go-web-utils/echokit/middleware"
	echokitSchemas "github.com/nrf24l01/go-web-utils/echokit/schemas"
	"github.com/nrf24l01/go-web-utils/pgkit"

	gologger "github.com/nrf24l01/go-logger"
)

func main() {
	grantUnlimitedFlag := flag.String("grant-unlimited", "", "user UUID to grant unlimited balance")
	revokeUnlimitedFlag := flag.String("revoke-unlimited", "", "user UUID to revoke unlimited balance")
	flag.Parse()

	ctx := context.Background()

	// Logger create
	logger := gologger.NewLogger(os.Stdout, "bank",
		gologger.WithTypeColors(map[gologger.LogType]string{
			gologger.LogType("HTTP"):  gologger.BgCyan,
			gologger.LogType("DB"):    gologger.BgGreen,
			gologger.LogType("SETUP"): gologger.BgRed,
			gologger.LogType("AUTH"):  gologger.BgMagenta,
			gologger.LogType("CLI"): gologger.BgCyan,
		}),
	)
	log.Printf("Logger initialized")

	err := godotenv.Load(".env")
	if err != nil {
		logger.Log(gologger.LevelWarn, gologger.LogType("SETUP"), fmt.Sprintf("Failed to load .env file: %v", err), "")
	} else {
		logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), ".env file loaded", "")
	}

	// Configuration initialization
	config, err := config.BuildConfigFromEnv()
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to build config: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), "Configuration loaded", "")
	}

	// Data sources initialization
	db, err := pgkit.NewDB(ctx, config.PGConfig)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to connect to postgres: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), "Connected to Postgres database", "")
	}
	err = pgkit.RunMigrations(db.SQL, config.PGConfig)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to run migrations: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), "Migrations ran successfully", "")
	}

	// Keycloak key verifier init
	jwks, err := auth.RegisterJwks(config.KeyCloakConfig, logger, &ctx)
	if err != nil {
		logger.Log(gologger.LevelFatal, gologger.LogType("SETUP"), fmt.Sprintf("Failed to register JWKS: %v", err), "")
		return
	} else {
		logger.Log(gologger.LevelSuccess, gologger.LogType("SETUP"), "JWKS registered", "")
	}

	// Unlimited balance CLI commands
	if *grantUnlimitedFlag != "" {
		targetID, err := uuid.Parse(*grantUnlimitedFlag)
		if err != nil {
			logger.Log(gologger.LevelFatal, gologger.LogType("CLI"), fmt.Sprintf("invalid grant-unlimited uuid: %v", err), "")
			return
		}
		if err := postgres.GrantUnlimitedBalance(db, ctx, targetID); err != nil {
			logger.Log(gologger.LevelFatal, gologger.LogType("DB"), fmt.Sprintf("failed to grant unlimited balance: %v", err), "")
			return
		}
		logger.Log(gologger.LevelSuccess, gologger.LogType("CLI"), fmt.Sprintf("unlimited balance granted to %s", targetID.String()), "")
		return
	}
	if *revokeUnlimitedFlag != "" {
		targetID, err := uuid.Parse(*revokeUnlimitedFlag)
		if err != nil {
			logger.Log(gologger.LevelFatal, gologger.LogType("CLI"), fmt.Sprintf("invalid revoke-unlimited uuid: %v", err), "")
			return
		}
		if err := postgres.RevokeUnlimitedBalance(db, ctx, targetID); err != nil {
			logger.Log(gologger.LevelFatal, gologger.LogType("DB"), fmt.Sprintf("failed to revoke unlimited balance: %v", err), "")
			return
		}
		logger.Log(gologger.LevelSuccess, gologger.LogType("CLI"), fmt.Sprintf("unlimited balance revoked for %s", targetID.String()), "")
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
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS, echo.DELETE},
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
