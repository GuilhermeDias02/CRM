package contact

type Contact struct {
	ID    int
	Name  string
	Email string
}

var contacts = []Contact{
	{ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
	{ID: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
}

func AddContact(name, email string) {
	lastId := contacts[len(contacts)-1].ID
	contacts = append(contacts, Contact{ID: lastId + 1, Name: name, Email: email})
}

func GetContacts() []Contact {
	return contacts
}