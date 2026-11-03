// Package lib bundles Linkz's libraries into the binary. The libraries are
// written in lizard-script; this file only embeds them.
package lib

import (
	"embed"
	"sort"
	"strings"
)

// Files holds the source of the graphics support (gfx) and the editor (ui).
// The drawing libraries live in their own package, lib/draw.
//
//go:embed gfx ui
var Files embed.FS

// Dir lists the .liz files of a folder, in the order they load: by name,
// except that names in first load before the rest. A folder that isn't
// embedded is empty.
func Dir(dir string, first ...string) []string {
	entries, err := Files.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".liz") {
			names = append(names, e.Name())
		}
	}
	rank := func(n string) int {
		for i, f := range first {
			if f == n {
				return i
			}
		}
		return len(first)
	}
	sort.SliceStable(names, func(i, j int) bool {
		if ri, rj := rank(names[i]), rank(names[j]); ri != rj {
			return ri < rj
		}
		return names[i] < names[j]
	})
	for i := range names {
		names[i] = dir + "/" + names[i]
	}
	return names
}

// UIFiles are the libraries the editor's UI environment loads: the graphics
// support first, then the UI, with app.liz last.
func UIFiles() []string {
	files := Dir("gfx")
	var app []string
	for _, f := range Dir("ui") {
		if f == "ui/app.liz" {
			app = append(app, f)
		} else {
			files = append(files, f)
		}
	}
	return append(files, app...)
}
