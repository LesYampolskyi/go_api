package api

import (
	"api/types"

	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)

	if !ok {
		return nil, ErrUnAuthorized()
	}

	return user, nil
}
