package api

import (
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// TODO: this needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

// TODO: this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return fmt.Errorf("this room is not booked")
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("internal server error")
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "unauthorized",
		})
	}
	return c.JSON(booking)
}
