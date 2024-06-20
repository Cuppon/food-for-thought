package webserver

import "net/http"

type Route func(mux *http.ServeMux)

func NewServer(routes ...Route) http.Handler {
	mux := http.NewServeMux()

	for _, r := range routes {
		r(mux)
	}

	return mux
}
