package contact

import "testing"

func resetState() {
    contacts = ListeContacts{}
    lastId = 0
}

func TestNewContact_SuccessAndAdd(t *testing.T) {
    resetState()
    c, err := NewContact("Alice", "alice@example.com")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if c.Id != 1 {
        t.Fatalf("expected id 1, got %d", c.Id)
    }
    c.AddContact()
    if len(contacts) != 1 {
        t.Fatalf("expected 1 contact, got %d", len(contacts))
    }
}

func TestNewContact_Invalid(t *testing.T) {
    resetState()
    if _, err := NewContact("", "a@b.com"); err == nil {
        t.Fatalf("expected error for empty name")
    }
    if _, err := NewContact("Bob", "invalid"); err == nil {
        t.Fatalf("expected error for invalid email")
    }
}

func TestGetContactById(t *testing.T) {
    resetState()
    c, _ := NewContact("Bob", "bob@example.com")
    c.AddContact()
    got, err := GetContactById(1)
    if err != nil || got.Id != c.Id {
        t.Fatalf("expected to get contact id 1, err=%v", err)
    }
    if _, err := GetContactById(0); err == nil {
        t.Fatalf("expected error for id 0")
    }
    if _, err := GetContactById(999); err == nil {
        t.Fatalf("expected error for unknown id")
    }
}

func TestDeleteContact(t *testing.T) {
    resetState()
    c, _ := NewContact("Carol", "carol@example.com")
    c.AddContact()
    if err := c.DeleteContact(); err != nil {
        t.Fatalf("unexpected error deleting: %v", err)
    }
    if len(contacts) != 0 {
        t.Fatalf("expected 0 contacts after delete, got %d", len(contacts))
    }
    if err := c.DeleteContact(); err == nil {
        t.Fatalf("expected error deleting non-existent contact")
    }
}

func TestUpdateContact(t *testing.T) {
    resetState()
    c, _ := NewContact("Dave", "dave@example.com")
    c.AddContact()
    if err := c.UpdateContact("", "x@y.com"); err == nil {
        t.Fatalf("expected error for empty name")
    }
    if err := c.UpdateContact("New Dave", "bad"); err == nil {
        t.Fatalf("expected error for invalid email")
    }
    if err := c.UpdateContact("New Dave", "newdave@example.com"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if c.Name != "New Dave" || c.Email != "newdave@example.com" {
        t.Fatalf("update did not apply: %+v", c)
    }
}

func TestIsValidEmail(t *testing.T) {
    if !IsValidEmail("a@b") { t.Fatalf("expected valid") }
    if IsValidEmail("abc") { t.Fatalf("expected invalid") }
}


