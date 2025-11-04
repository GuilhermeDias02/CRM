package contact

import "testing"

func TestNewContact_Success(t *testing.T) {
	store := NewMemoryStore()
	c, err := NewContact(store, "Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Id != 1 {
		t.Fatalf("expected id 1, got %d", c.Id)
	}
	all := GetContacts(store)
	if len(all) != 1 {
		t.Fatalf("expected 1 contact, got %d", len(all))
	}
	if all[c.Id].Name != "Alice" || all[c.Id].Email != "alice@example.com" {
		t.Fatalf("contact not saved correctly: %+v", all[c.Id])
	}
}

func TestNewContact_Invalid(t *testing.T) {
	store := NewMemoryStore()
	if _, err := NewContact(store, "", "a@b.com"); err == nil {
		t.Fatalf("expected error for empty name")
	}
	if _, err := NewContact(store, "Bob", "invalid"); err == nil {
		t.Fatalf("expected error for invalid email")
	}
}

func TestGetContactById(t *testing.T) {
	store := NewMemoryStore()
	c, _ := NewContact(store, "Bob", "bob@example.com")
	got, err := GetContactById(store, c.Id)
	if err != nil || got.Id != c.Id {
		t.Fatalf("expected to get contact id %d, err=%v", c.Id, err)
	}
	if _, err := GetContactById(store, 0); err == nil {
		t.Fatalf("expected error for id 0")
	}
	if _, err := GetContactById(store, 999); err == nil {
		t.Fatalf("expected error for unknown id")
	}
}

func TestDeleteContact(t *testing.T) {
	store := NewMemoryStore()
	c, _ := NewContact(store, "Carol", "carol@example.com")
	if err := DeleteContact(store, c.Id); err != nil {
		t.Fatalf("unexpected error deleting: %v", err)
	}
	all := GetContacts(store)
	if len(all) != 0 {
		t.Fatalf("expected 0 contacts after delete, got %d", len(all))
	}
	if err := DeleteContact(store, c.Id); err == nil {
		t.Fatalf("expected error deleting non-existent contact")
	}
}

func TestUpdateContact(t *testing.T) {
	store := NewMemoryStore()
	c, _ := NewContact(store, "Dave", "dave@example.com")
	if err := UpdateContact(store, c.Id, "", "x@y.com"); err != nil {
		t.Fatalf("empty name should not trigger update, got error: %v", err)
	}
	updated, _ := GetContactById(store, c.Id)
	if updated.Name != "Dave" || updated.Email != "x@y.com" {
		t.Fatalf("name should remain unchanged, email should be updated: %+v", updated)
	}
	if err := UpdateContact(store, c.Id, "New Dave", "bad"); err == nil {
		t.Fatalf("expected error for invalid email")
	}
	if err := UpdateContact(store, c.Id, "New Dave", "newdave@example.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated, _ = GetContactById(store, c.Id)
	if updated.Name != "New Dave" || updated.Email != "newdave@example.com" {
		t.Fatalf("update did not apply: %+v", updated)
	}
}

func TestUpdateContact_Partial(t *testing.T) {
	store := NewMemoryStore()
	c, _ := NewContact(store, "Eve", "eve@example.com")
	if err := UpdateContact(store, c.Id, "New Eve", ""); err != nil {
		t.Fatalf("unexpected error updating name only: %v", err)
	}
	updated, _ := GetContactById(store, c.Id)
	if updated.Name != "New Eve" || updated.Email != "eve@example.com" {
		t.Fatalf("partial update failed: %+v", updated)
	}
	if err := UpdateContact(store, c.Id, "", "neweve@example.com"); err != nil {
		t.Fatalf("unexpected error updating email only: %v", err)
	}
	updated, _ = GetContactById(store, c.Id)
	if updated.Name != "New Eve" || updated.Email != "neweve@example.com" {
		t.Fatalf("partial update failed: %+v", updated)
	}
}

func TestIsValidEmail(t *testing.T) {
	if !IsValidEmail("a@b") {
		t.Fatalf("expected valid")
	}
	if IsValidEmail("abc") {
		t.Fatalf("expected invalid")
	}
}


