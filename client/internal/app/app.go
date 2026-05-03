package app

import "github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"

type App struct {
	
}

func (a *App) Start() {

	session.Connect()

	


}
