package contact

import (
	"errors"
	"strings"
)

type Contact struct {
	Id    int
	Name  string
	Email string
}

type ListeContacts = map[int]*Contact

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

func GetContacts(store Storer) map[int]*Contact {
	return store.GetAll()
}

func GetContactById(store Storer, id int) (*Contact, error) {
	return store.GetByID(id)
}

func DeleteContact(store Storer, id int) error {
	return store.Delete(id)
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@")
}

func UpdateContact(store Storer, id int, name string, email string) error {
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
