package api

import (
	"api/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandle struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandle {
	return &HotelHandle{
		store: store,
	}
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandle) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	filter := bson.M{"hotelID": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrNoResourceNotFound("hotel")
	}
	return c.JSON(rooms)
}

func (h *HotelHandle) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidID()
	}

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), oid)
	if err != nil {
		return ErrNoResourceNotFound("hotel")
	}
	return c.JSON(hotel)
}

func (h *HotelHandle) HandleGetHotels(c *fiber.Ctx) error {
	// var queryParams HotelQueryParams
	// if err := c.QueryParser(&queryParams); err != nil {
	// 	return err
	// }

	// fmt.Println(queryParams)

	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return ErrNoResourceNotFound("hotels")
	}
	return c.JSON(hotels)
}
