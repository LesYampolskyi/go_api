package api

import (
	"api/db"
	"fmt"

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
	return c.JSON(booking)
}