package menus

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Home() {
	scanner := bufio.NewScanner(os.Stdin)
	for loop := true; loop; {
		loop = false
		message := `Pick an Option (1/2/3): `

		fmt.Print(message)
		scanner.Scan()
		option := strings.TrimSpace(scanner.Text())

		if option == "1" || option == "one" {

		} else if option == "2" || option == "two" {

		} else if option == "3" || option == "three" {

		} else {
			loop = true
		}

	}

}
