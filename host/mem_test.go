package host

import (
	"runtime"
	"testing"

	"lizardscript/object"
)

func TestReadMemory(t *testing.T) {
	m, err := ReadMemory()
	if err != nil {
		t.Fatal(err)
	}
	if m.Total == 0 || m.Free == 0 {
		t.Errorf("memory %+v: want both above 0", m)
	}
	if m.Free > m.Total {
		t.Errorf("free %d is more than total %d", m.Free, m.Total)
	}
	if m.Total < 64<<20 || m.Total > 64<<40 {
		t.Errorf("total %d bytes is not plausible", m.Total)
	}
}

func TestOSMemoryAndCPUs(t *testing.T) {
	h := osHash()
	call := func(name string) object.Object {
		return hashGet(h, name).(*object.Builtin).Fn()
	}
	mem, ok := call("memory").(*object.Hash)
	if !ok {
		t.Fatal("os.memory() is not a hash")
	}
	total := hashGet(mem, "total").(*object.Integer).Value
	free := hashGet(mem, "free").(*object.Integer).Value
	if total <= 0 || free <= 0 || free > total {
		t.Errorf("os.memory() = total %d, free %d", total, free)
	}
	if n := call("cpus").(*object.Integer).Value; n < 1 || n > 4096 || int(n) != runtime.NumCPU() {
		t.Errorf("os.cpus() = %d", n)
	}
	old := AppVersion
	defer func() { AppVersion = old }()
	AppVersion = "9.9.9"
	if got := call("app_version").(*object.String).Value; got != "9.9.9" {
		t.Errorf("os.app_version() = %q", got)
	}
}
