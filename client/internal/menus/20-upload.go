package menus

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/diskOperations"

)

func upload() {
	ClearScreen()
	scanner := bufio.NewScanner(os.Stdin)
	for loop := true; loop; {
		message := `
		
		Currently connected to Coordinator: ` + coordinatorURL + `
		Currently connected to Download Node: ` + downloadNodeURL + `
		Currently connected to Upload Node: ` + uploadNodeURL + `
		
		- - - Distributed-Cloud-Storage-System - - - 
			      - - - Upload Menu - - -
		`
		fmt.Print(message)
		files, err := diskOperations.ListFilesForUpload()
		if err != nil {
			fmt.Println(err.Error())
		}
		message = `		Exit to return to previous menu
		
		Pick an Option (1/2/3/4/exit): `

		fmt.Print(message)
		scanner.Scan()
		option := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if option == "close" || option == "exit" {
			loop = false
		} else {
			fileId, err := strconv.Atoi(option)
			if err != nil {
				println(err.Error())
			}
			fileName := files[fileId]

			encryptAndUpload(fileName)


		}

	}

}
