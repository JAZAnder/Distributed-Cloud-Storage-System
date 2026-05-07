package login

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/session"

)

func login() {

}

func checkNodeLogin() bool {
	scanner := bufio.NewScanner(os.Stdin)
	currentSession := session.GetSession()
	if currentSession.Node.DownloadloadNodeURL == "" {
		if os.Getenv("DOWNLOADNODE") == "" {
			fmt.Print("\nEnter the IPFS Node URL: ")
			scanner.Scan()
			currentSession.Coordinator.CoordinatorURL = strings.TrimSpace(scanner.Text())
			return true
		}else{
			currentSession.Coordinator.CoordinatorURL = os.Getenv("DOWNLOADNODE")
			return true
		}
		return false
	}

	return true
}
