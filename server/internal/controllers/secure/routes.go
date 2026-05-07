package secure

import "github.com/gorilla/mux"

func AddSecurityRoutes(a *mux.Router) {

	a.HandleFunc("/api/config", getPublicPrams).Methods("GET")
	a.HandleFunc("/api/key", mintUserKey).Methods("POST")

}