package host

import (
	"image"
	stdcolor "image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"lizardscript/object"
)

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestListFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "b.lzdraw"), "bb")
	writeFile(t, filepath.Join(dir, "a.lzdraw"), "a")
	writeFile(t, filepath.Join(dir, ".hidden.lzdraw"), "h")
	writeFile(t, filepath.Join(dir, "notes.txt"), "n")
	os.Mkdir(filepath.Join(dir, "sub.lzdraw"), 0o755) // a folder with the extension
	when := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	os.Chtimes(filepath.Join(dir, "a.lzdraw"), when, when)

	got, err := ListFiles(dir, ".lzdraw")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Name != "a.lzdraw" || got[1].Name != "b.lzdraw" {
		t.Fatalf("ListFiles = %+v, want a.lzdraw then b.lzdraw", got)
	}
	if got[0].MTime != when.Unix() || got[1].Size != 2 || got[0].Path != filepath.Join(dir, "a.lzdraw") {
		t.Errorf("wrong details: %+v", got)
	}
	if all, _ := ListFiles(dir, ""); len(all) != 3 {
		t.Errorf("no extension should list every visible file, got %d", len(all))
	}
	if none, err := ListFiles(filepath.Join(dir, "missing"), ".lzdraw"); err != nil || len(none) != 0 {
		t.Errorf("a missing folder is an empty list, got %v, %v", none, err)
	}

	// The same through lizard.
	h := fsHash()
	res, ok := call(t, h, "list", newStr(dir), newStr(".lzdraw")).(*object.Array)
	if !ok || len(res.Elements) != 2 {
		t.Fatalf("fs.list = %v", res)
	}
	row := res.Elements[0].(*object.Hash)
	if hashGet(row, "name").Inspect() != "a.lzdraw" || hashGet(row, "mtime").Inspect() != "1793534400" {
		t.Errorf("row = %s", row.Inspect())
	}
	if call(t, h, "list", newStr(dir)).Type() != "ERROR" {
		t.Error("fs.list needs two arguments")
	}
}

func TestFSHashFile(t *testing.T) {
	dir := t.TempDir()
	a, b, c := filepath.Join(dir, "a"), filepath.Join(dir, "b"), filepath.Join(dir, "c")
	writeFile(t, a, "same text")
	writeFile(t, b, "same text")
	writeFile(t, c, "other text")
	ha, _ := HashFile(a)
	hb, _ := HashFile(b)
	hc, _ := HashFile(c)
	if len(ha) != 16 || ha != hb || ha == hc {
		t.Errorf("hashes: %q %q %q", ha, hb, hc)
	}
	if got := call(t, fsHash(), "hash", newStr(a)).Inspect(); got != ha {
		t.Errorf("fs.hash = %q, want %q", got, ha)
	}
	if call(t, fsHash(), "hash", newStr(filepath.Join(dir, "nope"))).Type() != "ERROR" {
		t.Error("hashing a missing file should be an error")
	}
}

