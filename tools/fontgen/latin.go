package main

// Latin-1 glyphs for the pixel font. basicfont only has ASCII, so letters
// are composed (ASCII base + a small accent) and a few symbols are drawn by
// hand. The cell is 7 wide and 13 tall: capitals sit on rows 2-10, lower
// case on rows 5-10, so accents live above (rows 0-3) and cedillas below.

const (
	cellW = 7
	cellH = 13
)

// glyph is one cell of set pixels, indexed [y][x].
type glyph [cellH][cellW]bool

// pat builds a glyph from rows of '#' and '.', placed from row top down.
func pat(top int, rows ...string) glyph {
	var g glyph
	for i, row := range rows {
		if len(row) != cellW {
			panic("fontgen: bad pattern row " + row)
		}
		for x, c := range row {
			if c == '#' {
				g[top+i][x] = true
			}
		}
	}
	return g
}

// overlay merges b into a copy of a.
func overlay(a, b glyph) glyph {
	for y := range a {
		for x := range a[y] {
			a[y][x] = a[y][x] || b[y][x]
		}
	}
	return a
}

// touches reports whether any pixel of acc is next to (4-neighbour) or on a
// pixel of base.
func touches(base, acc glyph) bool {
	for y := range acc {
		for x := range acc[y] {
			if !acc[y][x] {
				continue
			}
			if base[y][x] {
				return true
			}
			if y > 0 && base[y-1][x] || y < cellH-1 && base[y+1][x] ||
				x > 0 && base[y][x-1] || x < cellW-1 && base[y][x+1] {
				return true
			}
		}
	}
	return false
}

// accent has the form used over lower case (rows 1-3, a blank row 4 above
// the x-height), over capitals (rows 0-1) and a flat one-row form for
// capitals whose top row would touch the two-row form.
type accent struct{ lower, upper, flat glyph }

var accents = map[string]accent{
	"acute": {
		pat(2, "....#..", "...#..."),
		pat(0, "....#..", "...#..."),
		pat(0, "....##."),
	},
	"grave": {
		pat(2, "..#....", "...#..."),
		pat(0, "..#....", "...#..."),
		pat(0, "..##..."),
	},
	"circ": {
		pat(2, "...#...", "..#.#.."),
		pat(0, "...#...", "..#.#.."),
		pat(0, "..###.."),
	},
	"tilde": {
		pat(2, "..##.#.", ".#..#.."),
		pat(0, "..##.#.", ".#..#.."),
		pat(0, ".#.##.."),
	},
	"diaer": {
		pat(3, "..#.#.."),
		pat(1, "..#.#.."),
		pat(0, "..#.#.."),
	},
	"ring": {
		pat(1, "...#...", "..#.#..", "...#..."),
		pat(0, "..##...", ".#..#.."),
		pat(0, "..##..."),
	},
}

// withAccent puts the named accent on base. isUpper picks the capital form.
func withAccent(base glyph, name string, isUpper bool) glyph {
	a := accents[name]
	if !isUpper {
		return overlay(base, a.lower)
	}
	if touches(base, a.upper) {
		return overlay(base, a.flat)
	}
	return overlay(base, a.upper)
}

// dotless removes the dot of i.
func dotless(g glyph) glyph {
	for x := range g[3] {
		g[3][x] = false
	}
	return g
}

// cedilla hangs a small hook under a letter.
func cedilla(g glyph) glyph {
	return overlay(g, pat(11, "...#...", "..#...."))
}

// rotate180 turns a glyph half a turn inside rows top..bottom and cols 0..5.
func rotate180(g glyph, top, bottom int) glyph {
	var out glyph
	for y := top; y <= bottom; y++ {
		for x := 0; x <= 5; x++ {
			out[top+bottom-y][5-x] = g[y][x]
		}
	}
	return out
}

// shiftRight moves a glyph one column right.
func shiftRight(g glyph) glyph {
	var out glyph
	for y := range g {
		for x := 0; x < cellW-1; x++ {
			out[y][x+1] = g[y][x]
		}
	}
	return out
}

func fallbackBox() glyph {
	return pat(2,
		".#####.",
		".#...#.",
		".#...#.",
		".#...#.",
		".#...#.",
		".#...#.",
		".#...#.",
		".#...#.",
		".#####.")
}

