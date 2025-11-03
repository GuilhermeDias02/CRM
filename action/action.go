package action

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/GuilhermeDias02/CRM/contact"
)

func DisplayContacts() {
	//Lister les contacts actuels
	fmt.Println("\nVos contacts: ")
	for index, val := range contact.GetContacts() {
		fmt.Printf("\tId: %d, Name: %s, Email: %s\n", index, val.Name, val.Email)
	}
}

func AddContactForm() {
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

func DeleteContactForm() {
	//Ajouter un contact
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nSupprimer un contact: ")

	fmt.Print("\tQuel est son id: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		fmt.Println("\nL'id selectionné n'est pas un numéro")
		return
	}

	if contactErr := contact.DeleteContact(idInt); contactErr != nil {
		fmt.Printf("\nCe contact n'existe pas: %s", contactErr)
	}
}
