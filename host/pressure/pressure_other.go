//go:build !darwin && !windows

package pressure

// Start does nothing: only macOS and Windows report pressure.
func Start() bool { return false }
