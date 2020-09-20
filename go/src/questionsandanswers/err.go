package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var ErrNoSuchQuestion = errors.New("No such question")

type HTTPError struct {
	Status int
	Message string
	err error
}

func SendError(w http.ResponseWriter, err error) {
	herr, isa := err.(*HTTPError)
	if !isa {
		herr = InternalServerError(err)
	}
	herr.WriteError(w)
}

func (herr *HTTPError) getMessage() string {
	if herr.err != nil {
		return herr.err.Error()
	}
	return herr.Message
}

func (herr *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", herr.Status, herr.getMessage())
}

func (herr *HTTPError) WriteError(w http.ResponseWriter) {
	if herr.err != nil {
		log.Println("ERROR:", herr.err)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(herr.Status)
	w.Write([]byte(herr.Message))
}

func BadRequest(err error) *HTTPError {
	return &HTTPError{
		Status: http.StatusBadRequest,
		Message: "Invalid request",
		err: err,
	}
}

func NotFound(err error) *HTTPError {
	return &HTTPError{
		Status: http.StatusNotFound,
		Message: "Resource not found",
		err: err,
	}
}

func InternalServerError(err error) *HTTPError {
	return &HTTPError{
		Status: http.StatusInternalServerError,
		Message: "Internal Server Error",
		err: err,
	}
}

func BadGateway(err error) *HTTPError {
	return &HTTPError{
		Status: http.StatusBadGateway,
		Message: "Downstream server error",
		err: err,
	}
}
