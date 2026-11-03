//go:build !darwin

package pressure

// Start does nothing: only macOS reports pressure.
func Start() bool { return false }
