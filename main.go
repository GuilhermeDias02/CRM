package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/GuilhermeDias02/CRM/contact"
)

func main() {
	for {
		//Lister les contacts actuels
		fmt.Println("\nVos contacts: ")
		for index, val := range contact.GetContacts() {
			fmt.Printf("\tId: %d, Name: %s, Email: %s\n", index, val.Name, val.Email)
		}

		//Ajouter un contact
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\nAjouter un contact: ")

		fmt.Print("\tQuel est son nom: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Print("\tQuel est son email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		contact.AddContact(name, email)
	}
}
