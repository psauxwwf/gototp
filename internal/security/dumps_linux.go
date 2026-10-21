//go:build linux

package security

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// DisableDumps must run before reading any secrets. PR_SET_DUMPABLE also
// restricts unprivileged ptrace and /proc/PID/mem access, unlike RLIMIT_CORE
// alone (which does not prevent piped core dumps).
func DisableDumps() error {
	if err := unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0); err != nil {
		return fmt.Errorf("disable process dumps: %w", err)
	}
	if err := unix.Setrlimit(unix.RLIMIT_CORE, &unix.Rlimit{}); err != nil {
		return fmt.Errorf("disable core files: %w", err)
	}
	return nil
}
