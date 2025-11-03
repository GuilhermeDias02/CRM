package main

import (
    "testing"
)

func TestUpdateContact_Success_NameOnly(t *testing.T) {
	contacts = map[int]Contact{}
	nextID = 1
	id := addContact("Alice", "alice@example.com")

	newName := "Alice Cooper"
	if err := updateContactByID(id, &newName, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := contacts[id]
	if got.Name != "Alice Cooper" {
		t.Errorf("expected name updated to 'Alice Cooper', got %q", got.Name)
	}
	if got.Email != "alice@example.com" {
		t.Errorf("email should remain unchanged, got %q", got.Email)
	}
}

func TestUpdateContact_Error_InvalidEmail(t *testing.T) {
	contacts = map[int]Contact{}
	nextID = 1
	id := addContact("Bob", "bob@example.com")

	badEmail := "invalid-email"
	err := updateContactByID(id, nil, &badEmail)
	if err == nil {
		t.Fatalf("expected error for invalid email, got nil")
	}
}

func TestUpdateContact_Error_UnknownID(t *testing.T) {
	contacts = map[int]Contact{}
	nextID = 1

	newName := "Charlie"
	err := updateContactByID(999, &newName, nil)
	if err == nil {
		t.Fatalf("expected error for unknown ID, got nil")
	}
}

func TestHandleAction_Update(t *testing.T) {
	called := 0
	orig := updateContactFunc
	updateContactFunc = func() { called++ }
	defer func() { updateContactFunc = orig }()

	if !handleAction(5) { t.Fatalf("expected continue for action 5") }
	if called != 1 { t.Fatalf("expected update called once, got %d", called) }
}


