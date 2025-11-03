package contact

import "fmt"

type Contact struct {
	Name  string
	Email string
}

var contacts = map[int]Contact{
	1: {Name: "John Doe", Email: "john.doe@example.com"},
	2: {Name: "Jane Doe", Email: "jane.doe@example.com"},
}

func AddContact(name string, email string) {
	lastId := len(contacts) //modifiable avec une vrai bdd pour mieux gérer les suppressions
	contacts[lastId+1] = Contact{
		Name:  name,
		Email: email,
	}
}

func GetContacts() map[int]Contact {
	return contacts
}

func DeleteContact(id int) error {
	if _, exists := contacts[id]; exists {
		delete(contacts, id)
		return nil
	}
	return fmt.Errorf("Contact doesn't exist")
}
