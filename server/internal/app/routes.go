package app

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/controllers/metadata"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/controllers/secure"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/server/internal/controllers/users"

)

func (a *App) initializeRoutes() {

	users.AddUserRoutes(a.Router)
	metadata.AddMetadataRoutes(a.Router)
	secure.AddSecurityRoutes(a.Router)
	AddStaticRoutes(a.Router)
}



func AddStaticRoutes(a *mux.Router) {

	a.PathPrefix("/").HandlerFunc(serveIndex)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
