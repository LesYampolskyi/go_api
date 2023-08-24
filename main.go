package main

import (
	"api/api"
	"api/db"
	"api/middleware"
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initialization
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		userHandler    = api.NewUserHandler(userStore)
		authHandler    = api.NewAuthHandler(userStore)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin = apiv1.Group("/admin", middleware.AdminAuth)
	)

	fmt.Println(client)
	listenAddr := flag.String("listenAddr", ":5000", "Listen address of API server")
	flag.Parse()

	if err != nil {
		panic(err)
	}
	// A handler should only do:
	// 	- serialization of incoming request (JSON)
	//  - do some data fetching from db
	//	- call some buisness logic
	//  - return the data back to the user

	// auth handlers
	auth.Post("/auth", authHandler.HandleAuthentication)

	// Versioned api route
	// user handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// admin handlers
	admin.Get("/bookings", bookingHandler.HandleGetBookings)
	// booking handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelingBooking)
	app.Listen(*listenAddr)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "working fine"})
}