func TestMkdirRenameCopy(t *testing.T) {
	dir := t.TempDir()
	h := fsHash()
	deep := filepath.Join(dir, "a", "b")
	if call(t, h, "mkdir", newStr(deep)).Type() == "ERROR" {
		t.Fatal("mkdir failed")
	}
	if info, err := os.Stat(deep); err != nil || !info.IsDir() {
		t.Error("mkdir did not make the folders")
	}
	call(t, h, "mkdir", newStr(deep)) // again is fine

	one, two := filepath.Join(dir, "one.lzdraw"), filepath.Join(dir, "two.lzdraw")
	writeFile(t, one, "one")
	writeFile(t, two, "two")
	if call(t, h, "rename", newStr(one), newStr(two)).Type() != "ERROR" {
		t.Error("rename over an existing file should be refused")
	}
	if got, _ := os.ReadFile(two); string(got) != "two" {
		t.Error("the existing file was changed")
	}
	three := filepath.Join(dir, "three.lzdraw")
	if call(t, h, "rename", newStr(one), newStr(three)).Type() == "ERROR" {
		t.Fatal("rename failed")
	}
	if _, err := os.Stat(one); err == nil {
		t.Error("the old name is still there")
	}
	if got, _ := os.ReadFile(three); string(got) != "one" {
		t.Errorf("renamed content = %q", got)
	}

	copy := filepath.Join(dir, "copy.lzdraw")
	if call(t, h, "copy", newStr(three), newStr(copy)).Type() == "ERROR" {
		t.Fatal("copy failed")
	}
	if got, _ := os.ReadFile(copy); string(got) != "one" {
		t.Errorf("copy content = %q", got)
	}
	if got, _ := os.ReadFile(three); string(got) != "one" {
		t.Error("the original changed")
	}
	if call(t, h, "copy", newStr(three), newStr(copy)).Type() != "ERROR" {
		t.Error("copy over an existing file should be refused")
	}
}

func TestTrash(t *testing.T) {
	if _, err := trashDir(); err != nil {
		t.Skip("no Trash folder on this system (Windows uses the Recycle Bin)")
	}
	home := t.TempDir()
	t.Setenv("HOME", home)
	dir := t.TempDir()
	path := filepath.Join(dir, "art.lzdraw")
	writeFile(t, path, "first")
	if call(t, fsHash(), "trash", newStr(path)).Type() == "ERROR" {
		t.Fatal("trash failed")
	}
	if _, err := os.Stat(path); err == nil {
		t.Error("the file is still where it was")
	}
	trash, _ := trashDir()
	if got, err := os.ReadFile(filepath.Join(trash, "art.lzdraw")); err != nil || string(got) != "first" {
		t.Fatalf("the file is not in the Trash: %v", err)
	}
	// A second file with the same name does not replace the first.
	writeFile(t, path, "second")
	if err := Trash(path); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(filepath.Join(trash, "art.lzdraw")); string(got) != "first" {
		t.Error("the first file in the Trash was replaced")
	}
	if got, _ := os.ReadFile(filepath.Join(trash, "art 2.lzdraw")); string(got) != "second" {
		t.Error("the second file did not get a new name")
	}
	if err := Trash(filepath.Join(dir, "missing.lzdraw")); err == nil {
		t.Error("trashing a missing file should fail")
	}
}

func TestReveal(t *testing.T) {
	old := runCommand
	defer func() { runCommand = old }()
	var ran []string
	runCommand = func(name string, args ...string) error {
		ran = append(ran, name+" "+strings.Join(args, " "))
		return nil
	}
	path := filepath.Join(t.TempDir(), "a.lzdraw")
	writeFile(t, path, "x")
	if call(t, osHash(), "reveal", newStr(path)).Type() == "ERROR" {
		t.Fatal("reveal failed")
	}
	if len(ran) != 1 || !strings.Contains(ran[0], path) {
		t.Errorf("ran %v, want a command that names the file", ran)
	}
	if call(t, osHash(), "reveal", newStr(path+"x")).Type() != "ERROR" {
		t.Error("revealing a missing file should be an error")
	}
	if home := call(t, osHash(), "home").Inspect(); home == "" {
		t.Error("os.home is empty")
	}
}

func TestTimePrims(t *testing.T) {
	h := timeHash()
	now := call(t, h, "now").(*object.Integer).Value
	if d := time.Now().Unix() - now; d < 0 || d > 2 {
		t.Errorf("time.now = %d is not now", now)
	}
	if tz := call(t, h, "tz"); tz.Type() != "INTEGER" {
		t.Errorf("time.tz = %s", tz.Inspect())
	}
}

