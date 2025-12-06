package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/GuilhermeDias02/CRM/internal/contact"
	"github.com/spf13/cobra"
)

var (
	updateID    uint
	updateName  string
	updateEmail string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Mettre à jour un contact existant",
	Long: `La commande 'update' permet de modifier le nom et/ou l'email d'un contact existant.
L'ID du contact est obligatoire. Si les flags --name et --email ne sont pas fournis,
le mode interactif sera activé.`,
	Run: func(cmd *cobra.Command, args []string) {
		var id uint
		var name, email string

		// Récupérer l'ID
		if updateID > 0 {
			id = updateID
		} else {
			// Mode interactif pour l'ID
			idStr := strings.TrimSpace(readLine("ID du contact à mettre à jour: "))
			idParsed, err := strconv.Atoi(idStr)
			idInt := uint(idParsed)
			if err != nil || idInt <= 0 {
				fmt.Println("ID invalide")
				os.Exit(1)
			}
			id = idInt
		}

		// Si les flags sont fournis, utiliser les valeurs des flags
		if updateName != "" || updateEmail != "" {
			name = strings.TrimSpace(updateName)
			email = strings.TrimSpace(updateEmail)
		} else {
			// Mode interactif
			name = strings.TrimSpace(readLine("Nouveau nom (laisser vide pour garder): "))
			email = strings.TrimSpace(readLine("Nouvel email (laisser vide pour garder): "))
		}

		if errUpdate := contact.UpdateContact(store, id, name, email); errUpdate != nil {
			fmt.Printf("Erreur lors de la mise à jour du contact: %s\n", errUpdate)
			os.Exit(1)
		}
		fmt.Println("Contact mis à jour.")
	},
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

func init() {
	rootCmd.AddCommand(updateCmd)

	// Définition des drapeaux pour la commande 'update'
	updateCmd.Flags().UintVarP(&updateID, "id", "i", 0, "ID du contact à mettre à jour")
	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Nouveau nom du contact")
	updateCmd.Flags().StringVarP(&updateEmail, "email", "e", "", "Nouvel email du contact")
}


