package users

import "github.com/gorilla/mux"

func AddUserRoutes(a *mux.Router) {

	a.HandleFunc("/api/login", login).Methods("POST")

}