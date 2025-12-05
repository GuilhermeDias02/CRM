package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/GuilhermeDias02/CRM/internal/contact"
	"github.com/spf13/cobra"
)

var (
	addName  string
	addEmail string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Ajouter un nouveau contact",
	Long: `La commande 'add' permet d'ajouter un nouveau contact au CRM.
Si les flags --name et --email ne sont pas fournis, le mode interactif sera activé.`,
	Run: func(cmd *cobra.Command, args []string) {
		var name, email string

		// Si les flags sont fournis, utiliser les valeurs des flags
		if addName != "" && addEmail != "" {
			name = strings.TrimSpace(addName)
			email = strings.TrimSpace(addEmail)
		} else {
			// Mode interactif
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("\nAjouter un contact: ")

			fmt.Print("\tQuel est son nom: ")
			nameInput, _ := reader.ReadString('\n')
			name = strings.TrimSpace(nameInput)

			fmt.Print("\tQuel est son email: ")
			emailInput, _ := reader.ReadString('\n')
			email = strings.TrimSpace(emailInput)
		}

		saved, err := contact.NewContact(store, name, email)
		if err != nil {
			fmt.Printf("\nErreur à la création d'un nouveau contact: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Contact ajouté avec ID %d\n", saved.Id)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Définition des drapeaux pour la commande 'add' (optionnels)
	addCmd.Flags().StringVarP(&addName, "name", "n", "", "Nom du contact")
	addCmd.Flags().StringVarP(&addEmail, "email", "e", "", "Email du contact")
}


