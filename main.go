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

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\n1 - Ajouter un contact ")
		fmt.Println("2 - Supprimer un contact ")
		fmt.Println("5 - Mettre à jour un contact ")
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
