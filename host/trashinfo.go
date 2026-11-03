package host

import (
	"strings"
	"time"
)

// The freedesktop Trash spec (Linux) keeps a small .trashinfo next to every
// trashed file so a file manager can put it back. These two helpers are plain
// text work, so they live here and are tested on every system.

// trashInfoText is the content of Trash/info/<name>.trashinfo.
func trashInfoText(origPath string, when time.Time) string {
	return "[Trash Info]\nPath=" + percentEncodePath(origPath) +
		"\nDeletionDate=" + when.Format("2006-01-02T15:04:05") + "\n"
}

// percentEncodePath escapes a path the way the spec wants: everything but
// letters, digits, "-_.~" and "/" becomes %XX.
func percentEncodePath(p string) string {
	const hex = "0123456789ABCDEF"
	var b strings.Builder
	for i := 0; i < len(p); i++ {
		c := p[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9',
			c == '-', c == '_', c == '.', c == '~', c == '/':
			b.WriteByte(c)
		default:
			b.WriteByte('%')
			b.WriteByte(hex[c>>4])
			b.WriteByte(hex[c&15])
		}
	}
	return b.String()
}
