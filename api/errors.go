package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(*Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)

}

type Error struct {
	Code    int    "json:code"
	Message string "json:message"
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func ErrResourceNotFound(res string) *Error {
	return NewError(http.StatusNotFound, res+" resource not found")
}

func ErrUnAuthorized() *Error {
	return NewError(http.StatusUnauthorized, "unathorized")
}

func ErrBadRequest() *Error {
	return NewError(http.StatusBadRequest, "invalid JSON request")
}

func ErrInvalidID() *Error {
	return NewError(http.StatusBadRequest, "invalid id given")
}
