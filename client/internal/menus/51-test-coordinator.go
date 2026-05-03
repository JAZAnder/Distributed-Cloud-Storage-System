package menus

import (
	"bufio"
	"fmt"
	"os"

	"github.com/JAZAnder/Distributed-Cloud-Storage-System/client/internal/logic/httpRequests"

)

func testCoordinator() {
	scanner := bufio.NewScanner(os.Stdin)
	ClearScreen()
	body, err := httpRequests.CoordinatorRequests("GET", "/api/whoami", "")
	message := ``
	if err != nil {
		message = "\n\n Error:" + err.Error() + "\n\n"
	} else {

		message = "\n\n\t"+string(body)+"\n\n"
	}

	fmt.Println(message)
	fmt.Println("\nPress Enter to Continue: ")
	scanner.Scan()

}