// A kept PNG is drawn back with blit, scaled by nearest neighbor.
func TestLoadPNGAndBlit(t *testing.T) {
	f := NewFBUI(8, 8)
	f.off = NewBuf(2, 2) // red, green / blue, transparent
	b := f.off
	set := func(x, y int, c stdcolor.RGBA) { b.Img.SetRGBA(x, y, c) }
	set(0, 0, stdcolor.RGBA{255, 0, 0, 255})
	set(1, 0, stdcolor.RGBA{0, 255, 0, 255})
	set(0, 1, stdcolor.RGBA{0, 0, 255, 255})
	path := filepath.Join(t.TempDir(), "t.png")
	if err := f.SavePNG(path, nil); err != nil {
		t.Fatal(err)
	}
	f.off = nil

	fb := f.Builtins()
	sz := call(t, fb, "load_png", newStr(path), newStr("t"))
	if sz.Inspect() != "(2, 2)" {
		t.Fatalf("load_png gave %s, want (2, 2)", sz.Inspect())
	}
	call(t, fb, "fill", newPoint(0, 0), newPoint(8, 8), color(10, 10, 10, 255))
	call(t, fb, "blit", newStr("t"), newPoint(1, 1), newPoint(4, 4)) // twice as big
	px := func(x, y int) [3]uint8 {
		p := f.ui.Img.Pix[y*f.ui.Img.Stride+x*4:]
		return [3]uint8{p[0], p[1], p[2]}
	}
	for _, c := range []struct {
		x, y int
		want [3]uint8
	}{{1, 1, [3]uint8{255, 0, 0}}, {2, 2, [3]uint8{255, 0, 0}}, {3, 1, [3]uint8{0, 255, 0}},
		{1, 3, [3]uint8{0, 0, 255}}, {4, 4, [3]uint8{10, 10, 10}}, {0, 0, [3]uint8{10, 10, 10}}} {
		if got := px(c.x, c.y); got != c.want {
			t.Errorf("pixel (%d,%d) = %v, want %v", c.x, c.y, got, c.want)
		}
	}
	// The transparent corner of the picture shows what is below it.
	if got := px(3, 3); got != [3]uint8{10, 10, 10} {
		t.Errorf("transparent pixel changed the background: %v", got)
	}
	// 1:1 when no size is given, and the clip is respected.
	call(t, fb, "clip", newPoint(0, 0), newPoint(1, 8))
	call(t, fb, "blit", newStr("t"), newPoint(0, 5))
	if got := px(0, 5); got != [3]uint8{255, 0, 0} {
		t.Errorf("1:1 blit pixel = %v", got)
	}
	if got := px(1, 5); got != [3]uint8{10, 10, 10} {
		t.Errorf("the blit ignored the clip: %v", got)
	}
	for _, bad := range [][]object.Object{
		{newStr("none"), newPoint(0, 0)},
		{newStr("t")},
		{newStr("t"), newPoint(0, 0), newPoint(-1, 2)},
	} {
		if call(t, fb, "blit", bad...).Type() != "ERROR" {
			t.Errorf("blit %v should be an error", bad)
		}
	}
	if call(t, fb, "load_png", newStr(filepath.Join(t.TempDir(), "no.png")), newStr("x")).Type() != "ERROR" {
		t.Error("loading a missing PNG should be an error")
	}
}

func TestBlitAlphaBlends(t *testing.T) {
	dst := NewBuf(1, 1)
	dst.Img.SetRGBA(0, 0, stdcolor.RGBA{100, 100, 100, 255})
	src := NewBuf(1, 1)
	src.Img.SetRGBA(0, 0, stdcolor.RGBA{128, 0, 0, 128}) // half red, premultiplied
	blitBuf(dst, src, 0, 0, 1, 1)
	got := dst.Img.RGBAAt(0, 0)
	if got.R < 177 || got.R > 179 || got.G < 49 || got.G > 51 || got.A != 255 {
		t.Errorf("half-transparent red over gray = %v", got)
	}
	var _ image.Image = dst.Img
}
