// Package draw bundles the drawing libraries: shapes, pens, geometry and the
// SVG writer. They are written in lizard-script; this file only embeds them.
package draw

import "embed"

// Files holds the source of every drawing library.
//
//go:embed *.liz
var Files embed.FS

// Order is the order the libraries load in, since a library can use what the
// ones before it define. here.liz must come first: the host replaces its
// here() with the real one after loading it.
var Order = []string{"here.liz", "geometry.liz", "pens.liz", "shapes.liz", "gradient.liz", "blend.liz", "svg.liz"}
