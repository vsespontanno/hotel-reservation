package api

import (
	"fmt"
	"hotel-reservation/db"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenArr, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrUnAuthorized()
		}
		if len(tokenArr) == 0 {
			return ErrUnAuthorized()
		}
		token := tokenArr[0]
		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return ErrUnAuthorized()
		}

		//set the current authenticated	user to the context.
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse JWT token:", err)
		return nil, ErrUnAuthorized()
	}

	if !token.Valid {
		fmt.Println("Invalid token")
		return nil, ErrUnAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, ErrUnAuthorized()
	}
	return claims, nil

}
