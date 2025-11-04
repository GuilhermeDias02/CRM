package action

import (
	"bytes"
	"testing"
)

func TestHandleFlags_AddSuccess(t *testing.T) {
    // reset contact package state by adding via flags; just ensure handled and code
    var out, errB bytes.Buffer
    handled, code := HandleFlags([]string{"-add", "-name", "Alice", "-email", "alice@example.com"}, &out, &errB)
    if !handled || code != 0 {
        t.Fatalf("expected handled=true code=0, got handled=%v code=%d, err=%q", handled, code, errB.String())
    }
    if !bytes.Contains(out.Bytes(), []byte("Contact ajouté avec ID")) {
        t.Fatalf("expected success message, got %q", out.String())
    }
}

func TestHandleFlags_MissingArgs(t *testing.T) {
    var out, errB bytes.Buffer
    handled, code := HandleFlags([]string{"-add", "-name", "Alice"}, &out, &errB)
    if !handled || code == 0 {
        t.Fatalf("expected handled=true and non-zero code, got handled=%v code=%d", handled, code)
    }
}

func TestHandleFlags_InvalidEmail(t *testing.T) {
    var out, errB bytes.Buffer
    handled, code := HandleFlags([]string{"-add", "-name", "Bob", "-email", "invalid"}, &out, &errB)
    if !handled || code == 0 {
        t.Fatalf("expected handled=true and non-zero code, got handled=%v code=%d", handled, code)
    }
}

func TestHandleAction_DefaultAndQuit(t *testing.T) {
    if cont := HandleAction(99); !cont {
        t.Fatalf("expected continue=true for unknown action")
    }
    if cont := HandleAction(4); cont {
        t.Fatalf("expected continue=false for quit action")
    }
}


