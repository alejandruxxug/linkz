package host

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"lizardscript/object"
)

// The primitives under the home screen's gallery: listing and moving files,
// the Trash, Finder, the clock, and keeping a PNG as a buffer lizard can
// draw. Nothing here decides what the gallery looks like or when a thumbnail
// is rendered: lizard does the rendering, Go only copies pixels it made.

// addFSPrims adds the file primitives of the gallery to the fs hash.
func addFSPrims(h *object.Hash) {
	add := func(name string, fn object.BuiltinFunction) {
		hashSet(h, name, &object.Builtin{Name: "fs." + name, Fn: fn})
	}
	add("list", func(args ...object.Object) object.Object {
		s, err := strArgs(args, 2, "2 (dir, ext)")
		if err != nil {
			return err
		}
		files, lerr := ListFiles(s[0], s[1])
		if lerr != nil {
			return errorf("%v", lerr)
		}
		out := &object.Array{}
		for _, f := range files {
			out.Elements = append(out.Elements, fromGo(map[string]any{
				"name": f.Name, "path": f.Path, "mtime": f.MTime, "size": f.Size,
			}))
		}
		return out
	})
	add("hash", func(args ...object.Object) object.Object {
		s, err := strArgs(args, 1, "1 (path)")
		if err != nil {
			return err
		}
		sum, herr := HashFile(s[0])
		if herr != nil {
			return errorf("%v", herr)
		}
		return newStr(sum)
	})
	add("mkdir", func(args ...object.Object) object.Object {
		s, err := strArgs(args, 1, "1 (path)")
		if err != nil {
			return err
		}
		if merr := os.MkdirAll(s[0], 0o755); merr != nil {
			return errorf("%v", merr)
		}
		return null
	})
	add("rename", func(args ...object.Object) object.Object {
		s, err := strArgs(args, 2, "2 (from, to)")
		if err != nil {
			return err
		}
		if rerr := RenameFile(s[0], s[1]); rerr != nil {
			return errorf("%v", rerr)
		}
		return null
	})
	add("copy", func(args ...object.Object) object.Object {
		s, err := strArgs(args, 2, "2 (from, to)")
		if err != nil {
			return err
		}
		if cerr := CopyFile(s[0], s[1]); cerr != nil {
			return errorf("%v", cerr)
		}
		return null
	})
	add("trash", func(args ...object.Object) object.Object {
		s, err := strArgs(args, 1, "1 (path)")
		if err != nil {
			return err
		}
		if terr := Trash(s[0]); terr != nil {
			return errorf("%v", terr)
		}
		return null
	})
}

// addOSPrims adds os.home and os.reveal.
func addOSPrims(h *object.Hash) {
	hashSet(h, "home", &object.Builtin{Name: "os.home", Fn: func(args ...object.Object) object.Object {
		home, err := os.UserHomeDir()
		if err != nil {
			return errorf("%v", err)
		}
		return newStr(home)
	}})
	hashSet(h, "reveal", &object.Builtin{Name: "os.reveal", Fn: func(args ...object.Object) object.Object {
		s, err := strArgs(args, 1, "1 (path)")
		if err != nil {
			return err
		}
		if rerr := Reveal(s[0]); rerr != nil {
			return errorf("%v", rerr)
		}
		return null
	}})
}

// addTimePrims adds the wall clock: time.now is seconds since 1970 (an
// INTEGER), time.tz the local offset from UTC in seconds, so lizard can tell
// calendar days apart.
func addTimePrims(h *object.Hash) {
	hashSet(h, "now", &object.Builtin{Name: "time.now", Fn: func(args ...object.Object) object.Object {
		return &object.Integer{Value: time.Now().Unix()}
	}})
	hashSet(h, "tz", &object.Builtin{Name: "time.tz", Fn: func(args ...object.Object) object.Object {
		_, off := time.Now().Zone()
		return newInt(off)
	}})
}

// FileInfo is one row of ListFiles.
type FileInfo struct {
	Name  string // file name with its extension
	Path  string
	MTime int64 // seconds since 1970
	Size  int64
}

// ListFiles lists the regular files of dir whose name ends in ext (any file
// when ext is ""), by name. Hidden files (a leading dot) are skipped. A
// folder that doesn't exist is an empty list, not an error.
func ListFiles(dir, ext string) ([]FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []FileInfo
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") || (ext != "" && !strings.HasSuffix(name, ext)) {
			continue
		}
		info, err := e.Info()
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		out = append(out, FileInfo{name, filepath.Join(dir, name), info.ModTime().Unix(), info.Size()})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// HashFile is a short hash of a file's content: the first 16 hex digits of
// its SHA-256. Same text, same hash, wherever the file is.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil))[:16], nil
}

// RenameFile moves a file, and refuses to replace one that is already there.
func RenameFile(from, to string) error {
	if _, err := os.Lstat(to); err == nil {
		return fmt.Errorf("%s already exists", to)
	}
	return os.Rename(from, to)
}

