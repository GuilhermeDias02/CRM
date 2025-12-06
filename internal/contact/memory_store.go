package contact

import (
	"errors"
	"sync"
)

type MemoryStore struct {
	contacts map[uint]*Contact
	nextID   uint
	mu       sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		contacts: make(map[uint]*Contact),
		nextID:   1,
	}
}

func (m *MemoryStore) GetAll() map[uint]*Contact {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[uint]*Contact)
	for k, v := range m.contacts {
		result[k] = v
	}
	return result
}

func (m *MemoryStore) GetByID(id uint) (*Contact, error) {
	if id <= 0 {
		return nil, errors.New("ID must be greater than 0")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	c, exists := m.contacts[id]
	if !exists {
		return nil, errors.New("contact doesn't exist")
	}
	return c, nil
}

func (m *MemoryStore) Save(c *Contact) (*Contact, error) {
	if c.Name == "" {
		return nil, errors.New("le nom du contact est obligatoire")
	}
	if !IsValidEmail(c.Email) {
		return nil, errors.New("le mail est vide ou incorrecte")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if c.Id == 0 {
		c.Id = m.nextID
		m.nextID++
	}

	m.contacts[c.Id] = c
	return c, nil
}

func (m *MemoryStore) Update(id uint, name *string, email *string) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	c, exists := m.contacts[id]
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

	return nil
}

func (m *MemoryStore) Delete(id uint) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.contacts[id]; !exists {
		return errors.New("contact doesn't exist")
	}

	delete(m.contacts, id)
	return nil
}

