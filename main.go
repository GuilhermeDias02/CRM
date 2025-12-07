package main

import (
	"fmt"
	"os"

	"github.com/GuilhermeDias02/CRM/cmd"
	"github.com/GuilhermeDias02/CRM/internal/config"
	"github.com/GuilhermeDias02/CRM/internal/contact"
)

func main() {
	// Charger la configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors du chargement de la configuration: %v\n", err)
		os.Exit(1)
	}

	// Sélectionner le storer selon la configuration
	var store contact.Storer
	switch cfg.Storer.Type {
	case "memory":
		store = contact.NewMemoryStore()
		fmt.Println("Utilisation du storer en mémoire (MemoryStore)")

	case "json":
		jsonStore, err := contact.NewJsonStore(cfg.Storer.JSON.FilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de l'initialisation du JsonStore: %v\n", err)
			os.Exit(1)
		}
		store = jsonStore
		fmt.Printf("Utilisation du storer JSON (JsonStore) avec le fichier: %s\n", cfg.Storer.JSON.FilePath)

	case "gorm":
		gormStore, err := contact.NewGormStore(cfg.Storer.GORM.DBPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de l'initialisation du GormStore: %v\n", err)
			os.Exit(1)
		}
		store = gormStore
		fmt.Printf("Utilisation du storer GORM (GormStore) avec la base de données: %s\n", cfg.Storer.GORM.DBPath)

	default:
		fmt.Fprintf(os.Stderr, "Type de storer non supporté: %s\n", cfg.Storer.Type)
		os.Exit(1)
	}

	cmd.Execute(store)
}