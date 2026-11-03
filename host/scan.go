package host

import (
	"regexp"
	"strings"

	"lizardscript/lexer"
	"lizardscript/token"
)

// The source scanner: the editor's view of a program as lines, statements and
// argument spans. It reads the lexer's tokens (so positions are the
// language's own), not the AST, which has no end positions.

// Arg is one argument of a call, or one element of an array or point.
type Arg struct {
	Text    string
	Line    int
	EndLine int
	Start   int // column of the first character, 1-based
	End     int // column just after the last character, on EndLine
	Kind    string
	Items   []Arg
}

// Stmt is one statement.
type Stmt struct {
	Line, EndLine, Column, Depth int
	Kind                         string // layer call let function assign block other
	Name, Callee                 string
	Args                         []Arg
	Children                     []*Stmt
	Params                       []string
	Doc                          *Doc
	Label                        string // trailing comment of the first line, without flags
	Flags                        []string
}

// Doc is the -~> comment block above a function.
type Doc struct {
	Summary string
	Params  map[string]string
	Example string
	HasEx   bool
}

// LineInfo is one line of the file.
type LineInfo struct {
	Line    int
	Indent  int
	Text    string
	Comment string
	Name    string
	Flags   []string
}

// Scan is the result of scanning a source.
type Scan struct {
	Lines  []LineInfo
	Stmts  []*Stmt
	Errors []SrcError
	runes  [][]rune
}

// SrcError is an error with its position.
type SrcError struct {
	Line, Column int
	Message      string
}

var flagRe = regexp.MustCompile(`\[(\w+)\]`)

// ScanSource scans code.
func ScanSource(code string) *Scan {
	sc := &Scan{}
	for i, text := range strings.Split(strings.ReplaceAll(code, "\r\n", "\n"), "\n") {
		sc.runes = append(sc.runes, []rune(text))
		sc.Lines = append(sc.Lines, lineInfo(i+1, text))
	}

	var toks []token.Token
	l := lexer.New(code)
	for {
		t := l.NextToken()
		if t.Type == token.EOF {
			break
		}
		if t.Type == token.ILLEGAL {
			sc.Errors = append(sc.Errors, SrcError{t.Line, t.Column, "illegal token " + t.Literal})
			continue
		}
		toks = append(toks, t)
	}
	p := &scanner{sc: sc, toks: toks}
	sc.Stmts = p.block(0, -1)
	return sc
}

// lineInfo finds the indent and the -~> comment of a line, skipping strings.
func lineInfo(n int, text string) LineInfo {
	li := LineInfo{Line: n, Text: text, Flags: []string{}}
	for _, r := range text {
		if r != ' ' && r != '\t' {
			break
		}
		li.Indent++
	}
	inStr := false
	rs := []rune(text)
	for i := 0; i < len(rs); i++ {
		switch {
		case inStr && rs[i] == '\\':
			i++
		case rs[i] == '"':
			inStr = !inStr
		case !inStr && i+2 < len(rs) && string(rs[i:i+3]) == "-~>":
			li.Comment = strings.TrimSpace(string(rs[i+3:]))
		}
		if li.Comment != "" {
			break
		}
	}
	for _, m := range flagRe.FindAllStringSubmatch(li.Comment, -1) {
		if m[1] == "hidden" || m[1] == "locked" {
			li.Flags = append(li.Flags, m[1])
		}
	}
	name := flagRe.ReplaceAllStringFunc(li.Comment, func(m string) string {
		if m == "[hidden]" || m == "[locked]" {
			return ""
		}
		return m
	})
	li.Name = strings.TrimSpace(name)
	return li
}

type scanner struct {
	sc   *Scan
	toks []token.Token
}

func isOpen(t token.TokenType) bool {
	return t == token.LPAREN || t == token.LBRACKET || t == token.LBRACE
}
func isClose(t token.TokenType) bool {
	return t == token.RPAREN || t == token.RBRACKET || t == token.RBRACE
}

