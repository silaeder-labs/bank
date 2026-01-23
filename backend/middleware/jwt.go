package middleware

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	echokitSchemas "github.com/nrf24l01/go-web-utils/echokit/schemas"
	"github.com/silaeder-labs/bank/backend/handlers"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(h *handlers.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "missing authorization header", nil))
			}

			// Remove bearer
			if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "invalid token format", nil))
			}
			tokenString := authHeader[7:]

			// Jwt parse
			token, err := jwt.Parse(tokenString, h.Jwks.Keyfunc)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "invalid token", nil))
			}

			// Check keys
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "invalid token", nil))
			}

			// Load claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "invalid token claims", nil))
			}

			// Verify standard claims
			if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "token expired", nil))
			}
			if !claims.VerifyIssuer(h.Config.KeyCloakConfig.ISSUER_URL, true) {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "invalid token issuer", nil))
			}

			// Извлекаем user_id
			userID, ok := claims["sub"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "Wrong claims", nil))
			}

			// Передаем user_id в контекст
			c.Set("userID", userID)

			return next(c)
		}
	}
}
