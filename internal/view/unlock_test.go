package view

import (
	"errors"
	"gototp/internal/gototp"
	"gototp/internal/models"
	"runtime"
	"runtime/secret"
	"testing"
)

func TestUnlockInputRunsInSecretMode(t *testing.T) {
	wantErr := errors.New("input aborted")
	wantSecret := runtime.GOOS == "linux" && (runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64")
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)

	unlocked, err := unlock(func() (string, error) {
		if secret.Enabled() != wantSecret {
			t.Error("unexpected secret mode during password input")
		}
		childSecret := make(chan bool, 1)
		go func() { childSecret <- secret.Enabled() }()
		if <-childSecret != wantSecret {
			t.Error("password input goroutines did not inherit secret mode")
		}
		return "", wantErr
	})
	runtime.ReadMemStats(&after)
	if unlocked != nil || !errors.Is(err, wantErr) {
		t.Fatalf("unlock = (%v, %v), want (nil, input error)", unlocked, err)
	}
	if after.NumGC <= before.NumGC {
		t.Error("unlock did not force garbage collection on input error")
	}
	if secret.Enabled() {
		t.Error("secret mode leaked into caller")
	}
}

func TestUnlockExistingDatabase(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("config uses AppData instead of XDG_CONFIG_HOME on Windows")
	}
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	const passphrase = "test-only database passphrase"
	// Create a fixture using the existing database format/derivation path.
	original, err := gototp.New(passphrase)
	if err != nil {
		t.Fatal(err)
	}
	original.Data.Keystore.Add(models.Input{
		Name: "example", Period: "30", Digit: "6", Algorithm: "SHA1",
		Secret: "JBSWY3DPEHPK3PXP",
	})
	if err := original.Save(); err != nil {
		t.Fatal(err)
	}

	unlocked, err := unlock(func() (string, error) { return passphrase, nil })
	if err != nil {
		t.Fatal(err)
	}
	if len(unlocked.Data.Keystore.Keys) != 1 || unlocked.Data.Keystore.Keys[0] != original.Data.Keystore.Keys[0] {
		t.Fatal("unlocked database differs from original")
	}
	// The derived key and database must survive the forced secret-heap GC.
	if err := unlocked.Save(); err != nil {
		t.Fatal(err)
	}
	if err := unlocked.Load(); err != nil {
		t.Fatal(err)
	}
	if len(unlocked.Data.Keystore.Keys) != 1 {
		t.Fatal("database did not survive save/reload")
	}

	wrong, err := unlock(func() (string, error) { return "wrong passphrase", nil })
	if wrong != nil || err == nil {
		t.Fatal("wrong passphrase must not unlock the database")
	}
}