// continues reports whether a statement goes on after token i, which ends a
// line at bracket depth 0.
func (p *scanner) continues(i int) bool {
	if i+1 >= len(p.toks) {
		return false
	}
	switch p.toks[i].Type {
	case token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.DOUBLE_SLASH, token.POW,
		token.MODULE, token.ASSIGN, token.EQ, token.NOT_EQ, token.EXACT, token.NOT_EXACT,
		token.LT, token.GT, token.LTE, token.GTE, token.AND, token.OR, token.COMMA,
		token.DOT, token.COLON, token.BANG:
		return true
	}
	switch p.toks[i+1].Type {
	case token.ELSE, token.ELIF, token.DOT, token.AND, token.OR:
		return true
	}
	return false
}

// extent returns the index after the last token of the statement that starts
// at i.
func (p *scanner) extent(i int) int {
	depth := 0
	for j := i; j < len(p.toks); j++ {
		t := p.toks[j]
		switch {
		case isOpen(t.Type):
			depth++
		case isClose(t.Type):
			if depth == 0 { // a } that ends the enclosing block
				return j
			}
			depth--
		case t.Type == token.SEMICOLON && depth == 0:
			return j + 1
		}
		if depth == 0 && j+1 < len(p.toks) && p.toks[j+1].Line > t.Line && !p.continues(j) {
			// A { opened on this line and closed later keeps depth > 0.
			return j + 1
		}
	}
	return len(p.toks)
}

// block scans statements from index i. With layerCol >= 0 it stops at the
// first statement whose column is not deeper (a layer body). It also stops
// at a closing bracket. It returns the statements; p.pos has the next index.
func (p *scanner) block(i, layerCol int) []*Stmt {
	return p.blockAt(&i, layerCol, 0)
}

func (p *scanner) blockAt(pos *int, layerCol, depth int) []*Stmt {
	var out []*Stmt
	for *pos < len(p.toks) {
		t := p.toks[*pos]
		if isClose(t.Type) {
			return out
		}
		if t.Type == token.SEMICOLON {
			*pos++
			continue
		}
		if layerCol >= 0 && t.Column <= layerCol {
			return out
		}
		out = append(out, p.statement(pos, depth))
	}
	return out
}

func (p *scanner) statement(pos *int, depth int) *Stmt {
	i := *pos
	t := p.toks[i]

	// layer name
	if t.Type == token.IDENT && t.Literal == "layer" && i+1 < len(p.toks) {
		n := p.toks[i+1]
		endsLine := i+2 >= len(p.toks) || p.toks[i+2].Line > t.Line || isClose(p.toks[i+2].Type)
		if (n.Type == token.IDENT || n.Type == token.STRING) && n.Line == t.Line && endsLine {
			s := &Stmt{Line: t.Line, EndLine: t.Line, Column: t.Column, Depth: depth, Kind: "layer", Name: n.Literal}
			*pos = i + 2
			s.Children = p.blockAt(pos, t.Column, depth+1)
			if len(s.Children) > 0 {
				s.EndLine = s.Children[len(s.Children)-1].EndLine
			}
			p.decorate(s)
			return s
		}
	}

	end := p.extent(i)
	if end <= i {
		end = i + 1
	}
	s := &Stmt{Line: t.Line, EndLine: p.toks[end-1].Line, Column: t.Column, Depth: depth, Kind: "other"}
	toks := p.toks[i:end]
	*pos = end

	switch {
	case t.Type == token.LET && len(toks) > 1:
		s.Kind, s.Name = "let", toks[1].Literal
	case t.Type == token.FUNCTION && len(toks) > 2 && toks[1].Type == token.IDENT:
		s.Kind, s.Name = "function", toks[1].Literal
		s.Params = params(toks[2:])
	case t.Type == token.IF || t.Type == token.WHILE || t.Type == token.FOR:
		s.Kind, s.Name = "block", strings.ToLower(string(t.Type))
	case t.Type == token.IDENT:
		p.callOrAssign(s, toks)
	}
	if s.Kind == "block" || s.Kind == "function" || s.Kind == "let" {
		s.Children = p.bodies(i, end, depth+1)
	}
	p.decorate(s)
	return s
}

