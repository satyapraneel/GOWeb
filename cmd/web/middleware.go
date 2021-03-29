package main

import (
	"fmt"
	"net/http"
)

func secureHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("X-XSS-Protection", "1; mode=block")
		rw.Header().Set("X-Frame-Options", "deny")
		next.ServeHTTP(rw, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(rw, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				rw.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				app.serverError(rw, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(rw, r)
	})
}
