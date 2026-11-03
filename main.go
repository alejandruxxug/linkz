// Command linkz draws vector art written in lizard-script.
//
//	linkz edit [--no-splash] [drawing.lzdraw]  opens the native editor
//	linkz render drawing.lzdraw [out.svg]  writes the SVG
//	linkz trace in.png out.lzdraw [--colors N] [--detail D]  traces a picture into code
//	linkz paint prog.liz out.png [WxH]     runs a program headless on a buffer
//	linkz debug [--pid N] Command [args]   drives a running editor
//	linkz calibrate                        measures this machine, prints JSON
//	linkz mcp                              MCP server on stdin/stdout, for AI agents
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"linkz/host"
	"linkz/host/window"
)

const usage = `usage:
  linkz edit [--no-splash] [file.lzdraw]
  linkz render file.lzdraw [out.svg]
  linkz trace in.png out.lzdraw [--colors 8] [--detail 1.5]
  linkz paint prog.liz out.png [WxH]
  linkz debug [--pid N] Command [args...]
  linkz calibrate
  linkz mcp`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "edit":
		err = edit(os.Args[2:])
	case "render":
		err = render(os.Args[2:])
	case "trace":
		err = traceCmd(os.Args[2:])
	case "paint":
		err = paint(os.Args[2:])
	case "debug":
		err = debug(os.Args[2:])
	case "calibrate":
		err = calibrate()
	case "mcp":
		err = host.ServeMCP(os.Stdin, os.Stdout, os.Stderr)
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// version is the app's version. A new one measures the machine again.
// Release builds set it with -ldflags "-X main.version=...".
var version = "0.1.0-dev"

// calibrate measures this machine, keeps the result in ~/.linkz/machine.liz
// and prints it as JSON.
func calibrate() error {
	host.AppVersion = version
	text, err := host.Recalibrate(os.Stderr)
	if err != nil {
		return err
	}
	fmt.Println(text)
	return nil
}

func render(args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("usage: linkz render drawing.lzdraw [out.svg]")
	}
	in := args[0]
	out := strings.TrimSuffix(in, ".lzdraw") + ".svg"
	if len(args) == 2 {
		out = args[1]
	}
	source, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	svg, err := host.Render(string(source), os.Stdout)
	if err != nil {
		return err
	}
	if err := os.WriteFile(out, []byte(svg), 0o644); err != nil {
		return err
	}
	fmt.Println("wrote", out)
	return nil
}

func paint(args []string) error {
	if len(args) < 2 || len(args) > 3 {
		return fmt.Errorf("usage: linkz paint prog.liz out.png [WxH]")
	}
	w, h := 512, 512
	if len(args) == 3 {
		ws, hs, ok := strings.Cut(strings.ToLower(args[2]), "x")
		var e1, e2 error
		w, e1 = strconv.Atoi(ws)
		h, e2 = strconv.Atoi(hs)
		if !ok || e1 != nil || e2 != nil || w <= 0 || h <= 0 || w*h > 1<<26 {
			return fmt.Errorf("size must look like 512x512")
		}
	}
	source, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	return host.Paint(string(source), args[1], w, h, os.Stdout)
}

func edit(args []string) error {
	path, splash := "", true
	for _, a := range args {
		switch {
		case a == "--no-splash":
			splash = false
		case strings.HasPrefix(a, "--"):
			return fmt.Errorf("unknown flag %s (usage: linkz edit [--no-splash] [file.lzdraw])", a)
		default:
			path = a
		}
	}
	app := host.NewApp(1024, 720)
	settings, err := host.LoadSettings()
	if err != nil {
		fmt.Fprintln(os.Stderr, "linkz: settings ignored, using defaults:", err)
	}
	app.Settings = settings
	host.AppVersion = version
	if machine, err := host.Machine(false, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "linkz: machine not measured:", err)
	} else {
		app.Machine = machine
	}
	if host.SkipLoadingAnim(settings) {
		splash = false
	}
	dbg := host.NewDebug(app)

	// A kill must not leave its debug file behind.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		dbg.Stop()
		os.Exit(130)
	}()
	defer dbg.Stop()

	title := "Linkz"
	if path != "" {
		title += " - " + path
	}

	var handler window.Handler = app
	win := window.Ebiten{}
	var launcher *host.Launcher
	if splash {
		sp, err := host.NewSplash()
		if err != nil {
			fmt.Fprintln(os.Stderr, "linkz: no splash:", err)
			splash = false
		} else {
			launcher = host.NewLauncher(sp, app, path)
			launcher.OnReady = func() {
				if err := dbg.Start(); err != nil {
					fmt.Fprintln(os.Stderr, "linkz: backdoor:", err)
				}
			}
			handler, win = launcher, window.Ebiten{Frameless: true}
		}
	}
	w, h := 1024, 720
	if splash {
		w, h = launcher.Splash.WindowSize()
	} else {
		if err := app.Load(path); err != nil {
			return err
		}
		if err := dbg.Start(); err != nil {
			return err
		}
	}
	err = win.Run(title, w, h, handler)
	time.Sleep(150 * time.Millisecond) // let the Quit reply go out
	if err == nil && launcher != nil {
		err = launcher.Err
	}
	return err
}

// traceCmd turns a picture into drawing code: linkz trace in.png out.lzdraw
// [--colors N] [--detail D] [--min-area N] [--max-side N].
func traceCmd(args []string) error {
	const use = "usage: linkz trace in.png out.lzdraw [--colors 8] [--detail 1.5] [--min-area 8] [--max-side 256] [--scale 1] [--corner 55]"
	var files []string
	options := map[string]any{}
	for i := 0; i < len(args); i++ {
		name, ok := strings.CutPrefix(args[i], "--")
		if !ok {
			files = append(files, args[i])
			continue
		}
		key := map[string]string{"colors": "colors", "detail": "detail", "min-area": "min_area", "max-side": "max_side", "scale": "scale", "corner": "corner"}[name]
		if key == "" || i+1 >= len(args) {
			return fmt.Errorf("%s\n%s", "bad option "+args[i], use)
		}
		i++
		n, err := strconv.ParseFloat(args[i], 64)
		if err != nil {
			return fmt.Errorf("--%s needs a number, got %q", name, args[i])
		}
		options[key] = n // whole numbers become INTEGERs in lizard
	}
	if len(files) != 2 {
		return fmt.Errorf("%s", use)
	}
	code, stats, err := host.Trace(files[0], options, os.Stdout)
	if err != nil {
		return err
	}
	if err := os.WriteFile(files[1], []byte(code), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s: %v colors, %v regions (%v holes), %v points, %.0f ms\n",
		files[1], stats["colors"], stats["regions"], stats["holes"], stats["points"], stats["ms"])
	if st, ok := stats["stages"].(map[string]any); ok {
		fmt.Printf("  load %.0f, quantize %.0f, clean %.0f, rings %.0f, simplify %.0f, code %.0f ms\n",
			st["load"], st["quantize"], st["clean"], st["rings"], st["simplify"], st["code"])
	}
	return nil
}
