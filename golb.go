package golb

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Glob struct {
	pattern string
	seps    []rune
}

func Compile(pattern string, separators ...rune) *Glob {
	return &Glob{
		pattern: pattern,
		seps:    separators,
	}
}

func Must(pattern string, separators ...rune) *Glob {
	return Compile(pattern, separators...)
}

func (g *Glob) Match(s string) bool {
	return g.match(g.pattern, s)
}

func (g *Glob) match(pattern, s string) bool {
	for len(pattern) > 0 {
		switch pattern[0] {
		case '*':
			if len(pattern) > 1 && pattern[1] == '*' {
				return g.matchSuper(pattern[2:], s)
			}
			return g.matchStar(pattern[1:], s)

		case '?':
			if len(s) == 0 {
				return false
			}
			_, size := utf8.DecodeRuneInString(s)
			if g.hasSeparator(s[:size]) {
				return false
			}
			pattern = pattern[1:]
			s = s[size:]

		case '[':
			end := strings.Index(pattern, "]")
			if end == -1 {
				if len(s) == 0 || s[0] != '[' {
					return false
				}
				pattern = pattern[1:]
				s = s[1:]
			} else {
				if !g.matchClass(pattern[1:end], s) {
					return false
				}
				pattern = pattern[end+1:]
				_, size := utf8.DecodeRuneInString(s)
				s = s[size:]
			}

		case '{':
			end := g.findClosingBrace(pattern)
			if end == -1 {
				if len(s) == 0 || s[0] != '{' {
					return false
				}
				pattern = pattern[1:]
				s = s[1:]
			} else {
				return g.matchAlternatives(pattern[1:end], pattern[end+1:], s)
			}

		case '\\':
			if len(pattern) == 1 {
				return len(s) == 1 && s[0] == '\\'
			}
			if len(s) == 0 || s[0] != pattern[1] {
				return false
			}
			pattern = pattern[2:]
			s = s[1:]

		default:
			if len(s) == 0 || s[0] != pattern[0] {
				return false
			}
			pattern = pattern[1:]
			s = s[1:]
		}
	}
	return len(s) == 0
}

func (g *Glob) matchStar(pattern, s string) bool {
	for i := 0; i <= len(s); i++ {
		if g.hasSeparator(s[:i]) {
			break
		}
		if g.match(pattern, s[i:]) {
			return true
		}
	}
	return false
}

func (g *Glob) matchSuper(pattern, s string) bool {
	for i := 0; i <= len(s); i++ {
		if g.match(pattern, s[i:]) {
			return true
		}
	}
	return false
}

func (g *Glob) matchClass(class, s string) bool {
	if len(s) == 0 {
		return false
	}

	r, _ := utf8.DecodeRuneInString(s)
	negate := len(class) > 0 && class[0] == '!'
	if negate {
		class = class[1:]
	}

	matched := false
	for i := 0; i < len(class); {
		if i+2 < len(class) && class[i+1] == '-' {
			start := rune(class[i])
			end := rune(class[i+2])
			if r >= start && r <= end {
				matched = true
				break
			}
			i += 3
		} else {
			if r == rune(class[i]) {
				matched = true
				break
			}
			i++
		}
	}

	return matched != negate
}

func (g *Glob) findClosingBrace(pattern string) int {
	depth := 0
	for i, c := range pattern {
		switch c {
		case '{':
			depth++
		case '}':
			if depth == 0 {
				return i
			}
			depth--
		case '\\':
			i++
		}
	}
	return -1
}

func (g *Glob) matchAlternatives(alternatives, rest, s string) bool {
	alts := g.splitAlternatives(alternatives)
	for _, alt := range alts {
		if g.match(alt+rest, s) {
			return true
		}
	}
	return false
}

func (g *Glob) splitAlternatives(s string) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	for i, c := range s {
		switch c {
		case '{':
			depth++
			current.WriteRune(c)
		case '}':
			depth--
			current.WriteRune(c)
		case ',':
			if depth == 0 {
				parts = append(parts, current.String())
				current.Reset()
			} else {
				current.WriteRune(c)
			}
		case '\\':
			current.WriteRune(c)
			if i+1 < len(s) {
				current.WriteByte(s[i+1])
			}
		default:
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

func (g *Glob) hasSeparator(s string) bool {
	if len(g.seps) == 0 {
		return false
	}
	for _, r := range s {
		return slices.Contains(g.seps, r)
	}
	return false
}

func QuoteMeta(s string) string {
	var result strings.Builder
	for _, c := range s {
		if isSpecial(byte(c)) {
			result.WriteByte('\\')
		}
		result.WriteRune(c)
	}
	return result.String()
}

func isSpecial(c byte) bool {
	switch c {
	case '*', '?', '[', ']', '{', '}', '\\':
		return true
	}
	return false
}

func Match(pattern, s string, separators ...rune) bool {
	g := Compile(pattern, separators...)
	return g.Match(s)
}
