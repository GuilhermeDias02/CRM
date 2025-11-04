package contact

import (
	"errors"
	"testing"
)

// mockStorer est une implémentation de test de l'interface Storer
// pour valider que l'interface est correctement définie et utilisable.
type mockStorer struct {
	contacts map[int]*Contact
	nextID   int
}

func newMockStorer() *mockStorer {
	return &mockStorer{
		contacts: make(map[int]*Contact),
		nextID:   1,
	}
}

func (m *mockStorer) GetAll() map[int]*Contact {
	result := make(map[int]*Contact)
	for k, v := range m.contacts {
		result[k] = v
	}
	return result
}

func (m *mockStorer) GetByID(id int) (*Contact, error) {
	if id <= 0 {
		return nil, errors.New("ID must be greater than 0")
	}
	c, exists := m.contacts[id]
	if !exists {
		return nil, errors.New("contact doesn't exist")
	}
	return c, nil
}

func (m *mockStorer) Save(c *Contact) (*Contact, error) {
	if c.Id == 0 {
		c.Id = m.nextID
		m.nextID++
	}
	if c.Name == "" {
		return nil, errors.New("name is required")
	}
	if !IsValidEmail(c.Email) {
		return nil, errors.New("invalid email")
	}
	m.contacts[c.Id] = c
	return c, nil
}

