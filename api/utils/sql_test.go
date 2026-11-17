package utils

import "testing"

// TestRebind vérifie la conversion des paramètres MySQL (?) en paramètres PostgreSQL ($1, $2...).
func TestRebind(t *testing.T) {
	got := Rebind("SELECT * FROM t WHERE a = ? AND b = ?")
	want := "SELECT * FROM t WHERE a = $1 AND b = $2"
	if got != want {
		t.Fatalf("obtenu %q, attendu %q", got, want)
	}
}
