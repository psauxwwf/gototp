//go:build linux

package security

import (
	"os"
	"os/exec"
	"testing"

	"golang.org/x/sys/unix"
)

func TestDisableDumps(t *testing.T) {
	// The hard core limit cannot be restored by an unprivileged process.
	// Apply it only in a child, not in the test runner.
	if os.Getenv("GOTOTP_TEST_DISABLE_DUMPS") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestDisableDumps$")
		cmd.Env = append(os.Environ(), "GOTOTP_TEST_DISABLE_DUMPS=1")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("child process: %v\n%s", err, output)
		}
		return
	}

	if err := DisableDumps(); err != nil {
		t.Fatal(err)
	}
	dumpable, err := unix.PrctlRetInt(unix.PR_GET_DUMPABLE, 0, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if dumpable != 0 {
		t.Fatalf("process is still dumpable: %d", dumpable)
	}
	var limit unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_CORE, &limit); err != nil {
		t.Fatal(err)
	}
	if limit.Cur != 0 || limit.Max != 0 {
		t.Fatalf("core limits = %+v, want zero", limit)
	}
}
