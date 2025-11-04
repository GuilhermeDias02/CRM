package app

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/GuilhermeDias02/CRM/internal/action"
	"github.com/GuilhermeDias02/CRM/internal/contact"
)

func App(store contact.Storer) {
	if handled, code := action.HandleFlags(store, os.Args[1:], os.Stdout, os.Stderr); handled {
		if code != 0 {
			os.Exit(code)
		}

		action.DisplayContacts(store)
		return
	}

	for {
		action.DisplayContacts(store)

		fmt.Println("\n1 - Ajouter un contact ")
		fmt.Println("2 - Supprimer un contact ")
		fmt.Println("3 - Mettre à jour un contact ")
		fmt.Println("4 - Quitter l'application ")
		fmt.Print("Veuillez choisir une action: ")

		reader := bufio.NewReader(os.Stdin)
		actionStr, _ := reader.ReadString('\n')
		actionStr = strings.TrimSpace(actionStr)
		actionInt, err := strconv.Atoi(actionStr)

		if err != nil {
			fmt.Println("Action indisponible")
			continue
		}

		if !action.HandleAction(store, actionInt) {
			return
		}
	}
}