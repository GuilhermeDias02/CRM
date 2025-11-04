package contact

import (
	"testing"
)

// TestMemoryStore_NewMemoryStore vérifie que NewMemoryStore crée un store vide.
func TestMemoryStore_NewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	if store == nil {
		t.Fatalf("NewMemoryStore should not return nil")
	}
	if len(store.GetAll()) != 0 {
		t.Fatalf("expected empty store, got %d contacts", len(store.GetAll()))
	}
}

// TestMemoryStore_GetAll vérifie que GetAll retourne tous les contacts.
func TestMemoryStore_GetAll(t *testing.T) {
	store := NewMemoryStore()
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

// TestMemoryStore_GetAll_ReturnsCopy vérifie que GetAll retourne une copie.
func TestMemoryStore_GetAll_ReturnsCopy(t *testing.T) {
	store := NewMemoryStore()
	store.Save(&Contact{Name: "Charlie", Email: "charlie@example.com"})

	all := store.GetAll()
	all[999] = &Contact{Id: 999, Name: "Hacked", Email: "hacked@example.com"}

	if len(store.GetAll()) != 1 {
		t.Fatalf("modifying returned map should not affect store, got %d contacts", len(store.GetAll()))
	}
}

// TestMemoryStore_GetByID_Success vérifie que GetByID récupère un contact existant.
func TestMemoryStore_GetByID_Success(t *testing.T) {
	store := NewMemoryStore()
	saved, _ := store.Save(&Contact{Name: "David", Email: "david@example.com"})

	got, err := store.GetByID(saved.Id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "David" {
		t.Errorf("expected David, got %s", got.Name)
	}
	if got.Email != "david@example.com" {
		t.Errorf("expected david@example.com, got %s", got.Email)
	}
}

// TestMemoryStore_GetByID_NotFound vérifie que GetByID retourne une erreur pour un ID inexistant.
func TestMemoryStore_GetByID_NotFound(t *testing.T) {
	store := NewMemoryStore()
	_, err := store.GetByID(999)
	if err == nil {
		t.Fatalf("expected error for non-existent ID")
	}
}

// TestMemoryStore_GetByID_InvalidID vérifie que GetByID retourne une erreur pour un ID invalide.
func TestMemoryStore_GetByID_InvalidID(t *testing.T) {
	store := NewMemoryStore()
	_, err := store.GetByID(0)
	if err == nil {
		t.Fatalf("expected error for ID 0")
	}
	_, err = store.GetByID(-1)
	if err == nil {
		t.Fatalf("expected error for negative ID")
	}
}

// TestMemoryStore_Save_GeneratesID vérifie que Save génère un ID si Contact.Id == 0.
func TestMemoryStore_Save_GeneratesID(t *testing.T) {
	store := NewMemoryStore()
	c := &Contact{Name: "Eve", Email: "eve@example.com"}
	if c.Id != 0 {
		t.Fatalf("contact should start with Id 0")
	}

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
	if saved.Id != 1 {
		t.Errorf("expected first ID to be 1, got %d", saved.Id)
	}
}

// TestMemoryStore_Save_IncrementsID vérifie que Save incrémente correctement les IDs.
func TestMemoryStore_Save_IncrementsID(t *testing.T) {
	store := NewMemoryStore()
	c1, _ := store.Save(&Contact{Name: "Frank", Email: "frank@example.com"})
	c2, _ := store.Save(&Contact{Name: "Grace", Email: "grace@example.com"})
	c3, _ := store.Save(&Contact{Name: "Henry", Email: "henry@example.com"})

	if c1.Id != 1 || c2.Id != 2 || c3.Id != 3 {
		t.Errorf("expected IDs 1, 2, 3, got %d, %d, %d", c1.Id, c2.Id, c3.Id)
	}
}

// TestMemoryStore_Save_WithExistingID vérifie que Save accepte un contact avec ID existant.
func TestMemoryStore_Save_WithExistingID(t *testing.T) {
	store := NewMemoryStore()
	c := &Contact{Id: 42, Name: "Iris", Email: "iris@example.com"}
	saved, err := store.Save(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if saved.Id != 42 {
		t.Errorf("expected ID 42, got %d", saved.Id)
	}

	got, _ := store.GetByID(42)
	if got.Name != "Iris" {
		t.Errorf("expected Iris, got %s", got.Name)
	}
}

// TestMemoryStore_Save_ValidationErrors vérifie que Save valide les données.
func TestMemoryStore_Save_ValidationErrors(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.Save(&Contact{Name: "", Email: "test@example.com"})
	if err == nil {
		t.Fatalf("expected error for empty name")
	}

	_, err = store.Save(&Contact{Name: "Test", Email: "invalid"})
	if err == nil {
		t.Fatalf("expected error for invalid email")
	}

	_, err = store.Save(&Contact{Name: "Test", Email: ""})
	if err == nil {
		t.Fatalf("expected error for empty email")
	}
}

// TestMemoryStore_Update_Success vérifie que Update modifie un contact existant.
func TestMemoryStore_Update_Success(t *testing.T) {
	store := NewMemoryStore()
	c, _ := store.Save(&Contact{Name: "Jack", Email: "jack@example.com"})

	newName := "Jack Updated"
	err := store.Update(c.Id, &newName, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := store.GetByID(c.Id)
	if updated.Name != "Jack Updated" {
		t.Errorf("expected 'Jack Updated', got %s", updated.Name)
	}
	if updated.Email != "jack@example.com" {
		t.Errorf("email should remain unchanged, got %s", updated.Email)
	}
}

// TestMemoryStore_Update_PartialFields vérifie que Update permet la mise à jour partielle.
func TestMemoryStore_Update_PartialFields(t *testing.T) {
	store := NewMemoryStore()
	c, _ := store.Save(&Contact{Name: "Kate", Email: "kate@example.com"})

	err := store.Update(c.Id, nil, nil)
	if err != nil {
		t.Fatalf("update with nil fields should not error, got: %v", err)
	}

	updated, _ := store.GetByID(c.Id)
	if updated.Name != "Kate" || updated.Email != "kate@example.com" {
		t.Errorf("contact should remain unchanged")
	}

	newEmail := "kate.new@example.com"
	err = store.Update(c.Id, nil, &newEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ = store.GetByID(c.Id)
	if updated.Email != "kate.new@example.com" {
		t.Errorf("expected 'kate.new@example.com', got %s", updated.Email)
	}
	if updated.Name != "Kate" {
		t.Errorf("name should remain unchanged, got %s", updated.Name)
	}
}

// TestMemoryStore_Update_Errors vérifie que Update retourne des erreurs appropriées.
func TestMemoryStore_Update_Errors(t *testing.T) {
	store := NewMemoryStore()

	err := store.Update(999, nil, nil)
	if err == nil {
		t.Fatalf("expected error for non-existent ID")
	}

	err = store.Update(0, nil, nil)
	if err == nil {
		t.Fatalf("expected error for ID 0")
	}

	c, _ := store.Save(&Contact{Name: "Liam", Email: "liam@example.com"})
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

// TestMemoryStore_Delete_Success vérifie que Delete supprime un contact.
func TestMemoryStore_Delete_Success(t *testing.T) {
	store := NewMemoryStore()
	c, _ := store.Save(&Contact{Name: "Mia", Email: "mia@example.com"})

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

// TestMemoryStore_Delete_Errors vérifie que Delete retourne des erreurs appropriées.
func TestMemoryStore_Delete_Errors(t *testing.T) {
	store := NewMemoryStore()

	err := store.Delete(999)
	if err == nil {
		t.Fatalf("expected error for non-existent ID")
	}

	err = store.Delete(0)
	if err == nil {
		t.Fatalf("expected error for ID 0")
	}

	err = store.Delete(-1)
	if err == nil {
		t.Fatalf("expected error for negative ID")
	}
}

// TestMemoryStore_Delete_DoesNotReuseID vérifie que les IDs ne sont pas réutilisés après suppression.
func TestMemoryStore_Delete_DoesNotReuseID(t *testing.T) {
	store := NewMemoryStore()
	c1, _ := store.Save(&Contact{Name: "Noah", Email: "noah@example.com"})
	c2, _ := store.Save(&Contact{Name: "Olivia", Email: "olivia@example.com"})

	if c1.Id != 1 || c2.Id != 2 {
		t.Fatalf("setup failed: expected IDs 1, 2")
	}

	store.Delete(c1.Id)

	c3, _ := store.Save(&Contact{Name: "Paul", Email: "paul@example.com"})
	if c3.Id != 3 {
		t.Errorf("expected ID 3 (not reused), got %d", c3.Id)
	}
}

// TestMemoryStore_ImplementsStorer vérifie que MemoryStore implémente bien l'interface Storer.
func TestMemoryStore_ImplementsStorer(t *testing.T) {
	var store Storer = NewMemoryStore()
	if store == nil {
		t.Fatalf("MemoryStore should implement Storer interface")
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

// TestMemoryStore_Integration vérifie un scénario d'utilisation complet.
func TestMemoryStore_Integration(t *testing.T) {
	store := NewMemoryStore()

	c1, _ := store.Save(&Contact{Name: "Quinn", Email: "quinn@example.com"})
	c2, _ := store.Save(&Contact{Name: "Rachel", Email: "rachel@example.com"})

	if len(store.GetAll()) != 2 {
		t.Fatalf("expected 2 contacts")
	}

	got, _ := store.GetByID(c1.Id)
	if got.Name != "Quinn" {
		t.Errorf("expected Quinn, got %s", got.Name)
	}

	newName := "Quinn Updated"
	store.Update(c1.Id, &newName, nil)

	got, _ = store.GetByID(c1.Id)
	if got.Name != "Quinn Updated" {
		t.Errorf("expected 'Quinn Updated', got %s", got.Name)
	}

	store.Delete(c2.Id)

	if len(store.GetAll()) != 1 {
		t.Fatalf("expected 1 contact after delete, got %d", len(store.GetAll()))
	}

	c3, _ := store.Save(&Contact{Name: "Sam", Email: "sam@example.com"})
	if c3.Id != 3 {
		t.Errorf("expected ID 3, got %d", c3.Id)
	}
}

