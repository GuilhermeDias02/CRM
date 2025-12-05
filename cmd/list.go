package cmd

import (
	"fmt"

	"github.com/GuilhermeDias02/CRM/internal/contact"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lister tous les contacts",
	Long:  `La commande 'list' affiche tous les contacts enregistrés dans le CRM.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nVos contacts: ")
		contacts := contact.GetContacts(store)
		if len(contacts) == 0 {
			fmt.Println("\tAucun contact enregistré.")
			return
		}
		for index, val := range contacts {
			fmt.Printf("\tId: %d, Name: %s, Email: %s\n", index, val.Name, val.Email)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}


