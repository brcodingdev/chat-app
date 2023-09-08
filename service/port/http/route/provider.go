package route

import "github.com/gorilla/mux"

// RegisterRoute register /v1 routes
func RegisterRoute() *mux.Router {
	return mux.NewRouter().PathPrefix("/v1").Subrouter()
}
