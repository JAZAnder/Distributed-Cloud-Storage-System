package menus

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"

)

var coordinatorURL string
var downloadNodeURL string
var uploadNodeURL string
func Home() {
	currentSession := session.GetSession()
	coordinatorURL = currentSession.Coordinator.CoordinatorURL
	downloadNodeURL = currentSession.Node.DownloadloadNodeURL 
	uploadNodeURL = currentSession.Node.UploadNodeURL

	ClearScreen()

	scanner := bufio.NewScanner(os.Stdin)
	
	for loop := true; loop; {
		message := `
		
		Currently connected to Coordinator: ` + coordinatorURL + `
		Currently connected to Download Node: ` + downloadNodeURL + `
		Currently connected to Upload Node: ` + uploadNodeURL + `
		
		- - - Distributed-Cloud-Storage-System - - - 
			      - - - Main Menu - - -
		
		1)Download File via DCSS
		
		2)Upload File via DCSS

		3)Update File via DCSS

		4)Download File via CID

		5)View Connection Details

		6)Close Connection and Exit
		
		Pick an Option (1/2/3): `

		fmt.Print(message)
		scanner.Scan()
		option := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if option == "1" || option == "download" {
			download()

		} else if option == "2" || option == "upload" {
			upload()

		} else if option == "3" || option == "update" {
			update()

		} else if option == "4" || option == "cid" {
			cid()

		} else if option == "5" || option == "view" || option == "details" {
			view()

		} else if option == "6" || option == "close" || option == "exit" {
			loop = false
		}

	}

}

func ClearScreen() {
	var cmd *exec.Cmd
	// Identify the operating system to choose the correct command [Conversation History]
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