// CopyFile copies a file, and refuses to replace one that is already there.
func CopyFile(from, to string) error {
	if _, err := os.Lstat(to); err == nil {
		return fmt.Errorf("%s already exists", to)
	}
	data, err := os.ReadFile(from)
	if err != nil {
		return err
	}
	mode := os.FileMode(0o644)
	if info, err := os.Stat(from); err == nil {
		mode = info.Mode().Perm()
	}
	return os.WriteFile(to, data, mode)
}

// runCommand runs a program; tests replace it.
var runCommand = func(name string, args ...string) error {
	return exec.Command(name, args...).Run()
}

// Reveal shows a file in Finder (the folder, elsewhere).
func Reveal(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	if runtime.GOOS == "darwin" {
		return runCommand("open", "-R", path)
	}
	return runCommand("xdg-open", filepath.Dir(path))
}

// trashDir is the folder the Trash is, under the user's home.
func trashDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, ".Trash"), nil
	}
	return filepath.Join(home, ".local", "share", "Trash", "files"), nil
}

// Trash moves a file into the Trash. It never deletes: a name already in the
// Trash gets " 2", " 3", ... before its extension. If the move is not
// allowed it asks Finder to do it, and if that fails too, the file stays.
func Trash(path string) error {
	if _, err := os.Lstat(path); err != nil {
		return err
	}
	dir, err := trashDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err == nil {
		base := filepath.Base(path)
		ext := filepath.Ext(base)
		stem := strings.TrimSuffix(base, ext)
		dest := filepath.Join(dir, base)
		for n := 2; ; n++ {
			if _, err := os.Lstat(dest); err != nil {
				break
			}
			dest = filepath.Join(dir, fmt.Sprintf("%s %d%s", stem, n, ext))
		}
		if err := os.Rename(path, dest); err == nil {
			return nil
		}
	}
	if runtime.GOOS == "darwin" {
		script := fmt.Sprintf(`tell application "Finder" to delete POSIX file %q`, path)
		if err := runCommand("osascript", "-e", script); err == nil {
			return nil
		}
	}
	return fmt.Errorf("could not move %s to the Trash", path)
}

// LoadPNG reads a PNG into a buffer.
func LoadPNG(path string) (*Buf, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", path, err)
	}
	b := NewBuf(img.Bounds().Dx(), img.Bounds().Dy())
	draw.Draw(b.Img, b.Img.Rect, img, img.Bounds().Min, draw.Src)
	return b, nil
}

// loadPNG is fb.load_png(path, name): keeps the PNG under name, like fb.keep
// does for an offscreen buffer, and gives back its size as a point.
func (f *FB) loadPNG(args ...object.Object) object.Object {
	s, err := strArgs(args, 2, "2 (path, name)")
	if err != nil {
		return err
	}
	b, lerr := LoadPNG(s[0])
	if lerr != nil {
		return errorf("%v", lerr)
	}
	f.mu.Lock()
	if f.kept == nil {
		f.kept = map[string]*Buf{}
	}
	f.kept[s[1]] = b
	f.mu.Unlock()
	return newPoint(b.W(), b.H())
}

// blit is fb.blit(name, pos) or fb.blit(name, pos, size): draws a kept buffer
// onto the current target at pos, scaled to size with nearest neighbor (its
// own size when none is given). Transparent pixels leave what is below.
func (f *FB) blit(args ...object.Object) object.Object {
	if len(args) != 2 && len(args) != 3 {
		return errorf("wrong number of arguments: want 2 or 3 (name, pos, size), got %d", len(args))
	}
	name, ok := args[0].(*object.String)
	x, y, okp := pointXY(args[1])
	if !ok || !okp {
		return errorf("want (STRING name, POINT pos[, POINT size])")
	}
	src := f.Kept(name.Value)
	if src == nil {
		return errorf("no kept buffer called %q", name.Value)
	}
	w, h := float64(src.W()), float64(src.H())
	if len(args) == 3 {
		var oks bool
		w, h, oks = pointXY(args[2])
		if !oks || w < 0 || h < 0 {
			return errorf("size must be a POINT with no negative part")
		}
	}
	f.mark()
	blitBuf(f.current(), src, int(x), int(y), int(w), int(h))
	return null
}

// blitBuf draws src scaled to w by h at (x, y) in dst, inside dst's clip.
// Pixels are premultiplied on both sides, so source-over is s + d(1 - sa).
func blitBuf(dst, src *Buf, x, y, w, h int) {
	if w <= 0 || h <= 0 || src.W() == 0 || src.H() == 0 {
		return
	}
	area := image.Rect(x, y, x+w, y+h).Intersect(dst.clip)
	for dy := area.Min.Y; dy < area.Max.Y; dy++ {
		sy := (dy - y) * src.H() / h
		for dx := area.Min.X; dx < area.Max.X; dx++ {
			sx := (dx - x) * src.W() / w
			s := src.Img.Pix[sy*src.Img.Stride+sx*4:]
			if s[3] == 0 {
				continue
			}
			d := dst.Img.Pix[dy*dst.Img.Stride+dx*4:]
			if s[3] == 255 {
				d[0], d[1], d[2], d[3] = s[0], s[1], s[2], 255
				continue
			}
			inv := 255 - int(s[3])
			for i := 0; i < 4; i++ {
				d[i] = uint8(min(255, int(s[i])+(int(d[i])*inv+127)/255))
			}
		}
	}
}