func (m *mockStorer) Update(id int, name *string, email *string) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}
	c, exists := m.contacts[id]
	if !exists {
		return errors.New("contact doesn't exist")
	}
	if name != nil {
		if *name == "" {
			return errors.New("name cannot be empty")
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

func (m *mockStorer) Delete(id int) error {
	if id <= 0 {
		return errors.New("ID must be greater than 0")
	}
	if _, exists := m.contacts[id]; !exists {
		return errors.New("contact doesn't exist")
	}
	delete(m.contacts, id)
	return nil
}

// TestStorerInterface_GetAll vérifie que GetAll retourne tous les contacts.
func TestStorerInterface_GetAll(t *testing.T) {
	store := newMockStorer()
	c1, _ := store.Save(&Contact{Name: "Alice", Email: "alice@example.com"})
	c2, _ := store.Save(&Contact{Name: "Bob", Email: "bob@example.com"})

	all := store.GetAll()
	if len(all) != 2 {
		t.Fatalf("expected 2 contacts, got %d", len(all))
	}
	if all[c1.Id].Name != "Alice" {
		t.Errorf("expected Alice, got %s", all[c1.Id].Name)
	}
	if all[c2.Id].Name != "Bob" {
		t.Errorf("expected Bob, got %s", all[c2.Id].Name)
	}
}

// TestStorerInterface_GetByID vérifie que GetByID récupère un contact par ID.
func TestStorerInterface_GetByID(t *testing.T) {
	store := newMockStorer()
	c, _ := store.Save(&Contact{Name: "Charlie", Email: "charlie@example.com"})

	got, err := store.GetByID(c.Id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Charlie" {
		t.Errorf("expected Charlie, got %s", got.Name)
	}

	_, err = store.GetByID(999)
	if err == nil {
		t.Fatalf("expected error for non-existent ID")
	}

	_, err = store.GetByID(0)
	if err == nil {
		t.Fatalf("expected error for invalid ID")
	}
}

// TestStorerInterface_Save vérifie que Save persiste un contact et génère un ID si nécessaire.
func TestStorerInterface_Save(t *testing.T) {
	store := newMockStorer()

	c := &Contact{Name: "David", Email: "david@example.com"}
	saved, err := store.Save(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved.Id == 0 {
		t.Fatalf("expected ID to be generated, got 0")
	}
	if saved.Id != c.Id {
		t.Errorf("expected ID %d, got %d", saved.Id, c.Id)
	}

	all := store.GetAll()
	if len(all) != 1 {
		t.Fatalf("expected 1 contact, got %d", len(all))
	}

	_, err = store.Save(&Contact{Name: "", Email: "test@example.com"})
	if err == nil {
		t.Fatalf("expected error for empty name")
	}

	_, err = store.Save(&Contact{Name: "Test", Email: "invalid"})
	if err == nil {
		t.Fatalf("expected error for invalid email")
	}
}

// TestStorerInterface_SaveWithExistingID vérifie que Save avec un ID existant fonctionne.
func TestStorerInterface_SaveWithExistingID(t *testing.T) {
	store := newMockStorer()
	c := &Contact{Id: 42, Name: "Eve", Email: "eve@example.com"}
	saved, err := store.Save(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved.Id != 42 {
		t.Errorf("expected ID 42, got %d", saved.Id)
	}
}

// TestStorerInterface_Update vérifie que Update modifie un contact existant.
func TestStorerInterface_Update(t *testing.T) {
	store := newMockStorer()
	c, _ := store.Save(&Contact{Name: "Frank", Email: "frank@example.com"})

	newName := "Frank Updated"
	err := store.Update(c.Id, &newName, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := store.GetByID(c.Id)
	if updated.Name != "Frank Updated" {
		t.Errorf("expected 'Frank Updated', got %s", updated.Name)
	}
	if updated.Email != "frank@example.com" {
		t.Errorf("email should remain unchanged, got %s", updated.Email)
	}

	newEmail := "frank.new@example.com"
	err = store.Update(c.Id, nil, &newEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ = store.GetByID(c.Id)
	if updated.Email != "frank.new@example.com" {
		t.Errorf("expected 'frank.new@example.com', got %s", updated.Email)
	}
}

// TestStorerInterface_UpdatePartialFields vérifie que Update permet de modifier seulement certains champs.
func TestStorerInterface_UpdatePartialFields(t *testing.T) {
	store := newMockStorer()
	c, _ := store.Save(&Contact{Name: "Grace", Email: "grace@example.com"})

	err := store.Update(c.Id, nil, nil)
	if err != nil {
		t.Fatalf("update with nil fields should not error, got: %v", err)
	}

	updated, _ := store.GetByID(c.Id)
	if updated.Name != "Grace" || updated.Email != "grace@example.com" {
		t.Errorf("contact should remain unchanged")
	}
}

// TestStorerInterface_UpdateErrors vérifie que Update retourne des erreurs appropriées.
func TestStorerInterface_UpdateErrors(t *testing.T) {
	store := newMockStorer()

	err := store.Update(999, nil, nil)
	if err == nil {
		t.Fatalf("expected error for non-existent ID")
	}

	err = store.Update(0, nil, nil)
	if err == nil {
		t.Fatalf("expected error for invalid ID")
	}

	c, _ := store.Save(&Contact{Name: "Henry", Email: "henry@example.com"})
	emptyName := ""
	err = store.Update(c.Id, &emptyName, nil)
	if err == nil {
		t.Fatalf("expected error for empty name")
	}

	invalidEmail := "invalid"
	err = store.Update(c.Id, nil, &invalidEmail)
	if err == nil {
		t.Fatalf("expected error for invalid email")
	}
}

// TestStorerInterface_Delete vérifie que Delete supprime un contact.
func TestStorerInterface_Delete(t *testing.T) {
	store := newMockStorer()
	c, _ := store.Save(&Contact{Name: "Iris", Email: "iris@example.com"})

	err := store.Delete(c.Id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	all := store.GetAll()
	if len(all) != 0 {
		t.Fatalf("expected 0 contacts after delete, got %d", len(all))
	}

	_, err = store.GetByID(c.Id)
	if err == nil {
		t.Fatalf("expected error when getting deleted contact")
	}
}

// TestStorerInterface_DeleteErrors vérifie que Delete retourne des erreurs appropriées.
func TestStorerInterface_DeleteErrors(t *testing.T) {
	store := newMockStorer()

	err := store.Delete(999)
	if err == nil {
		t.Fatalf("expected error for non-existent ID")
	}

	err = store.Delete(0)
	if err == nil {
		t.Fatalf("expected error for invalid ID")
	}
}

// TestStorerInterface_Conformity vérifie qu'une implémentation mock satisfait bien l'interface Storer.
func TestStorerInterface_Conformity(t *testing.T) {
	var store Storer = newMockStorer()
	if store == nil {
		t.Fatalf("mockStorer should implement Storer interface")
	}

	_, err := store.GetByID(1)
	if err == nil {
		t.Fatalf("expected error for non-existent contact")
	}

	all := store.GetAll()
	if all == nil {
		t.Fatalf("GetAll should not return nil")
	}
}

// TestStorerInterface_Integration vérifie un scénario d'utilisation complet de l'interface.
func TestStorerInterface_Integration(t *testing.T) {
	store := newMockStorer()

	c1, _ := store.Save(&Contact{Name: "Alice", Email: "alice@example.com"})
	c2, _ := store.Save(&Contact{Name: "Bob", Email: "bob@example.com"})

	if len(store.GetAll()) != 2 {
		t.Fatalf("expected 2 contacts")
	}

	got, _ := store.GetByID(c1.Id)
	if got.Name != "Alice" {
		t.Errorf("expected Alice, got %s", got.Name)
	}

	newName := "Alice Updated"
	store.Update(c1.Id, &newName, nil)

	got, _ = store.GetByID(c1.Id)
	if got.Name != "Alice Updated" {
		t.Errorf("expected 'Alice Updated', got %s", got.Name)
	}

	store.Delete(c2.Id)

	if len(store.GetAll()) != 1 {
		t.Fatalf("expected 1 contact after delete, got %d", len(store.GetAll()))
	}
}

