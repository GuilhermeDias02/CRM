package action

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/GuilhermeDias02/CRM/internal/contact"
)

func DisplayContacts() {
	//Lister les contacts actuels
	fmt.Println("\nVos contacts: ")
	for index, val := range *contact.GetContacts() {
		fmt.Printf("\tId: %d, Name: %s, Email: %s\n", index, val.Name, val.Email)
	}
}

func addContactForm() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nAjouter un contact: ")

	fmt.Print("\tQuel est son nom: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("\tQuel est son email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	newContact, err := contact.NewContact(name, email)

	if err != nil {
		fmt.Printf("\nErreur à la création d'un nouveau contact: %s", err)
	}

	newContact.AddContact()
}

func deleteContactForm() {
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

	contactToDelete, err := contact.GetContactById(idInt)

	if err != nil {
		fmt.Printf("\nCe contact n'existe pas: %s", err)
	}

	if contactErr := contactToDelete.DeleteContact(); contactErr != nil {
		fmt.Printf("\nCe contact n'existe pas: %s", contactErr)
	} else {
		fmt.Println("Contact supprimé")
	}
}

func readLine(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de lecture:", err)
		return strings.TrimSpace(text)
	}
	return strings.TrimSpace(text)
}

func HandleAction(actionInt int) bool {
	switch actionInt {
	case 1:
		addContactForm()
		return true
	case 2:
		deleteContactForm()
		return true
	case 3:
		updateContactForm()
		return true
	case 4:
		fmt.Println("Au revoir !")
		return false
	default:
		fmt.Println("Action indisponible")
		return true
	}
}

//bool true and code 0 if the contact was created
func HandleFlags(args []string, out io.Writer, errW io.Writer) (bool, int) {
	fs := flag.NewFlagSet("crm", flag.ContinueOnError)
	fs.SetOutput(errW)
	add := fs.Bool("add", false, "Ajouter un contact via flags")
	name := fs.String("name", "", "Nom du contact")
	email := fs.String("email", "", "Email du contact")
	if err := fs.Parse(args); err != nil {
		return true, 1
	}
	if !*add {
		return false, 0
	}
	trimName := strings.TrimSpace(*name)
	trimEmail := strings.TrimSpace(*email)
	if trimName == "" || trimEmail == "" {
		fmt.Fprintln(errW, "Usage: -add -name \"Nom\" -email \"Email\"")
		return true, 1
	}
	
	newContact, err := contact.NewContact(trimName, trimEmail)
	if err != nil {
		fmt.Printf("\n Erreur à la création de l'utilisateur: %s", err)
		return true, 1
	}

	newContact.AddContact();
	fmt.Fprintf(out, "Contact ajouté avec ID %d\n", newContact.Id)
	return true, 0
}

func updateContactForm() {
	idStr := strings.TrimSpace(readLine("ID du contact à mettre à jour: "))
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt <= 0 {
		fmt.Println("ID invalide")
		return
	}

	contactToUpdate, errContact := contact.GetContactById(idInt)

	if errContact != nil {
		fmt.Printf("Contact introuvable: %s", errContact)
		return
	}

	name := strings.TrimSpace(readLine("Nouveau nom (laisser vide pour garder): "))
	email := strings.TrimSpace(readLine("Nouvel email (laisser vide pour garder): "))
	
	if errUpdate := contactToUpdate.UpdateContact(name, email); errUpdate != nil {
		fmt.Printf("Erreur lors de la mise à jour du contact: %s", errUpdate)
		return
	}
	fmt.Println("Contact mis à jour.")
}
