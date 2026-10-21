package view

import (
	"gototp/internal/gototp"

	"runtime"
	"runtime/secret"
)

func readPassphrase() (string, error) {
	var passphrase string
	form := PasswordForm(
		"🔐 enter your secret", "",
		strValidator,
		&passphrase,
	)
	if err := form.Run(); err != nil {
		return "", err
	}
	return passphrase, nil
}

func unlock(read func() (string, error)) (*gototp.Gototp, error) {
	var (
		unlocked *gototp.Gototp
		err      error
	)
	// The input form must also run in secret mode: it creates its own copies
	// of the password, including in Bubble Tea's child goroutines.
	// Do erases temporary stacks/registers; GC erases unreachable secret heap
	// allocations before we enter the main menu, including on input errors.
	defer runtime.GC()
	secret.Do(func() {
		passphrase, readErr := read()
		if readErr != nil {
			err = readErr
			return
		}
		unlocked, err = gototp.New(passphrase)
	})
	return unlocked, err
}
