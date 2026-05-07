package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

)

type App struct {
	Router *mux.Router
}

func (a *App) Initialize() {

	
	//Creates Router
	a.Router = mux.NewRouter()
	a.initializeRoutes()

}

func (a *App) Run(addr string) {
	err := http.ListenAndServeTLS(":"+addr, "./certs/certificate.cer", "./certs/privateKey.key", a.Router)
	log.Println(err)
}
