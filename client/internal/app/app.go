package app

import (
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"
	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/menus"
)

type App struct {
}

func (a *App) Start() {

	session.Connect()
	menus.Home()

}
