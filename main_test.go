package main

import (
    "bytes"
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

func TestHandleAction_Quit(t *testing.T) {
	origAdd := addContactForm
	origDel := deleteContactForm
	addCalled := 0
	delCalled := 0
	addContactForm = func() { addCalled++ }
	deleteContactForm = func() { delCalled++ }
	defer func() { addContactForm = origAdd; deleteContactForm = origDel }()

	cont := handleAction(6)
	if cont {
		t.Fatalf("expected quit (false), got %v", cont)
	}
	if addCalled != 0 || delCalled != 0 {
		t.Fatalf("no action should be called, got add=%d del=%d", addCalled, delCalled)
	}
}

func TestHandleFlags_AddSuccess(t *testing.T) {
    contacts = map[int]Contact{}
    nextID = 1

    var out, errB bytes.Buffer
    handled, code := handleFlags([]string{"-add", "-name", "Alice", "-email", "alice@example.com"}, &out, &errB)
    if !handled || code != 0 {
        t.Fatalf("expected handled=true code=0, got handled=%v code=%d, err=%q", handled, code, errB.String())
    }
    if len(contacts) != 1 {
        t.Fatalf("expected 1 contact, got %d", len(contacts))
    }
    if !bytes.Contains(out.Bytes(), []byte("Contact ajouté avec ID")) {
        t.Fatalf("expected success message, got %q", out.String())
    }
}

func TestHandleFlags_MissingArgs(t *testing.T) {
    contacts = map[int]Contact{}
    nextID = 1

    var out, errB bytes.Buffer
    handled, code := handleFlags([]string{"-add", "-name", "Alice"}, &out, &errB)
    if !handled || code == 0 {
        t.Fatalf("expected handled=true error code!=0, got handled=%v code=%d", handled, code)
    }
}

func TestHandleFlags_InvalidEmail(t *testing.T) {
    contacts = map[int]Contact{}
    nextID = 1

    var out, errB bytes.Buffer
    handled, code := handleFlags([]string{"-add", "-name", "Bob", "-email", "invalid"}, &out, &errB)
    if !handled || code == 0 {
        t.Fatalf("expected handled=true error code!=0, got handled=%v code=%d", handled, code)
    }
}


