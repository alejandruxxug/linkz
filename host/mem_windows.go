//go:build windows

package host

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// memoryStatusEx is the MEMORYSTATUSEX struct of GlobalMemoryStatusEx.
type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

func platformMemory() (total, free uint64, err error) {
	proc := windows.NewLazySystemDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx")
	st := memoryStatusEx{Length: uint32(unsafe.Sizeof(memoryStatusEx{}))}
	if r, _, e := proc.Call(uintptr(unsafe.Pointer(&st))); r == 0 {
		return 0, 0, e
	}
	return st.TotalPhys, st.AvailPhys, nil
}
