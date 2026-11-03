package host

import (
	"runtime"

	"lizardscript/object"
)

// os.memory() and os.cpus(): what the machine has, so the editor can pick a
// sensible object limit. The numbers come from each OS (mem_*.go).

// MemoryInfo is the machine's RAM in bytes: all of it, and what is free to
// use now (free plus what the system can give back at once).
type MemoryInfo struct {
	Total, Free uint64
}

// ReadMemory asks the OS. Free is never more than Total.
func ReadMemory() (MemoryInfo, error) {
	total, free, err := platformMemory()
	if err != nil {
		return MemoryInfo{}, err
	}
	if free > total {
		free = total
	}
	return MemoryInfo{total, free}, nil
}

func addMemPrims(h *object.Hash) {
	hashSet(h, "memory", &object.Builtin{Name: "os.memory", Fn: func(args ...object.Object) object.Object {
		m, err := ReadMemory()
		if err != nil {
			return errorf("os.memory: %v", err)
		}
		out := object.NewHash()
		hashSet(out, "total", &object.Integer{Value: int64(m.Total)})
		hashSet(out, "free", &object.Integer{Value: int64(m.Free)})
		return out
	}})
	hashSet(h, "cpus", &object.Builtin{Name: "os.cpus", Fn: func(args ...object.Object) object.Object {
		return newInt(runtime.NumCPU())
	}})
	hashSet(h, "app_version", &object.Builtin{Name: "os.app_version", Fn: func(args ...object.Object) object.Object {
		return newStr(AppVersion)
	}})
}
