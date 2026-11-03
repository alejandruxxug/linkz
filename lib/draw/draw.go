// Package draw bundles the drawing libraries: shapes, pens, geometry and the
// SVG writer. They are written in lizard-script; this file only embeds them.
package draw

import (
	"embed"

	"lizardscript/evaluator"
)

// Files holds the source of every drawing library.
//
//go:embed *.liz
var Files embed.FS

// Docs holds the translated doc comments, one file per language
// (docs/es.liz defines docs_es). They are data for the tooltips, not library
// code: they are never loaded into a drawing's environment.
//
//go:embed docs/*.liz
var Docs embed.FS

// Order is the order the libraries load in, since a library can use what the
// ones before it define. here.liz must come first: the host replaces its
// here() with the real one after loading it.
var Order = []string{"here.liz", "geometry.liz", "pens.liz", "check.liz", "shapes.liz", "gradient.liz", "blend.liz", "svg.liz"}

// The libraries say which others they need with use lines, for the lizard
// port (lizardscript run .). The Go host loads every library into one
// environment, in the orders kept here and in lib.go, so those uses are
// already satisfied there: for it, each module reads as an empty file.
// It goes away with the Go host (B042).
func init() {
	evaluator.ReadModule = func(string) ([]byte, error) { return nil, nil }
}
