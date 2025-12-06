package contact

import (
	"errors"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(dbPath string) (*GormStore, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the Contact schema
	if err := db.AutoMigrate(&Contact{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &GormStore{db: db}, nil
}

func (g *GormStore) GetAll() map[uint]*Contact {
	var contacts []Contact
	g.db.Find(&contacts)

	result := make(map[uint]*Contact)
	for i := range contacts {
		contact := &contacts[i]
		result[uint(contact.Id)] = contact
	}
	return result
}

func (g *GormStore) GetByID(id uint) (*Contact, error) {
	if id <= 0 {
		return nil, errors.New("ID must be greater than 0")
	}

	var contact Contact
	result := g.db.First(&contact, uint(id))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("contact doesn't exist")
		}
		return nil, fmt.Errorf("failed to get contact: %w", result.Error)
	}

	return &contact, nil
}

func (g *GormStore) Save(c *Contact) (*Contact, error) {
	if c.Name == "" {
		return nil, errors.New("le nom du contact est obligatoire")
	}
	if !IsValidEmail(c.Email) {
		return nil, errors.New("le mail est vide ou incorrecte")
	}

	if err := g.db.Create(c).Error; err != nil {
		return nil, fmt.Errorf("failed to save contact: %w", err)
	}

	return c, nil
}

func (g *GormStore) Update(id uint, name *string, email *string) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}

	var contact Contact
	result := g.db.First(&contact, uint(id))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("contact doesn't exist")
		}
		return fmt.Errorf("failed to get contact: %w", result.Error)
	}

	updates := make(map[string]interface{})
	if name != nil {
		if *name == "" {
			return errors.New("invalid name")
		}
		updates["name"] = *name
	}
	if email != nil {
		if !IsValidEmail(*email) {
			return errors.New("invalid email")
		}
		updates["email"] = *email
	}

	if len(updates) > 0 {
		if err := g.db.Model(&contact).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update contact: %w", err)
		}
	}

	return nil
}

func (g *GormStore) Delete(id uint) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}

	result := g.db.Delete(&Contact{}, uint(id))
	if result.Error != nil {
		return fmt.Errorf("failed to delete contact: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("contact doesn't exist")
	}

	return nil
}
