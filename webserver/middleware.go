package webserver

import (
	"mime"
	"net/http"
	"time"
)

const (
	nonJsonContentType string = "Content-Type must be JSON"
	unauthorized       string = "Unauthorized"
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

const responseWaitDuration = time.Second * 2 // TODO: pull this from config, set on conf struct
// AuthorizeMiddleware requires being used over HTTPS connections only. It uses
// basic authentication to check a user and password combination.  It also
// guards against timing attacks by guaranteeing the middleware response will
// always take the same amount of configured time.
func AuthorizeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pw, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, unauthorized, http.StatusUnauthorized)
			return
		}

		validity := make(chan bool, 1)
		now := time.Now()

		go isValidLogin(validity, user, pw)
		// TODO (REN-145): add fingerprint whitelist check, update channel count

		isValid := <-validity

		elapsed := time.Since(now)
		remaining := responseWaitDuration - elapsed
		totalWaits := 0
		if remaining > 0 {
			totalWaits++
			time.Sleep(remaining)
		}

		if totalWaits == 0 {
			// TODO: error log this, as it means the response wait duration needs to be increased
		}

		if !isValid {
			http.Error(w, unauthorized, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
		return
	}
}

// TODO: pull from storage, make this a receiver on a conf struct
func isValidLogin(isValid chan bool, user string, password string) {
	actualUser := "getFromStorage"
	actualPw := "getFromStorage"
	isValid <- (user == actualUser) && (password == actualPw)
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
