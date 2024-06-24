package webserver

import (
	"mime"
	"net/http"
)

const (
	nonJsonContentType string = "Content-Type must be JSON"
	unknownContentType string = "Malformed Content-Type header"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func AddMiddleware(h http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	totalMiddleware := len(middleware) - 1

	if totalMiddleware == -1 {
		return h
	}

	for i := totalMiddleware; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}

// TODO: fully flesh this out later
func ValidateJSONMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mediaType, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, unknownContentType, http.StatusBadRequest)
				return
			}

			if mediaType != "application/json" {
				http.Error(w, nonJsonContentType, http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
