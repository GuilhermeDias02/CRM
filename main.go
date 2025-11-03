package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/GuilhermeDias02/CRM/action"
)

func main() {
	for {
		action.DisplayContacts()

		//Action à exécuter
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\n1 - Ajouter un contact ")
		fmt.Println("2 - Supprimer un contact ")
		fmt.Print("Veuillez choisir une action: ")
		actionStr, _ := reader.ReadString('\n')
		actionStr = strings.TrimSpace(actionStr)
		actionInt, err := strconv.Atoi(actionStr)

		if err != nil {
			fmt.Println("Action indisponible")
			continue
		}

		switch actionInt {
		case 1:
			action.AddContactForm()
		case 2:
			action.DeleteContactForm()
		default:
			fmt.Println("Action indisponible")
		}
	}
}
