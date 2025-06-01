package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func getJWTSecret() []byte {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return []byte(s)
	}
	return []byte("your-256-bit-secret")
}

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:    getJWTSecret(),
		SigningMethod: "HS256",
		ContextKey:    "user",
		TokenLookup:   "header:Authorization",
		ErrorHandler: func(c echo.Context, err error) error {
			auth := c.Request().Header.Get("Authorization")
			fmt.Println("🔍 Authorization header:", auth)
			fmt.Println("🔐 JWT validation error:", err)
			return echo.NewHTTPError(401, "invalid or expired jwt")
		},
	})
}

func UserIDFromContext(c echo.Context) (uint, error) {
	user := c.Get("user")
	if user == nil {
		return 0, errors.New("no token in context")
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return 0, errors.New("invalid token format")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}
	uidf, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id claim missing or invalid")
	}
	return uint(uidf), nil
}

func RoleFromContext(c echo.Context) (string, error) {
	user := c.Get("user")
	if user == nil {
		return "", errors.New("no token in context")
	}
	token, ok := user.(*jwt.Token)
	if !ok {
		return "", errors.New("invalid token format")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}
	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role claim missing or invalid")
	}
	return role, nil
}
