// Package lib bundles Linkz's libraries into the binary. The libraries are
// written in lizard-script; this file only embeds them.
package lib

import (
	"embed"
	"sort"
	"strings"

	"linkz/lib/draw"
)

// Files holds the source of the graphics support (gfx) and the editor (ui).
// Only the top-level .liz files of gfx are embedded, not its tests. The
// drawing libraries live in their own package, lib/draw.
//
//go:embed gfx/*.liz gfx/FONT_LICENSE ui
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

// DrawPrefix marks a UIFiles name that lives in lib/draw, not in Files.
const DrawPrefix = "draw/"

// splashFiles belong to the splash screen, which has its own environment;
// the editor's environment doesn't load them.
var splashFiles = map[string]bool{"ui/splash.liz": true, "ui/splash_art.liz": true}

// UIFiles are the libraries the editor's UI environment loads, in order: the
// drawing libraries (geometry, shapes, paint_order), which the graphics
// support uses, then the graphics support, then the UI, with app.liz last.
// Names starting with DrawPrefix are read with ReadUI.
func UIFiles() []string {
	var files []string
	for _, f := range draw.Order {
		files = append(files, DrawPrefix+f)
	}
	files = append(files, Dir("gfx")...)
	var app []string
	for _, f := range Dir("ui") {
		switch {
		case splashFiles[f]:
		case f == "ui/app.liz":
			app = append(app, f)
		default:
			files = append(files, f)
		}
	}
	return append(files, app...)
}

// ReadUI reads a file named by UIFiles.
func ReadUI(name string) ([]byte, error) {
	if rest, ok := strings.CutPrefix(name, DrawPrefix); ok {
		return draw.Files.ReadFile(rest)
	}
	return Files.ReadFile(name)
}
