// Command linkz draws vector art written in lizard-script.
//
//	linkz render drawing.lzdraw            writes drawing.svg
//	linkz render drawing.lzdraw out.svg    writes out.svg
package main

import (
	"fmt"
	"linkz/host"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 || os.Args[1] != "render" {
		fmt.Fprintln(os.Stderr, "usage: linkz render drawing.lzdraw [out.svg]")
		os.Exit(2)
	}

	in := os.Args[2]
	out := strings.TrimSuffix(in, ".lzdraw") + ".svg"
	if len(os.Args) == 4 {
		out = os.Args[3]
	}

	source, err := os.ReadFile(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	svg, err := host.Render(string(source), os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.WriteFile(out, []byte(svg), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("wrote", out)
}