// latin returns every Latin-1 glyph (and the extra symbols) built from the
// ASCII glyphs.
func latin(a map[rune]glyph) map[rune]glyph {
	m := map[rune]glyph{}

	// Accented letters: rune -> base letter + accent.
	type spec struct {
		r    rune
		base rune
		acc  string
	}
	specs := []spec{
		{'À', 'A', "grave"}, {'Á', 'A', "acute"}, {'Â', 'A', "circ"}, {'Ã', 'A', "tilde"}, {'Ä', 'A', "diaer"}, {'Å', 'A', "ring"},
		{'È', 'E', "grave"}, {'É', 'E', "acute"}, {'Ê', 'E', "circ"}, {'Ë', 'E', "diaer"},
		{'Ì', 'I', "grave"}, {'Í', 'I', "acute"}, {'Î', 'I', "circ"}, {'Ï', 'I', "diaer"},
		{'Ñ', 'N', "tilde"},
		{'Ò', 'O', "grave"}, {'Ó', 'O', "acute"}, {'Ô', 'O', "circ"}, {'Õ', 'O', "tilde"}, {'Ö', 'O', "diaer"},
		{'Ù', 'U', "grave"}, {'Ú', 'U', "acute"}, {'Û', 'U', "circ"}, {'Ü', 'U', "diaer"},
		{'Ý', 'Y', "acute"},
		{'à', 'a', "grave"}, {'á', 'a', "acute"}, {'â', 'a', "circ"}, {'ã', 'a', "tilde"}, {'ä', 'a', "diaer"}, {'å', 'a', "ring"},
		{'è', 'e', "grave"}, {'é', 'e', "acute"}, {'ê', 'e', "circ"}, {'ë', 'e', "diaer"},
		{'ì', 'i', "grave"}, {'í', 'i', "acute"}, {'î', 'i', "circ"}, {'ï', 'i', "diaer"},
		{'ñ', 'n', "tilde"},
		{'ò', 'o', "grave"}, {'ó', 'o', "acute"}, {'ô', 'o', "circ"}, {'õ', 'o', "tilde"}, {'ö', 'o', "diaer"},
		{'ù', 'u', "grave"}, {'ú', 'u', "acute"}, {'û', 'u', "circ"}, {'ü', 'u', "diaer"},
		{'ý', 'y', "acute"}, {'ÿ', 'y', "diaer"},
	}
	for _, s := range specs {
		base := a[s.base]
		if s.base == 'i' {
			base = dotless(base)
		}
		m[s.r] = withAccent(base, s.acc, s.base >= 'A' && s.base <= 'Z')
	}
	m['Ç'] = cedilla(a['C'])
	m['ç'] = cedilla(a['c'])

	// Inverted marks.
	m['¿'] = rotate180(a['?'], 2, 10)
	m['¡'] = rotate180(a['!'], 2, 10)

	// Drawn by hand.
	m['ß'] = pat(2,
		".###...",
		"#...#..",
		"#...#..",
		"#..#...",
		"#...#..",
		"#....#.",
		"#....#.",
		"#....#.",
		"#.###..")
	m['Æ'] = pat(2,
		"..#####",
		".##....",
		".#.#...",
		"#..#...",
		"#####..",
		"#..#...",
		"#..#...",
		"#..#...",
		"#..####")
	m['æ'] = pat(5,
		".##.##.",
		"...#..#",
		".######",
		"#..#...",
		"#..#..#",
		".##.##.")
	m['Ø'] = pat(2,
		".####.#",
		"#....#.",
		"#...##.",
		"#..#.#.",
		"#.#..#.",
		"##...#.",
		"#....#.",
		"#....#.",
		".####..")
	m['ø'] = pat(5,
		".####.#",
		"#...##.",
		"#..#.#.",
		"#.#..#.",
		"##...#.",
		".####..")
	m['Ð'] = overlay(shiftRight(a['D']), pat(6, "#......"))
	m['ð'] = pat(2,
		"..##...",
		"...##..",
		"..#.#..",
		".####..",
		"#....#.",
		"#....#.",
		"#....#.",
		"#....#.",
		".####..")
	m['Þ'] = pat(2,
		"#......",
		"####...",
		"#...#..",
		"#....#.",
		"#...#..",
		"####...",
		"#......",
		"#......",
		"#......")
	m['þ'] = pat(3,
		"#......",
		"#......",
		"#.###..",
		"##...#.",
		"#....#.",
		"##...#.",
		"#.###..",
		"#......",
		"#......")
	m['€'] = pat(2,
		"..###..",
		".#...#.",
		"###....",
		"#......",
		"####...",
		"#......",
		"###....",
		".#...#.",
		"..###..")
	m['«'] = pat(5,
		"..#..#.",
		".#..#..",
		"#..#...",
		".#..#..",
		"..#..#.")
	m['»'] = pat(5,
		".#..#..",
		"..#..#.",
		"...#..#",
		"..#..#.",
		".#..#..")
	m['°'] = pat(2,
		".##....",
		"#..#...",
		"#..#...",
		".##....")
	m['·'] = pat(6, "...#...")
	m['—'] = pat(6, "#######")
	m['…'] = pat(10, "#.#.#.#")
	m['“'] = pat(2,
		".#..#..",
		"#..#...",
		"##.##..",
		"##.##..")
	m['”'] = pat(2,
		"##.##..",
		"##.##..",
		".#..#..",
		"#..#...")
	m['‘'] = pat(2,
		"..#....",
		".#.....",
		"##.....",
		"##.....")
	m['’'] = pat(2,
		"##.....",
		"##.....",
		".#.....",
		"#......")

	// Common symbols.
	m['×'] = pat(5,
		"#....#.",
		".#..#..",
		"..##...",
		".#..#..",
		"#....#.")
	m['÷'] = pat(4,
		"...#...",
		".......",
		"#######",
		".......",
		"...#...")
	m['±'] = pat(3,
		"...#...",
		"...#...",
		"#######",
		"...#...",
		"...#...",
		".......",
		"#######")
	m['¢'] = pat(3,
		"...#...",
		".####..",
		"#..#..#",
		"#..#...",
		"#..#..#",
		".####..",
		"...#...")
	m['£'] = pat(2,
		"..###..",
		".#...#.",
		".#.....",
		"####...",
		".#.....",
		".#.....",
		"#.....#",
		"#######")
	m['¥'] = pat(2,
		"#....#.",
		".#..#..",
		"..##...",
		"######.",
		"..##...",
		"######.",
		"..##...",
		"..##...")
	m['©'] = pat(2,
		"..###..",
		".#...#.",
		"#.###.#",
		"#.#...#",
		"#.#...#",
		"#.###.#",
		".#...#.",
		"..###..")
	m['ª'] = pat(1,
		".###...",
		"....#..",
		".####..",
		".#..#..",
		".####..",
		".......",
		".####..")
	m['º'] = pat(1,
		".###...",
		".#..#..",
		".#..#..",
		".###...",
		".......",
		".####..")
	m['¨'] = pat(2, "..#.#..")
	m['´'] = pat(2, "....#..", "...#...")
	return m
}
