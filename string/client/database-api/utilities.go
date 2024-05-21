package database_api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type JSONError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func HandleError(w http.ResponseWriter, err error) {
	encodedErr, _ := json.Marshal(JSONError{true, err.Error()})
	// Main switch to handle different types of errors
	var status int
	switch true {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case errors.Is(err, gorm.ErrInvalidData) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrInvalidTransaction) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrInvalidValue) ||
		errors.Is(err, gorm.ErrDuplicatedKey) ||
		errors.Is(err, gorm.ErrInvalidValueOfLength) ||
		errors.Is(err, &url.Error{}) ||
		errors.As(err, &sqlite3.Error{}) ||
		errors.Is(err, &json.SyntaxError{}) ||
		errors.Is(err, &json.UnmarshalTypeError{}) ||
		errors.Is(err, &json.UnsupportedValueError{}) ||
		errors.Is(err, &json.InvalidUnmarshalError{}) ||
		errors.Is(err, &json.MarshalerError{}) ||
		errors.Is(err, &json.UnsupportedTypeError{}) ||
		uuid.IsInvalidLengthError(err) ||
		err.Error() == "id is required" ||
		err.Error() == "invalid UUID format":
		status = http.StatusBadRequest
	case err.Error() == "method not allowed":
		status = http.StatusMethodNotAllowed
	case err.Error() == "ownUser already exists":
		status = http.StatusConflict

	default:
		status = http.StatusInternalServerError
	}
	http.Error(w, string(encodedErr), status)
}
