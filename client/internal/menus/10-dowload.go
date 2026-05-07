package menus

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/httpRequests"

)

func download() {
	ClearScreen()
	scanner := bufio.NewScanner(os.Stdin)

	for loop := true; loop; {
		listOfFiles, err := httpRequests.GetDirectory()
		if err != nil {
			message := "\n\n\t Error: " + err.Error() + "\n\n Press Enter to return to previous menu (Enter): "
			fmt.Print(message)
			scanner.Scan()
			return
		}

		message := `
		
		Currently connected to Coordinator: ` + coordinatorURL + `
		Currently connected to Download Node: ` + downloadNodeURL + `
		Currently connected to Upload Node: ` + uploadNodeURL + `
		
		- - - Distributed-Cloud-Storage-System - - - 
			      - - - Download Menu - - -
			`

		message += `		
		--- Available Files (Select a number to download) ---
				Choice | File Name
		`

		for _, f := range listOfFiles {
			message += `		[` + strconv.Itoa(int(f.ID)) + `]     ` + f.Name + "\n"

		}

		message += `


		Pick an File to download (exit to return): `

		fmt.Print(message)
		scanner.Scan()
		option := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if option == "exit" || option == "close" {
			loop = false
		} else {
			fileId, err := strconv.Atoi(option)
			if err != nil {
				fmt.Println(err.Error() + "\n\n\n Press Enter to continue: ")
				scanner.Scan()
			} else {
				err := downloadAndDecode(fileId)
				if err != nil {
					fmt.Println(err.Error() + "\n\n\n Press Enter to continue: ")
					scanner.Scan()
				}else{
					fmt.Println("Operation Completed Successfully! \n Press Enter to continue: ")
					scanner.Scan()
				}

			}

		}

	}

}
