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

func DisplayContacts(store contact.Storer) {
	fmt.Println("\nVos contacts: ")
	for index, val := range contact.GetContacts(store) {
		fmt.Printf("\tId: %d, Name: %s, Email: %s\n", index, val.Name, val.Email)
	}
}

func addContactForm(store contact.Storer) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nAjouter un contact: ")

	fmt.Print("\tQuel est son nom: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("\tQuel est son email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	saved, err := contact.NewContact(store, name, email)

	if err != nil {
		fmt.Printf("\nErreur à la création d'un nouveau contact: %s", err)
		return
	}

	fmt.Printf("Contact ajouté avec ID %d\n", saved.Id)
}

func deleteContactForm(store contact.Storer) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nSupprimer un contact: ")

	fmt.Print("\tQuel est son id: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	idParsed, err := strconv.Atoi(idStr)
	idInt := uint(idParsed)

	if err != nil {
		fmt.Println("\nL'id selectionné n'est pas un numéro")
		return
	}

	if err := contact.DeleteContact(store, idInt); err != nil {
		fmt.Printf("\nCe contact n'existe pas: %s", err)
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

func HandleAction(store contact.Storer, actionInt int) bool {
	switch actionInt {
	case 1:
		addContactForm(store)
		return true
	case 2:
		deleteContactForm(store)
		return true
	case 3:
		updateContactForm(store)
		return true
	case 4:
		fmt.Println("Au revoir !")
		return false
	default:
		fmt.Println("Action indisponible")
		return true
	}
}

func HandleFlags(store contact.Storer, args []string, out io.Writer, errW io.Writer) (bool, int) {
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
	
	saved, err := contact.NewContact(store, trimName, trimEmail)
	if err != nil {
		fmt.Fprintf(errW, "\nErreur à la création de l'utilisateur: %s", err)
		return true, 1
	}

	fmt.Fprintf(out, "Contact ajouté avec ID %d\n", saved.Id)
	return true, 0
}

func updateContactForm(store contact.Storer) {
	idStr := strings.TrimSpace(readLine("ID du contact à mettre à jour: "))
	idParsed, err := strconv.Atoi(idStr)
	idInt := uint(idParsed)
	if err != nil || idInt <= 0 {
		fmt.Println("ID invalide")
		return
	}

	name := strings.TrimSpace(readLine("Nouveau nom (laisser vide pour garder): "))
	email := strings.TrimSpace(readLine("Nouvel email (laisser vide pour garder): "))

	if errUpdate := contact.UpdateContact(store, idInt, name, email); errUpdate != nil {
		fmt.Printf("Erreur lors de la mise à jour du contact: %s", errUpdate)
		return
	}
	fmt.Println("Contact mis à jour.")
}
