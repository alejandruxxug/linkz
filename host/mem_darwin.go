//go:build darwin

package host

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// platformMemory reads the total from sysctl and the free memory from
// vm_stat: free plus inactive pages, times the page size.
func platformMemory() (total, free uint64, err error) {
	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		return 0, 0, err
	}
	total, err = strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	out, err = exec.Command("vm_stat").Output()
	if err != nil {
		return 0, 0, err
	}
	free, err = parseVMStat(string(out))
	return total, free, err
}

// parseVMStat turns vm_stat's text into free bytes (free + inactive pages).
func parseVMStat(text string) (uint64, error) {
	page, pages, found := uint64(0), uint64(0), 0
	for _, line := range strings.Split(text, "\n") {
		switch {
		case strings.Contains(line, "page size of"):
			for _, f := range strings.Fields(line) {
				if n, err := strconv.ParseUint(f, 10, 64); err == nil {
					page = n
				}
			}
		case strings.HasPrefix(line, "Pages free:"), strings.HasPrefix(line, "Pages inactive:"):
			_, num, _ := strings.Cut(line, ":")
			n, err := strconv.ParseUint(strings.Trim(strings.TrimSpace(num), "."), 10, 64)
			if err != nil {
				return 0, err
			}
			pages += n
			found++
		}
	}
	if page == 0 || found != 2 {
		return 0, fmt.Errorf("unexpected vm_stat output")
	}
	return pages * page, nil
}
