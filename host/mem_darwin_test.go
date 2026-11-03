//go:build darwin

package host

import "testing"

func TestParseVMStat(t *testing.T) {
	text := `Mach Virtual Memory Statistics: (page size of 16384 bytes)
Pages free:                               10000.
Pages active:                            400000.
Pages inactive:                           50000.
Pages speculative:                         2000.
`
	got, err := parseVMStat(text)
	if err != nil {
		t.Fatal(err)
	}
	if want := uint64(60000 * 16384); got != want {
		t.Errorf("parseVMStat = %d, want %d", got, want)
	}
	if _, err := parseVMStat("nothing useful"); err == nil {
		t.Error("garbage input should be an error")
	}
}
