//go:build !linux

package security

// DisableDumps currently only hardens Linux processes. Other platforms retain
// their existing behavior; runtime/secret also has platform limitations.
func DisableDumps() error {
	return nil
}
