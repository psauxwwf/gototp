package main

import (
	"fmt"
	"os"

	"gototp/internal/security"
	"gototp/internal/view"
)

func main() {
	if err := security.DisableDumps(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_view := view.New()
	_view.Run()
}
