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

var contacts = ListeContacts{
	1: {Id: 1, Name: "John Doe", Email: "john.doe@example.com"},
	2: {Id: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
}

var lastId = 2

//Creates a new contact
func NewContact(name string, email string) (*Contact, error) {
	if name == "" {
		return nil, errors.New("le nom du contact est obligatoire")
	}
	if !IsValidEmail(email) {
		return nil, errors.New("le mail est vide ou incorrecte")
	}

	lastId++
	return &Contact{Id: lastId, Name: name, Email: email}, nil
}

//Add the reciever contact to the list of contacts
func (c *Contact) AddContact() {
	(*GetContacts())[c.Id] = c
}

//returns a pointer to the contact.contacts map
func GetContacts() *ListeContacts {
	return &contacts
}

func GetContactById(id int) (*Contact, error) {
	if (id <= 0) {
		return nil, errors.New("int must be bigger than 0")
	}

	contact, exists := (*GetContacts())[id]

	if !exists {
		return nil, errors.New("contact doesn't exist")
	}

	return contact, nil
}

//Delete reciever contact from the list of contacts
func (c Contact) DeleteContact() error {
	id := c.Id
	if _, exists := contacts[id]; exists {
		delete(*GetContacts(), id)
		return nil
	}
	return errors.New("contact doesn't exist")
}

func IsValidEmail(email string) bool {
	return strings.Contains(email, "@")
}

//Updates the contact and adds it to the list of contacts
func (c *Contact) UpdateContact(name string, email string) error {
	if name == "" {
		return errors.New("invalid name")
	}
	if !IsValidEmail(email) {
		return errors.New("invalid email")
	}

	id := (*c).Id
	(*c).Name = name
	(*c).Email = email

	//ici on ajoute le nouveau contact dans le dictionnaire mais je pense que ce n'est pas forcément nécessaire sous la forme actuelle vu que c'est déjà un tableau de pointeurs
	(*GetContacts())[id] = c

	return nil
}
