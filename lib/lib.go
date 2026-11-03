// Package lib bundles Linkz's libraries into the binary. The libraries are
// written in lizard-script; this file only embeds them.
package lib

import "embed"

// Files holds the source of every library.
//
//go:embed *.liz
var Files embed.FS

// Order is the order the libraries load in, since a library can use what the
// ones before it define. This goes away once lizard-script has use.
var Order = []string{"shapes.liz", "svg.liz"}
