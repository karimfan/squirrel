package main

import (
	"os"
	"path/filepath"
	"testing"

	"net/http/httptest"

	"squirrel/internal/server"
)

func TestEndToEnd(t *testing.T) {
	ms := server.NewMemoryStore()
	srv := server.NewServer(ms)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	serverURL = ts.URL + "/graphql"
	dataDir = t.TempDir()
	sessionFile = filepath.Join(dataDir, "session_token")

	if err := runCreateAccount([]string{"--name", "Jane", "--email", "jane@example.com", "--password", "secret", "--squirrel", "Sprinkles"}); err != nil {
		t.Fatalf("create account: %v", err)
	}

	if err := runLogin([]string{"--email", "jane@example.com", "--password", "secret"}); err != nil {
		t.Fatalf("login: %v", err)
	}

	tok, err := os.ReadFile(filepath.Join(dataDir, "session_token"))
	if err != nil {
		t.Fatalf("read token: %v", err)
	}
	if string(tok) != "token-1" {
		t.Fatalf("unexpected token %s", string(tok))
	}

	if err := runAddItem([]string{"Buy milk", "--tags", "errands"}); err != nil {
		t.Fatalf("add note: %v", err)
	}
	if ms.Items[0].Type != "note" {
		t.Fatalf("expected note got %s", ms.Items[0].Type)
	}

	if err := runAddItem([]string{"https://example.com", "--tags", "articles"}); err != nil {
		t.Fatalf("add article: %v", err)
	}
	if ms.Items[1].Type != "article" {
		t.Fatalf("expected article got %s", ms.Items[1].Type)
	}
}
