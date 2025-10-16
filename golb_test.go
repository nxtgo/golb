package golb_test

import (
	"testing"

	"github.com/nxtgo/golb"
)

func TestBasicMatching(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"abc", "abc", true},
		{"abc", "abd", false},
		{"abc*", "abcdef", true},
		{"abc*", "af", false},
		{"*def", "abcdef", true},
		{"*def", "af", false},
		{"ab*ef", "abcdef", true},
		{"ab*ef", "af", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.input, func(t *testing.T) {
			g := golb.Compile(tt.pattern)
			got := g.Match(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
			}
		})
	}
}

func TestWildcards(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"a?c", "abc", true},
		{"a?c", "ac", false},
		{"a?c", "abcd", false},
		{"a**z", "a/b/c/d/z", true},
		{"a**z", "a/b/c/d/y", false},
		{"a**z", "az", true},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.input, func(t *testing.T) {
			g := golb.Compile(tt.pattern)
			got := g.Match(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
			}
		})
	}
}

func TestCharacterClasses(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"[abc]", "a", true},
		{"[abc]", "b", true},
		{"[abc]", "d", false},
		{"[0-9]", "5", true},
		{"[0-9]", "a", false},
		{"[0-9][0-9][0-9]", "123", true},
		{"[0-9][0-9][0-9]", "abc", false},
		{"[a-z]def", "xdef", true},
		{"[a-z]def", "1def", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.input, func(t *testing.T) {
			g := golb.Compile(tt.pattern)
			got := g.Match(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
			}
		})
	}
}

func TestNegatedClasses(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"[!abc]def", "xdef", true},
		{"[!abc]def", "adef", false},
		{"[^abc]def", "xdef", true},
		{"[^abc]def", "adef", false},
		{"[!0-9]", "a", true},
		{"[!0-9]", "5", false},
		{"[^0-9]", "a", true},
		{"[^0-9]", "5", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.input, func(t *testing.T) {
			g := golb.Compile(tt.pattern)
			got := g.Match(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
			}
		})
	}
}

func TestAlternatives(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"{a,b,c}", "a", true},
		{"{a,b,c}", "b", true},
		{"{a,b,c}", "d", false},
		{"{ab,cd}", "ab", true},
		{"{ab,cd}", "cd", true},
		{"{ab,cd}", "ac", false},
		{"{*yahoo.*}", "http://yahoo.com", true},
		{"{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}", "http://yahoo.com", true},
		{"{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}", "https://maps.google.com", true},
		{"{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}", "http://google.com", false},
		{"{https://*.gobwas.com,http://exclude.gobwas.com}", "https://safe.gobwas.com", true},
		{"{https://*.gobwas.com,http://exclude.gobwas.com}", "http://exclude.gobwas.com", true},
		{"{https://*.gobwas.com,http://exclude.gobwas.com}", "http://safe.gobwas.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.input, func(t *testing.T) {
			g := golb.Compile(tt.pattern)
			got := g.Match(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
			}
		})
	}
}

func TestComplexPatterns(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"[a-z][!a-x]*cat*[h][!b]*eyes*", "my cat has very bright eyes", true},
		{"[a-z][!a-x]*cat*[h][!b]*eyes*", "my dog has very bright eyes", false},
		{"https://*.google.*", "https://account.google.com", true},
		{"https://*.google.*", "https://google.com", false},
		{"*.{js,ts,jsx,tsx}", "app.js", true},
		{"*.{js,ts,jsx,tsx}", "app.ts", true},
		{"*.{js,ts,jsx,tsx}", "app.jsx", true},
		{"*.{js,ts,jsx,tsx}", "app.tsx", true},
		{"*.{js,ts,jsx,tsx}", "app.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.input, func(t *testing.T) {
			g := golb.Compile(tt.pattern)
			got := g.Match(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
			}
		})
	}
}

func TestFileMatching(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"*.go", "main.go", true},
		{"*.go", "main.js", false},
		{"**/*.go", "src/main.go", true},
		{"**/*.go", "src/utils/helper.go", true},
		{"**/*.go", "main.js", false},
		{"test/**/*_test.go", "test/unit/parser_test.go", true},
		{"test/**/*_test.go", "test/integration/api_test.go", true},
		{"test/**/*_test.go", "src/main.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.input, func(t *testing.T) {
			g := golb.Compile(tt.pattern)
			got := g.Match(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
			}
		})
	}
}

func TestEscaping(t *testing.T) {
	tests := []struct {
		pattern string
		input   string
		want    bool
	}{
		{"\\*", "*", true},
		{"\\*", "a", false},
		{"\\?", "?", true},
		{"\\?", "a", false},
		{"a\\*b", "a*b", true},
		{"a\\*b", "aXb", false},
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.input, func(t *testing.T) {
			g := golb.Compile(tt.pattern)
			got := g.Match(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.input, got, tt.want)
			}
		})
	}
}

func TestConvenienceFunction(t *testing.T) {
	if !golb.Match("*.go", "main.go") {
		t.Error("Match(*.go, main.go) should be true")
	}
	if golb.Match("*.go", "main.js") {
		t.Error("Match(*.go, main.js) should be false")
	}
}

func TestQuoteMeta(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"abc", "abc"},
		{"a*b", "a\\*b"},
		{"a?b", "a\\?b"},
		{"a[b", "a\\[b"},
		{"a]b", "a\\]b"},
		{"a{b", "a\\{b"},
		{"a}b", "a\\}b"},
		{"a\\b", "a\\\\b"},
		{"*?[]{}\\", "\\*\\?\\[\\]\\{\\}\\\\"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := golb.QuoteMeta(tt.input)
			if got != tt.want {
				t.Errorf("QuoteMeta(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func BenchmarkSimplePrefix(b *testing.B) {
	g := golb.Compile("abc*")
	for i := 0; i < b.N; i++ {
		g.Match("abcdef")
	}
}

func BenchmarkSimpleSuffix(b *testing.B) {
	g := golb.Compile("*def")
	for i := 0; i < b.N; i++ {
		g.Match("abcdef")
	}
}

func BenchmarkComplexPattern(b *testing.B) {
	g := golb.Compile("[a-z][!a-x]*cat*[h][!b]*eyes*")
	for i := 0; i < b.N; i++ {
		g.Match("my cat has very bright eyes")
	}
}

func BenchmarkAlternatives(b *testing.B) {
	g := golb.Compile("{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}")
	for i := 0; i < b.N; i++ {
		g.Match("http://yahoo.com")
	}
}

func BenchmarkCompilation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		golb.Compile("[a-z][!a-x]*cat*[h][!b]*eyes*")
	}
}
