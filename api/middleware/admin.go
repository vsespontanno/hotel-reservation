package middleware

import (
	"fmt"
	"hotel-reservation/types"

	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("user not logged in")
	}
	if !user.IsAdmin {
		return fmt.Errorf("user is not admin")
	}
	return c.Next()
}
