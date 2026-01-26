package middleware

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v3/jwk"
	gologger "github.com/nrf24l01/go-logger"
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

			traceID := ""
			if v := c.Get("traceId"); v != nil {
				if s, ok := v.(string); ok {
					traceID = s
				}
			}

			// Jwt parse with keyfunc (resolve key from JWKS cache)
			keyFunc := func(token *jwt.Token) (interface{}, error) {
				kid, ok := token.Header["kid"].(string)
				if !ok || kid == "" {
					return nil, fmt.Errorf("missing kid header")
				}

				set, err := h.Jwks.Lookup(context.Background(), h.Config.KeyCloakConfig.URL)
				if err != nil {
					return nil, err
				}

				key, ok := set.LookupKeyID(kid)
				if !ok {
					return nil, fmt.Errorf("unable to find JWK for kid %s", kid)
				}

				pub, err := key.PublicKey()
				if err != nil {
					return nil, err
				}

				raw, err := jwk.PublicRawKeyOf(pub)
				if err != nil {
					return nil, err
				}

				switch k := raw.(type) {
				case *rsa.PublicKey:
					return k, nil
				case rsa.PublicKey:
					return &k, nil
				case *rsa.PrivateKey:
					return &k.PublicKey, nil
				case rsa.PrivateKey:
					return &k.PublicKey, nil
				case *ecdsa.PublicKey:
					return k, nil
				case ecdsa.PublicKey:
					return &k, nil
				case *ecdsa.PrivateKey:
					return &k.PublicKey, nil
				case ecdsa.PrivateKey:
					return &k.PublicKey, nil
				default:
					return nil, fmt.Errorf("unsupported public key type: %T", raw)
				}
			}

			token, err := jwt.Parse(tokenString, keyFunc)
			if err != nil {
				h.Logger.Log(gologger.LevelError, gologger.LogType("AUTH"), fmt.Sprintf("Failed to parse token: %v", err), traceID)
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "invalid token", nil))
			}

			// Check keys
			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "SUS token", nil))
			}

			// Load claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "invalid token claims", nil))
			}

			// Извлекаем user_id
			userID, ok := claims["sub"].(string)
			userUUID, err := uuid.Parse(userID)
			if err != nil {
				h.Logger.Log(gologger.LevelError, gologger.LogType("AUTH"), fmt.Sprintf("Invalid user ID in claims: %v", err), traceID)
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "invalid user ID in claims", nil))
			}
			if !ok {
				return c.JSON(http.StatusUnauthorized, echokitSchemas.GenError(c, echokitSchemas.UNAUTHORIZED, "wrong claims", nil))
			}

			// Передаем user_id в контекст
			c.Set("userID", userUUID)

			return next(c)
		}
	}
}
