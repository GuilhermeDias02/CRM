package contact

type Contact struct {
	Name  string
	Email string
}

var contacts = map[int]Contact{
	1: {Name: "John Doe", Email: "john.doe@example.com"},
	2: {Name: "Jane Doe", Email: "jane.doe@example.com"},
}

func AddContact(name, email string) {
	lastId := len(contacts)
	contacts[lastId+1] = Contact{
		Name:  name,
		Email: email,
	}
}

func GetContacts() map[int]Contact {
	return contacts
}
