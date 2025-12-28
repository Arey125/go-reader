package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func ServerError(w http.ResponseWriter, err error) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
    panic(err)
}

func Forbiden(w http.ResponseWriter) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func HttpError(w http.ResponseWriter, statusCode int) {
	fmt.Printf("http error %s, stack trace:\n%s", http.StatusText(statusCode), debug.Stack())
	http.Error(w, http.StatusText(statusCode), statusCode)
}
