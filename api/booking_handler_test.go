package api

import (
	"api/fixture"
	"api/middleware"
	"api/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		nonAuthUser    = fixture.AddUser(db.Store, "jimmy", "carr", false)
		user           = fixture.AddUser(db.Store, "foo", "buz", false)
		hotel          = fixture.AddHotel(db.Store, "hotel 1", "a", 4, nil)
		room           = fixture.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixture.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	// app.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var bookingResponse *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResponse); err != nil {
		t.Fatal(err)
	}

	if bookingResponse.ID != booking.ID {
		t.Fatalf("expected %s got %s ", booking.ID, bookingResponse.ID)
	}

	if bookingResponse.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResponse.UserID)
	}
	fmt.Println(bookingResponse)

	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}

}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		user      = fixture.AddUser(db.Store, "foo", "buz", false)
		adminUser = fixture.AddUser(db.Store, "admin", "admin", true)
		hotel     = fixture.AddHotel(db.Store, "hotel 1", "a", 4, nil)
		room      = fixture.AddRoom(db.Store, "small", true, 4.4, hotel.ID)

		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 5)
		booking        = fixture.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatal("not 200 response, get", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}

	have := bookings[0]

	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s ", booking.ID, have.ID)
	}

	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}

	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}

}
