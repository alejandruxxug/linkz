//go:build !darwin && !windows

package host

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// platformMemory reads /proc/meminfo: MemTotal and MemAvailable (kB).
func platformMemory() (total, free uint64, err error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	return parseMeminfo(string(data))
}

func parseMeminfo(text string) (total, free uint64, err error) {
	found := 0
	for _, line := range strings.Split(text, "\n") {
		name, rest, ok := strings.Cut(line, ":")
		if !ok || (name != "MemTotal" && name != "MemAvailable") {
			continue
		}
		fields := strings.Fields(rest)
		if len(fields) == 0 {
			continue
		}
		kb, perr := strconv.ParseUint(fields[0], 10, 64)
		if perr != nil {
			return 0, 0, perr
		}
		if name == "MemTotal" {
			total = kb * 1024
		} else {
			free = kb * 1024
		}
		found++
	}
	if found != 2 {
		return 0, 0, fmt.Errorf("unexpected /proc/meminfo")
	}
	return total, free, nil
}
