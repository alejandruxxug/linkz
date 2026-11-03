package host

import (
	"strings"
	"testing"

	"lizardscript/object"
)

// Both engines give the same value for the same program, and the same
// kind of error.
func TestEvalEngines(t *testing.T) {
	prev := CurrentEngine
	defer func() { CurrentEngine = prev }()

	tests := []struct {
		name, src string
		want      string
	}{
		{"value", "let a = 6; a * 7", "42"},
		{"loop", "let s = 0; for (let i = 0; i < 10; i += 1) { s += i }; s", "45"},
		{"error", "1 + true", "type mismatch"},
	}
	for _, engine := range []Engine{EngineVM, EngineTree} {
		for _, tt := range tests {
			t.Run(string(engine)+"/"+tt.name, func(t *testing.T) {
				CurrentEngine = engine
				prog, err := parseNamed(tt.src, "")
				if err != nil {
					t.Fatal(err)
				}
				res := Eval(prog, object.NewEnvironment())
				if e, ok := res.(*object.Error); ok {
					if tt.want != "type mismatch" || !strings.Contains(e.Message, tt.want) {
						t.Fatalf("error %q, want %q", e.Message, tt.want)
					}
					return
				}
				if res.Inspect() != tt.want {
					t.Errorf("got %s, want %s", res.Inspect(), tt.want)
				}
			})
		}
	}
}
