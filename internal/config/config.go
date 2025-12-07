package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config représente la configuration de l'application
type Config struct {
	Storer StorerConfig `mapstructure:"storer"`
}

// StorerConfig représente la configuration du storer
type StorerConfig struct {
	Type string      `mapstructure:"type"`
	JSON JSONConfig  `mapstructure:"json"`
	GORM GORMConfig  `mapstructure:"gorm"`
}

// JSONConfig représente la configuration pour le storer JSON
type JSONConfig struct {
	FilePath string `mapstructure:"file_path"`
}

// GORMConfig représente la configuration pour le storer GORM
type GORMConfig struct {
	DBPath string `mapstructure:"db_path"`
}

// LoadConfig charge la configuration depuis le fichier config.yaml
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./CRM")

	// Valeurs par défaut
	viper.SetDefault("storer.type", "gorm")
	viper.SetDefault("storer.json.file_path", "contacts.json")
	viper.SetDefault("storer.gorm.db_path", "contacts.db")

	// Lire le fichier de configuration
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Le fichier de configuration n'existe pas, on utilise les valeurs par défaut
			log.Printf("Fichier config.yaml non trouvé, utilisation des valeurs par défaut")
		} else {
			return nil, fmt.Errorf("erreur lors de la lecture du fichier de configuration: %w", err)
		}
	} else {
		log.Printf("Configuration chargée depuis: %s", viper.ConfigFileUsed())
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing de la configuration: %w", err)
	}

	// Validation du type de storer
	if config.Storer.Type != "memory" && config.Storer.Type != "json" && config.Storer.Type != "gorm" {
		return nil, fmt.Errorf("type de storer invalide: %s. Options valides: memory, json, gorm", config.Storer.Type)
	}

	return &config, nil
}

