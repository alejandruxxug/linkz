// Command linkz draws vector art written in lizard-script.
//
//	linkz edit [drawing.lzdraw]            opens the native editor
//	linkz render drawing.lzdraw [out.svg]  writes the SVG
//	linkz paint prog.liz out.png [WxH]     runs a program headless on a buffer
//	linkz debug Command [args]             drives a running editor
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
  linkz edit [file.lzdraw]
  linkz render file.lzdraw [out.svg]
  linkz paint prog.liz out.png [WxH]
  linkz debug Command [args...]`

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
	case "paint":
		err = paint(os.Args[2:])
	case "debug":
		err = debug(os.Args[2:])
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
	path := ""
	if len(args) > 0 {
		path = args[0]
	}
	app := host.NewApp(1024, 720)
	if err := app.Load(path); err != nil {
		return err
	}
	dbg := host.NewDebug(app)
	if err := dbg.Start(); err != nil {
		return err
	}
	defer dbg.Stop()

	// A kill must not leave debug.json behind.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		dbg.Stop()
		os.Exit(130)
	}()

	title := "Linkz"
	if path != "" {
		title += " - " + path
	}
	err := window.Ebiten{}.Run(title, 1024, 720, app)
	time.Sleep(150 * time.Millisecond) // let the Quit reply go out
	return err
}
