package main

import (
	"fmt"
	"os"

	"github.com/GuilhermeDias02/CRM/cmd"
	"github.com/GuilhermeDias02/CRM/internal/contact"
)

func main() {
	store, err := contact.NewGormStore("contacts.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'initialisation du store: %v\n", err)
		os.Exit(1)
	}
	cmd.Execute(store)
}