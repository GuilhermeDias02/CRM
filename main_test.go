package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// helper entrypoint that invokes main() in a subprocess when GO_WANT_HELPER_PROCESS is set
func TestInvokeMain(t *testing.T) {
    if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
        return
    }
    // extract args after "--" and feed them to main via os.Args
    sep := 0
    for i, a := range os.Args {
        if a == "--" { sep = i + 1; break }
    }
    if sep > 0 && sep < len(os.Args) {
        os.Args = append([]string{os.Args[0]}, os.Args[sep:]...)
    } else {
        os.Args = os.Args[:1]
    }
    main()
    os.Exit(0)
}

func TestMain_HandleFlagsError_ExitsWithCode(t *testing.T) {
    // invalid email -> action.HandleFlags returns handled=true, code=1 -> main os.Exit(1)
    args := []string{"-test.run=TestInvokeMain", "--", "-add", "-name", "Alice", "-email", "invalid"}
    cmd := exec.Command(os.Args[0], args...)
    cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
    err := cmd.Run()
    if err == nil {
        t.Fatalf("expected non-zero exit status")
    }
    if ee, ok := err.(*exec.ExitError); ok {
        if ee.ExitCode() == 0 {
            t.Fatalf("expected exit code != 0, got 0")
        }
    } else {
        t.Fatalf("unexpected error type: %v", err)
    }
}

func TestMain_HandleFlagsSuccess_DisplaysContacts(t *testing.T) {
    // valid add -> main should print contacts and exit 0
    args := []string{"-test.run=TestInvokeMain", "--", "-add", "-name", "Bob", "-email", "bob@example.com"}
    cmd := exec.Command(os.Args[0], args...)
    cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
    out, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("expected exit code 0, got error: %v, output: %s", err, string(out))
    }
    s := string(out)
    if !strings.Contains(s, "Vos contacts:") {
        t.Fatalf("expected contacts to be displayed, got: %s", s)
    }
}
