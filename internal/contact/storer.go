package contact

type Storer interface {
	GetAll() map[int]*Contact

	GetByID(id int) (*Contact, error)

	Save(*Contact) (*Contact, error)

	Update(id int, name *string, email *string) error

	Delete(id int) error
}

