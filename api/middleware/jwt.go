package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	tokenArr, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if len(tokenArr) == 0 {
		return fmt.Errorf("unauthorized")
	}
	token := tokenArr[0]
	claims, err := validateToken(token)
	if err != nil {
		return err
	}

	expires := int64(claims["expires"].(float64))

	if time.Now().Unix() > expires {
		return fmt.Errorf("token expired")
	}

	fmt.Println(expires)

	fmt.Println(claims)
	return c.Next()

}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("NEVER SHOW YOUR SECRET: ", secret)
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse JWT token:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		fmt.Println("Invalid token")
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil

}
