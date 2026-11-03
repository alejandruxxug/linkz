// Package trace bundles the image tracing library: it turns a picture into
// lizard-script drawing code. It is written in lizard-script; this file only
// embeds it. It needs the drawing libraries (lib/draw) loaded first, for the
// geometry helpers, and `img` from the host to read files.
package trace

import "embed"

// Files holds the source of the tracing library.
//
//go:embed *.liz
var Files embed.FS

// Order is the order the files load in.
var Order = []string{"trace.liz"}
