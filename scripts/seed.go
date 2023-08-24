package main

import (
	"api/api"
	"api/db"
	"api/fixture"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		Room:    db.NewMongoRoomStore(client, hotelStore),
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Hotel:   db.NewMongoHotelStore(client),
	}

	userNew := fixture.AddUser(store, "foo", "buz", false)
	fmt.Println("user -> ", api.CreateTokenFromUser(userNew))

	admin := fixture.AddUser(store, "admin", "admin", true)
	fmt.Println("admin -> ", api.CreateTokenFromUser(admin))

	hotelItem := fixture.AddHotel(store, "hotel 1", "location 1", 4, nil)

	room := fixture.AddRoom(store, "large", true, 88.44, hotelItem.ID)
	booking := fixture.AddBooking(store, userNew.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking -> ", booking)
	return
}
