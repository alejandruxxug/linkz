module linkz

go 1.27.0

require (
	github.com/ebitengine/purego v0.11.0
	github.com/hajimehoshi/ebiten/v2 v2.10.5
	golang.org/x/image v0.47.0
	golang.org/x/sys v0.49.0
	lizardscript v0.0.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20260820040257-d11f821a26a6 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
)

// lizard-script is the sibling repo. Before the first release this becomes
// a real module path (github.com/alejandruxxug/lizard-script).
replace lizardscript => ../lizard-script
