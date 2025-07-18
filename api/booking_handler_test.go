package api

import (
	"encoding/json"
	"fmt"
	"hotel-reservation/db/fixtures"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestUserGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "tony", "soprano", false)
		user           = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)

	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 statuscode but got %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookingResp)
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected booking id of %s, got %s", booking.ID, bookingResp.ID)
	}

	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected booking user id of %s, got %s", booking.UserID, bookingResp.UserID)
	}

	// request by non-auth user

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
		user      = fixtures.AddUser(db.Store, "james", "foo", false)
		adminUser = fixtures.AddUser(db.Store, "admin", "admin", true)
		hotel     = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from      = time.Now()
		till      = from.AddDate(0, 0, 5)
		booking   = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app       = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)

	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)

	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of %d, got %d", http.StatusOK, resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected booking id of %s, got %s", booking.ID, have.ID)
	}

	if have.UserID != booking.UserID {
		t.Fatalf("expected booking user id of %s, got %s", booking.UserID, have.UserID)
	}

	// test if non-admin can access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)

	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized code got %d", resp.StatusCode)
	}

}
