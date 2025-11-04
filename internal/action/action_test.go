package action

import (
	"bytes"
	"testing"

	"github.com/GuilhermeDias02/CRM/internal/contact"
)

func TestHandleFlags_AddSuccess(t *testing.T) {
	store := contact.NewMemoryStore()
	var out, errB bytes.Buffer
	handled, code := HandleFlags(store, []string{"-add", "-name", "Alice", "-email", "alice@example.com"}, &out, &errB)
	if !handled || code != 0 {
		t.Fatalf("expected handled=true code=0, got handled=%v code=%d, err=%q", handled, code, errB.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("Contact ajouté avec ID")) {
		t.Fatalf("expected success message, got %q", out.String())
	}
	if len(store.GetAll()) != 1 {
		t.Fatalf("expected 1 contact in store, got %d", len(store.GetAll()))
	}
}

func TestHandleFlags_MissingArgs(t *testing.T) {
	store := contact.NewMemoryStore()
	var out, errB bytes.Buffer
	handled, code := HandleFlags(store, []string{"-add", "-name", "Alice"}, &out, &errB)
	if !handled || code == 0 {
		t.Fatalf("expected handled=true and non-zero code, got handled=%v code=%d", handled, code)
	}
}

func TestHandleFlags_InvalidEmail(t *testing.T) {
	store := contact.NewMemoryStore()
	var out, errB bytes.Buffer
	handled, code := HandleFlags(store, []string{"-add", "-name", "Bob", "-email", "invalid"}, &out, &errB)
	if !handled || code == 0 {
		t.Fatalf("expected handled=true and non-zero code, got handled=%v code=%d", handled, code)
	}
}

func TestHandleAction_DefaultAndQuit(t *testing.T) {
	store := contact.NewMemoryStore()
	if cont := HandleAction(store, 99); !cont {
		t.Fatalf("expected continue=true for unknown action")
	}
	if cont := HandleAction(store, 4); cont {
		t.Fatalf("expected continue=false for quit action")
	}
}


