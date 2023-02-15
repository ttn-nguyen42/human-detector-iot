package custom

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

/*
Equivalent to a 500 Internal Server Error, specifying in unknown error that occurs in server side.
Should be returned in data layers (repositories) and logic layers (services)
when an unknown error occurred.
*/
type InternalServerError struct {
	Message string `default:"Internal server error"`
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func NewInternalServerError(message string) error {
	err := &InternalServerError{}
	if len(message) != 0 {
		err.Message = message
	}
	// Logs the newly found error to the console
	logrus.Debug(err.Error())
	return err
}

/*
Equivalent to a 404 Not Found error, specifying failed attempt to find an item in the database
despite correct user input (typically an ID).
*/
type ItemNotFoundError struct {
	Message string `default:"Not found"`
}

func (e *ItemNotFoundError) Error() string {
	return e.Message
}

func NewItemNotFoundError(message string) error {
	err := &ItemNotFoundError{}
	if len(message) != 0 {
		err.Message = message
	}
	return err
}

/*
An use-case to a 400 Bad Request error, specifying a failed attempt to parse user request,
most likely be due to a missing field in the body.
*/
type FieldMissingError struct {
	Field string
}

func NewFieldMissingError(field string) error {
	return &FieldMissingError{
		Field: fmt.Sprintf("%s: required field is missing", field),
	}
}

func (e *FieldMissingError) Error() string {
	return e.Field
}

/*
An use-case to a 400 Bad Request error,
specifying a failed attempt to parse an incorrect ID from user
*/
type BadIdError struct {
	Message string `default:"Bad request"`
}

func NewBadIdError(message string) error {
	err := &BadIdError{}
	if len(message) != 0 {
		err.Message = message
	}
	return err
}

func (e *BadIdError) Error() string {
	return e.Message
}

/*
An use-case to a 400 Bad Request error,
specifying a failed attempt to parse an incorrect query from user
*/
type BadQueryError struct {
	Message string `default:"Bad request"`
}

func NewBadQueryError(message string) error {
	err := &BadIdError{}
	if len(message) != 0 {
		err.Message = message
	}
	return err
}

func (e *BadQueryError) Error() string {
	return e.Message
}

/*
An use-case to a 400 Bad Request, specifying a validation error on user request
*/
type InvalidFormatError struct {
	Message string `default:"Invalid payload format"`
}

func NewInvalidFormatError(message string) error {
	err := &InvalidFormatError{}
	if len(message) != 0 {
		err.Message = message
	}
	return err
}

func (e *InvalidFormatError) Error() string {
	return e.Message
}

/*
An use-case to 403 Unauthorized, specifying an unauthorized action
*/
type UnauthorizedError struct {
	Message string `default:"Unauthorized"`
}

func NewUnauthorizedError(message string) error {
	err := &UnauthorizedError{}
	if len(message) != 0 {
		err.Message = message
	}
	return err
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

/*
When a device attempts to register itself twice
*/
type AlreadyRegisteredError struct {
	Message string `default:"Already registered"`
}

func NewAlreadyRegisteredError(message string) error {
	err := &AlreadyRegisteredError{}
	if len(message) != 0 {
		err.Message = message
	}
	return err
}

func (e *AlreadyRegisteredError) Error() string {
	return e.Message
}
