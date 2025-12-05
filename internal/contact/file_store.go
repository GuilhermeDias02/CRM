package contact

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const defaultDataFile = "contacts.json"

type FileStore struct {
	filePath string
	contacts map[int]*Contact
	nextID   int
	mu       sync.RWMutex
}

type fileData struct {
	Contacts map[int]*Contact `json:"contacts"`
	NextID   int              `json:"nextID"`
}

func NewFileStore(filePath string) (*FileStore, error) {
	if filePath == "" {
		// Utiliser le répertoire de travail actuel
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("impossible de déterminer le répertoire de travail: %w", err)
		}
		filePath = filepath.Join(wd, defaultDataFile)
	}

	store := &FileStore{
		filePath: filePath,
		contacts: make(map[int]*Contact),
		nextID:   1,
	}

	// Charger les données existantes
	if err := store.load(); err != nil {
		return nil, fmt.Errorf("erreur lors du chargement des données: %w", err)
	}

	return store, nil
}

func (f *FileStore) load() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Vérifier si le fichier existe
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		// Le fichier n'existe pas encore, on commence avec un store vide
		return nil
	}

	// Lire le fichier
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture du fichier: %w", err)
	}

	// Si le fichier est vide, on commence avec un store vide
	if len(data) == 0 {
		return nil
	}

	// Décoder le JSON
	var fileData fileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return fmt.Errorf("erreur lors du décodage JSON: %w", err)
	}

	// Restaurer les données
	if fileData.Contacts != nil {
		f.contacts = fileData.Contacts
	}
	if fileData.NextID > 0 {
		f.nextID = fileData.NextID
	}

	return nil
}

func (f *FileStore) save() error {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Préparer les données à sauvegarder
	data := fileData{
		Contacts: f.contacts,
		NextID:   f.nextID,
	}

	// Encoder en JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage JSON: %w", err)
	}

	// Écrire dans le fichier
	if err := os.WriteFile(f.filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier: %w", err)
	}

	return nil
}

func (f *FileStore) GetAll() map[int]*Contact {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make(map[int]*Contact)
	for k, v := range f.contacts {
		result[k] = v
	}
	return result
}

func (f *FileStore) GetByID(id int) (*Contact, error) {
	if id <= 0 {
		return nil, errors.New("ID must be greater than 0")
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	c, exists := f.contacts[id]
	if !exists {
		return nil, errors.New("contact doesn't exist")
	}
	return c, nil
}

func (f *FileStore) Save(c *Contact) (*Contact, error) {
	if c.Name == "" {
		return nil, errors.New("le nom du contact est obligatoire")
	}
	if !IsValidEmail(c.Email) {
		return nil, errors.New("le mail est vide ou incorrecte")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if c.Id == 0 {
		c.Id = f.nextID
		f.nextID++
	}

	f.contacts[c.Id] = c

	// Sauvegarder après modification
	if err := f.saveUnlocked(); err != nil {
		return nil, fmt.Errorf("erreur lors de la sauvegarde: %w", err)
	}

	return c, nil
}

func (f *FileStore) Update(id int, name *string, email *string) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	c, exists := f.contacts[id]
	if !exists {
		return errors.New("contact doesn't exist")
	}

	if name != nil {
		if *name == "" {
			return errors.New("invalid name")
		}
		c.Name = *name
	}
	if email != nil {
		if !IsValidEmail(*email) {
			return errors.New("invalid email")
		}
		c.Email = *email
	}

	// Sauvegarder après modification
	if err := f.saveUnlocked(); err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde: %w", err)
	}

	return nil
}

func (f *FileStore) Delete(id int) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.contacts[id]; !exists {
		return errors.New("contact doesn't exist")
	}

	delete(f.contacts, id)

	// Sauvegarder après modification
	if err := f.saveUnlocked(); err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde: %w", err)
	}

	return nil
}

// saveUnlocked est une version de save qui ne verrouille pas le mutex
// (à utiliser uniquement quand le mutex est déjà verrouillé)
func (f *FileStore) saveUnlocked() error {
	// Préparer les données à sauvegarder
	data := fileData{
		Contacts: f.contacts,
		NextID:   f.nextID,
	}

	// Encoder en JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors de l'encodage JSON: %w", err)
	}

	// Écrire dans le fichier
	if err := os.WriteFile(f.filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier: %w", err)
	}

	return nil
}
