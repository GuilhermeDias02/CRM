package cmd

import (
	"fmt"
	"os"

	"github.com/GuilhermeDias02/CRM/internal/contact"
	"github.com/spf13/cobra"
)

var store contact.Storer

var rootCmd = &cobra.Command{
	Use:   "crm",
	Short: "CRM est un outil de gestion de contacts en ligne de commande",
	Long: `CRM est un mini-CRM (Customer Relationship Management) accessible en ligne de commande,
permettant de gérer une liste de contacts avec les opérations CRUD de base.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Par défaut, afficher les contacts si aucune commande n'est spécifiée
		displayContacts()
	},
}

func Execute(storer contact.Storer) {
	store = storer
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func displayContacts() {
	fmt.Println("\nVos contacts: ")
	contacts := contact.GetContacts(store)
	if len(contacts) == 0 {
		fmt.Println("\tAucun contact enregistré.")
		return
	}
	for index, val := range contacts {
		fmt.Printf("\tId: %d, Name: %s, Email: %s\n", index, val.Name, val.Email)
	}
}


