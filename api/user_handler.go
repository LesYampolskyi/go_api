package api

import (
	"api/db"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Printf("STEP1")
	ctx := c.Context()
	user, err := h.userStore.GetUserByID(ctx, id)
	fmt.Println(user)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
