package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/GuilhermeDias02/CRM/action"
)

func main() {
	if handled, code := handleFlags(os.Args[1:], os.Stdout, os.Stderr); handled {
		if code != 0 {
			os.Exit(code)
		}
		return
	}

	for {
		action.DisplayContacts()

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\n1 - Ajouter un contact ")
		fmt.Println("2 - Supprimer un contact ")
		fmt.Println("5 - Mettre à jour un contact ")
		fmt.Println("6 - Quitter l'application ")
		fmt.Print("Veuillez choisir une action: ")
		actionStr, _ := reader.ReadString('\n')
		actionStr = strings.TrimSpace(actionStr)
		actionInt, err := strconv.Atoi(actionStr)
		if err != nil {
			fmt.Println("Action indisponible")
			continue
		}
		if !handleAction(actionInt) {
			return
		}
	}
}

var addContactForm = action.AddContactForm
var deleteContactForm = action.DeleteContactForm
var updateContactFunc = updateContact

func handleAction(actionInt int) bool {
	switch actionInt {
	case 1:
		addContactForm()
		return true
	case 2:
		deleteContactForm()
		return true
	case 5:
		updateContactFunc()
		return true
	case 6:
		fmt.Println("Au revoir !")
		return false
	default:
		fmt.Println("Action indisponible")
		return true
	}
}

type Contact struct {
	ID int
	Name string
	Email string
}

var contacts = map[int]Contact{}
var nextID = 1

func addContact(name, email string) int {
	id := nextID
	nextID++
	contacts[id] = Contact{ID: id, Name: name, Email: email}
	return id
}

func updateContact() {
	idStr := strings.TrimSpace(readLine("ID du contact à mettre à jour: "))
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		fmt.Println("ID invalide")
		return
	}

	name := strings.TrimSpace(readLine("Nouveau nom (laisser vide pour garder): "))
	email := strings.TrimSpace(readLine("Nouvel email (laisser vide pour garder): "))

	var namePtr *string
	var emailPtr *string
	if name != "" {
		namePtr = &name
	}
	if email != "" {
		emailPtr = &email
	}
	if err := updateContactByID(id, namePtr, emailPtr); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Contact mis à jour.")
}

func updateContactByID(id int, namePtr *string, emailPtr *string) error {
	if id <= 0 {
		return fmt.Errorf("ID invalide")
	}
	c, ok := contacts[id]
	if !ok {
		return fmt.Errorf("Aucun contact avec cet ID")
	}
	if namePtr != nil {
		c.Name = strings.TrimSpace(*namePtr)
	}
	if emailPtr != nil {
		email := strings.TrimSpace(*emailPtr)
		if !isValidEmail(email) {
			return fmt.Errorf("Email invalide")
		}
		c.Email = email
	}
	contacts[id] = c
	return nil
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

func isValidEmail(email string) bool {
	return strings.Contains(email, "@")
}

func handleFlags(args []string, out io.Writer, errW io.Writer) (bool, int) {
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
	if !isValidEmail(trimEmail) {
		fmt.Fprintln(errW, "Email invalide")
		return true, 1
	}
	id := addContact(trimName, trimEmail)
	fmt.Fprintf(out, "Contact ajouté avec ID %d\n", id)
	return true, 0
}
