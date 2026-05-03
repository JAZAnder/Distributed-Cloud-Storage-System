package menus

import (
	"bufio"
	"fmt"
	"os"
	"strings"

)

func update() {
	ClearScreen()
	scanner := bufio.NewScanner(os.Stdin)
	for loop := true; loop; {
		message := `
		
		Currently connected to Coordinator: ` + coordinatorURL + `
		Currently connected to Node: ` + nodeURL + `
		
		- - - Distributed-Cloud-Storage-System - - - 
			      - - - Update Menu - - -
		

			This menu has not be implemented yet.

		Pick an Option (exit): `

		fmt.Print(message)
		scanner.Scan()
		option := strings.ToLower(strings.TrimSpace(scanner.Text())) 

		if option == "1" || option == "download" {


		} else if option == "6" || option == "close" || option == "exit" {
			loop = false
		} 

	}

}
