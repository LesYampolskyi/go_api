package fixture

import (
	"api/db"
	"api/types"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hID,
	}

	insertRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertRoom
}

func AddHotel(store *db.Store, name, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomsIDS = rooms
	if rooms == nil {
		roomsIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomsIDS,
		Rating:   rating,
	}

	insertedHotel, err := store.Hotel.Insert(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel

}

func AddUser(store *db.Store, fn, ln string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