// params reads the names between the parentheses of a function.
func params(toks []token.Token) []string {
	names := []string{}
	if len(toks) == 0 || toks[0].Type != token.LPAREN {
		return names
	}
	for _, t := range toks[1:] {
		if t.Type == token.RPAREN {
			break
		}
		if t.Type == token.IDENT {
			names = append(names, t.Literal)
		}
	}
	return names
}

// callOrAssign fills in a statement that starts with a name: f(...), a.b(...)
// or x = ...
func (p *scanner) callOrAssign(s *Stmt, toks []token.Token) {
	j := 1
	name := toks[0].Literal
	for j+1 < len(toks) && toks[j].Type == token.DOT && toks[j+1].Type == token.IDENT {
		name += "." + toks[j+1].Literal
		j += 2
	}
	if j < len(toks) && toks[j].Type == token.LPAREN {
		s.Kind, s.Name, s.Callee = "call", name, name
		if close := matching(toks, j); close > 0 {
			s.Args = p.args(toks, j, close)
		}
		return
	}
	if j < len(toks) && (toks[j].Type == token.ASSIGN || toks[j].Type == token.LBRACKET) {
		// x = ... or x[i] = ...
		depth := 0
		for k := j; k < len(toks); k++ {
			if isOpen(toks[k].Type) {
				depth++
			} else if isClose(toks[k].Type) {
				depth--
			} else if depth == 0 && toks[k].Type == token.ASSIGN {
				s.Kind, s.Name = "assign", name
				return
			}
		}
	}
}

// matching returns the index of the bracket that closes toks[open], or -1.
func matching(toks []token.Token, open int) int {
	depth := 0
	for k := open; k < len(toks); k++ {
		if isOpen(toks[k].Type) {
			depth++
		} else if isClose(toks[k].Type) {
			depth--
			if depth == 0 {
				return k
			}
		}
	}
	return -1
}

// bodies scans the statements inside the { } bodies of toks[from:to].
func (p *scanner) bodies(from, to, depth int) []*Stmt {
	var out []*Stmt
	parenDepth := 0
	for k := from; k < to; k++ {
		switch p.toks[k].Type {
		case token.LPAREN, token.LBRACKET:
			parenDepth++
		case token.RPAREN, token.RBRACKET:
			parenDepth--
		case token.LBRACE:
			if parenDepth > 0 && !p.looksLikeBody(k) {
				continue
			}
			pos := k + 1
			out = append(out, p.blockAt(&pos, -1, depth)...)
			// skip to the closing brace
			k = pos
		}
	}
	return out
}

// looksLikeBody is false for the { of a hash literal inside parentheses. A
// function body inside a call (a callback) is rare in drawings; treat every
// brace inside parentheses as data.
func (p *scanner) looksLikeBody(k int) bool { return false }

// args splits toks[open+1:close] at commas of depth 0.
func (p *scanner) args(toks []token.Token, open, close int) []Arg {
	out := []Arg{}
	start := open + 1
	depth := 0
	for k := open + 1; k <= close; k++ {
		t := toks[k]
		if k < close && isOpen(t.Type) {
			depth++
			continue
		}
		if k < close && isClose(t.Type) {
			depth--
			continue
		}
		if (k == close || (t.Type == token.COMMA && depth == 0)) && k > start {
			out = append(out, p.arg(toks, start, k, t))
			start = k + 1
		} else if t.Type == token.COMMA && depth == 0 {
			start = k + 1
		}
	}
	return out
}

// arg builds the Arg for toks[from:to], which ends just before terminator.
func (p *scanner) arg(toks []token.Token, from, to int, terminator token.Token) Arg {
	first := toks[from]
	a := Arg{Line: first.Line, Start: first.Column}
	// The argument ends right before the terminator, minus the spaces.
	a.EndLine, a.End = terminator.Line, terminator.Column
	a.Text, a.End = p.trimmed(a.Line, a.Start, a.EndLine, a.End)
	if a.End <= a.Start && a.EndLine == a.Line {
		a.End = a.Start
	}

	sub := toks[from:to]
	switch {
	case len(sub) == 1 && (sub[0].Type == token.INT || sub[0].Type == token.FLOAT):
		a.Kind = "number"
	case len(sub) == 2 && sub[0].Type == token.MINUS && (sub[1].Type == token.INT || sub[1].Type == token.FLOAT):
		a.Kind = "number"
	case len(sub) == 1 && sub[0].Type == token.COLOR:
		a.Kind = "color"
	case len(sub) == 1 && sub[0].Type == token.STRING:
		a.Kind = "string"
	case sub[0].Type == token.LBRACKET && matching(sub, 0) == len(sub)-1:
		a.Kind = "array"
		a.Items = p.args(sub, 0, len(sub)-1)
	case sub[0].Type == token.LPAREN && matching(sub, 0) == len(sub)-1 && hasComma(sub):
		a.Kind = "point"
		a.Items = p.args(sub, 0, len(sub)-1)
	default:
		a.Kind = "expr"
	}
	return a
}

