package menus

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func view() {
	ClearScreen()
	scanner := bufio.NewScanner(os.Stdin)
	for loop := true; loop; {
		message := `
		
		Currently connected to Coordinator: ` + coordinatorURL + `
		Currently connected to Download Node: ` + downloadNodeURL + `
		Currently connected to Upload Node: ` + uploadNodeURL + `
		
		- - - Distributed-Cloud-Storage-System - - - 
			      - - - Configuration Menu - - -
		

			1) Test Coordinator Connection

		Pick an Option (exit): `

		fmt.Print(message)
		scanner.Scan()
		option := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if option == "1" {
			testCoordinator()

		} else if option == "6" || option == "close" || option == "exit" {
			loop = false
		}

	}

}

type whoAmIResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
