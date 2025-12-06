package contact

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type toSave struct {
	Contacts map[int]*Contact `json:"contacts"`
	NextID   int              `json:"nextId"`
}

type JsonStore struct {
	contacts map[int]*Contact
	nextID   int
	filePath string
	mu       sync.RWMutex
}

func NewJsonStore(filePath string) (*JsonStore, error) {
	store := &JsonStore{
		contacts: make(map[int]*Contact),
		nextID:   1,
		filePath: filePath,
	}

	// Try to load existing data from file
	if err := store.loadFromFile(); err != nil {
		// If file doesn't exist and that we cannot create it, we return an error
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load contacts from file: %w", err)
		}
	}

	return store, nil
}

func (j *JsonStore) loadFromFile() error {
	data, err := os.ReadFile(j.filePath)
	if err != nil {
		return err
	}

	loadedSave := toSave {
		Contacts: make(map[int]*Contact),
		NextID: 0,
	}

	if err := json.Unmarshal(data, &loadedSave); err != nil {
		return fmt.Errorf("erreur lors de la désérialisation %s: %w", j.filePath, err)
	}

	j.mu.Lock()
	defer j.mu.Unlock()

	j.contacts = loadedSave.Contacts
	if j.contacts == nil {
		j.contacts = make(map[int]*Contact)
	}
	j.nextID = loadedSave.NextID
	if j.nextID < 1 {
		j.nextID = 1
	}

	return nil
}

func (j *JsonStore) saveToFile() error {
	j.mu.RLock()
	defer j.mu.RUnlock()

	data := toSave{
		Contacts: j.contacts,
		NextID:   j.nextID,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser les données: %w", err)
	}

	if err := os.WriteFile(j.filePath, jsonData, os.ModePerm); err != nil {
		return fmt.Errorf("impossible d'écrire le fichier %s: %w", j.filePath, err)
	}

	return nil
}

func (j *JsonStore) GetAll() map[int]*Contact {
	j.mu.RLock()
	defer j.mu.RUnlock()

	result := make(map[int]*Contact)
	for k, v := range j.contacts {
		result[k] = v
	}
	return result
}

func (j *JsonStore) GetByID(id int) (*Contact, error) {
	if id <= 0 {
		return nil, errors.New("ID must be greater than 0")
	}

	j.mu.RLock()
	defer j.mu.RUnlock()

	c, exists := j.contacts[id]
	if !exists {
		return nil, errors.New("contact doesn't exist")
	}
	return c, nil
}

func (j *JsonStore) Save(c *Contact) (*Contact, error) {
	if c.Name == "" {
		return nil, errors.New("le nom du contact est obligatoire")
	}
	if !IsValidEmail(c.Email) {
		return nil, errors.New("le mail est vide ou incorrecte")
	}

	j.mu.Lock()
	if c.Id == 0 {
		c.Id = j.nextID
		j.nextID++
	}
	j.contacts[c.Id] = c
	j.mu.Unlock()

	// Save to file after modification
	if err := j.saveToFile(); err != nil {
		return c, fmt.Errorf("failed to save to file: %w", err)
	}

	return c, nil
}

func (j *JsonStore) Update(id int, name *string, email *string) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}

	j.mu.Lock()
	c, exists := j.contacts[id]
	if !exists {
		j.mu.Unlock()
		return errors.New("contact doesn't exist")
	}

	if name != nil {
		if *name == "" {
			j.mu.Unlock()
			return errors.New("invalid name")
		}
		c.Name = *name
	}
	if email != nil {
		if !IsValidEmail(*email) {
			j.mu.Unlock()
			return errors.New("invalid email")
		}
		c.Email = *email
	}
	j.mu.Unlock()

	// Save to file after modification
	if err := j.saveToFile(); err != nil {
		return fmt.Errorf("failed to save to file: %w", err)
	}

	return nil
}

func (j *JsonStore) Delete(id int) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}

	j.mu.Lock()
	if _, exists := j.contacts[id]; !exists {
		j.mu.Unlock()
		return errors.New("contact doesn't exist")
	}

	delete(j.contacts, id)
	j.mu.Unlock()

	// Save to file after modification
	if err := j.saveToFile(); err != nil {
		return fmt.Errorf("failed to save to file: %w", err)
	}

	return nil
}
