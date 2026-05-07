package menus

import (
	"bufio"
	"fmt"
	"os"
	"strings"

)

func cid() {
	ClearScreen()
	scanner := bufio.NewScanner(os.Stdin)
	for loop := true; loop; {
		message := `
		
		Currently connected to Coordinator: ` + coordinatorURL + `
		Currently connected to Download Node: ` + downloadNodeURL + `
		Currently connected to Upload Node: ` + uploadNodeURL + `
		
		- - - Distributed-Cloud-Storage-System - - - 
			      - - - Manual Download Menu - - -
		

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
