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
	deleteID int
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Supprimer un contact",
	Long: `La commande 'delete' permet de supprimer un contact existant par son ID.
L'ID peut être fourni via le flag --id ou demandé de manière interactive.`,
	Run: func(cmd *cobra.Command, args []string) {
		var id int

		if deleteID > 0 {
			id = deleteID
		} else {
			// Mode interactif
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("\nSupprimer un contact: ")

			fmt.Print("\tQuel est son id: ")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			idInt, err := strconv.Atoi(idStr)

			if err != nil {
				fmt.Println("\nL'id selectionné n'est pas un numéro")
				os.Exit(1)
			}
			id = idInt
		}

		if err := contact.DeleteContact(store, id); err != nil {
			fmt.Printf("\nCe contact n'existe pas: %s\n", err)
			os.Exit(1)
		}
		fmt.Println("Contact supprimé")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Définition des drapeaux pour la commande 'delete'
	deleteCmd.Flags().IntVarP(&deleteID, "id", "i", 0, "ID du contact à supprimer")
}


