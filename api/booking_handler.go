package api

import (
	"api/db"
	"fmt"
	"net/http"

	// "fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelingBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}

	return c.JSON(genericResp{Type: "msg", Msg: "updated"})

}

// TODO: admin
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	fmt.Println("error here0")
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		fmt.Println("error here1")
		return err
	}
	return c.JSON(bookings)
}

// TODO: user
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(c)
	if err != nil {
		return err
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	return c.JSON(booking)
}
