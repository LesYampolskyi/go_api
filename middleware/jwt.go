package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("JWT AUTH")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized x")
	}

	claims, err := parseToken(token)
	if err != nil {
		return err
	}

	fmt.Println(claims)
	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)

	if time.Now().Unix() > expires {
		return fmt.Errorf("token expired")
	}
	fmt.Println(expires)
	return c.Next()
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header)
			return nil, fmt.Errorf("unauthorized 1")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("SECRET::", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token: ", err)
		return nil, fmt.Errorf("unauthorized 2")
	}

	if !token.Valid {
		fmt.Println("invalid token: ", err)
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println(ok)
		return nil, fmt.Errorf("unauthorized 3")
	}
	fmt.Println("OK")
	return claims, nil
}