func hasComma(sub []token.Token) bool {
	depth := 0
	for _, t := range sub[1 : len(sub)-1] {
		if isOpen(t.Type) {
			depth++
		} else if isClose(t.Type) {
			depth--
		} else if t.Type == token.COMMA && depth == 0 {
			return true
		}
	}
	return false
}

// trimmed returns the source between two positions without trailing blanks,
// and the column where that text really ends.
func (p *scanner) trimmed(l1, c1, l2, c2 int) (string, int) {
	var parts []string
	for l := l1; l <= l2 && l <= len(p.sc.runes); l++ {
		rs := p.sc.runes[l-1]
		from, to := 0, len(rs)
		if l == l1 {
			from = c1 - 1
		}
		if l == l2 {
			to = min(c2-1, len(rs))
		}
		if from > to {
			from = to
		}
		parts = append(parts, string(rs[from:to]))
	}
	text := strings.Join(parts, "\n")
	cut := strings.TrimRight(text, " \t\r\n")
	removed := len([]rune(text)) - len([]rune(cut))
	return cut, c2 - removed
}

// decorate adds the trailing-comment label and the doc.
func (p *scanner) decorate(s *Stmt) {
	if s.Line-1 < len(p.sc.Lines) {
		li := p.sc.Lines[s.Line-1]
		s.Label, s.Flags = li.Name, li.Flags
		// A multi-line statement carries its comment on the last line.
		if s.Kind != "layer" && s.EndLine != s.Line && s.EndLine-1 < len(p.sc.Lines) {
			if end := p.sc.Lines[s.EndLine-1]; s.Label == "" && s.Kind != "block" && s.Kind != "function" {
				s.Label, s.Flags = end.Name, end.Flags
			}
		}
	}
	if s.Kind == "function" {
		s.Doc = p.doc(s)
	}
}

var (
	paramLineRe = regexp.MustCompile(`^(\w+):\s*(.*)$`)
	exampleRe   = regexp.MustCompile(`^example:\s*(.*)$`)
)

// doc collects the -~> lines right above a function.
func (p *scanner) doc(s *Stmt) *Doc {
	var lines []string
	for l := s.Line - 2; l >= 0; l-- {
		t := strings.TrimSpace(p.sc.Lines[l].Text)
		if !strings.HasPrefix(t, "-~>") {
			break
		}
		lines = append([]string{strings.TrimSpace(t[3:])}, lines...)
	}
	d := &Doc{Params: map[string]string{}}
	known := map[string]bool{}
	for _, n := range s.Params {
		known[n] = true
	}
	var summary []string
	summaryDone := false
	last := "" // "param:name" or "example"
	for _, ln := range lines {
		if m := exampleRe.FindStringSubmatch(ln); m != nil {
			d.Example, d.HasEx, last, summaryDone = m[1], true, "example", true
			continue
		}
		if m := paramLineRe.FindStringSubmatch(ln); m != nil && known[m[1]] {
			d.Params[m[1]] = m[2]
			last, summaryDone = "param:"+m[1], true
			continue
		}
		switch {
		case ln == "":
			if len(summary) > 0 {
				summaryDone = true
			}
			last = ""
		case !summaryDone:
			summary = append(summary, ln)
		case strings.HasPrefix(last, "param:"):
			d.Params[last[6:]] += " " + ln
		case last == "example":
			d.Example += "\n" + ln
		}
	}
	d.Summary = strings.Join(summary, " ")
	return d
}
