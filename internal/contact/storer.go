package contact

type Storer interface {
	GetAll() map[uint]*Contact

	GetByID(id uint) (*Contact, error)

	Save(*Contact) (*Contact, error)

	Update(id uint, name *string, email *string) error

	Delete(id uint) error
}
