package contact

import (
	"errors"
	"strings"
)

type Contact struct {
	Id    uint `gorm:"primaryKey"`
	Name  string `gorm:"not null;size:100"`
	Email string `gorm:"not null;size:100"`
}

type ListeContacts = map[uint]*Contact

func NewContact(store Storer, name string, email string) (*Contact, error) {
	if name == "" {
		return nil, errors.New("le nom du contact est obligatoire")
	}
	if !IsValidEmail(email) {
		return nil, errors.New("le mail est vide ou incorrecte")
	}

	newContact := &Contact{Name: name, Email: email}
	saved, err := store.Save(newContact)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func GetContacts(store Storer) map[uint]*Contact {
	return store.GetAll()
}

func GetContactById(store Storer, id uint) (*Contact, error) {
	return store.GetByID(id)
}

func DeleteContact(store Storer, id uint) error {
	return store.Delete(id)
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@")
}

func UpdateContact(store Storer, id uint, name string, email string) error {
	var namePtr *string
	var emailPtr *string

	if name != "" {
		namePtr = &name
	}
	if email != "" {
		emailPtr = &email
	}

	return store.Update(id, namePtr, emailPtr)
}
