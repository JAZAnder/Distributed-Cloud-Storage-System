package metadata

import "github.com/gorilla/mux"

func AddMetadataRoutes(a *mux.Router) {

	a.HandleFunc("/api/metadata", getMyMetadata).Methods("GET")
	a.HandleFunc("/api/metadata/{id:[0-9]+", getMetadataById).Methods("GET")
	a.HandleFunc("/api/metadata", uploadFile).Methods("POST")

}