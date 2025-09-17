package golb_test

import (
	"testing"

	"github.com/nxtgo/golb"
)

type BenchmarkCase struct {
	name        string
	pattern     string
	fixture     string
	shouldMatch bool
}

var benchmarkTable = []BenchmarkCase{
	// complex pattern with character classes and wildcards
	{"ComplexCat_Match", "[a-z][!a-x]*cat*[h][!b]*eyes*", "my cat has very bright eyes", true},
	{"ComplexCat_NoMatch", "[a-z][!a-x]*cat*[h][!b]*eyes*", "my dog has very bright eyes", false},

	// url patterns, simplified from regex escaping
	{"GoogleURL_Match", "https://*.google.*", "https://account.google.com", true},
	{"GoogleURL_NoMatch", "https://*.google.*", "https://google.com", false},

	// multiple alternatives (simplified from regex or)
	{"MultiURL_Match", "{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}", "http://yahoo.com", true},
	{"MultiURL_NoMatch", "{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}", "http://google.com", false},

	// gobwas domain patterns
	{"Gobwas_Match", "{https://*.gobwas.com,http://exclude.gobwas.com}", "https://safe.gobwas.com", true},
	{"Gobwas_NoMatch", "{https://*.gobwas.com,http://exclude.gobwas.com}", "http://safe.gobwas.com", false},

	// simple prefix patterns
	{"Prefix_Match", "abc*", "abcdef", true},
	{"Prefix_NoMatch", "abc*", "af", false},

	// simple suffix patterns
	{"Suffix_Match", "*def", "abcdef", true},
	{"Suffix_NoMatch", "*def", "af", false},

	// prefix + suffix patterns
	{"PrefixSuffix_Match", "ab*ef", "abcdef", true},
	{"PrefixSuffix_NoMatch", "ab*ef", "af", false},

	// additional glob-specific patterns for completeness
	{"SuperWildcard_Match", "a**z", "a/b/c/d/z", true},
	{"SuperWildcard_NoMatch", "a**z", "a/b/c/d/y", false},
	{"SingleChar_Match", "a?c", "abc", true},
	{"SingleChar_NoMatch", "a?c", "ac", false},
	{"Range_Match", "[0-9][0-9][0-9]", "123", true},
	{"Range_NoMatch", "[0-9][0-9][0-9]", "abc", false},
	{"NegatedClass_Match", "[!abc]def", "xdef", true},
	{"NegatedClass_NoMatch", "[!abc]def", "adef", false},
}

func BenchmarkComplexCat_Match(b *testing.B) {
	runSingleBenchmark(b, "[a-z][!a-x]*cat*[h][!b]*eyes*", "my cat has very bright eyes")
}

func BenchmarkComplexCat_NoMatch(b *testing.B) {
	runSingleBenchmark(b, "[a-z][!a-x]*cat*[h][!b]*eyes*", "my dog has very bright eyes")
}

func BenchmarkGoogleURL_Match(b *testing.B) {
	runSingleBenchmark(b, "https://*.google.*", "https://account.google.com")
}

func BenchmarkGoogleURL_NoMatch(b *testing.B) {
	runSingleBenchmark(b, "https://*.google.*", "https://google.com")
}

func BenchmarkMultiURL_Match(b *testing.B) {
	runSingleBenchmark(b, "{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}", "http://yahoo.com")
}

func BenchmarkMultiURL_NoMatch(b *testing.B) {
	runSingleBenchmark(b, "{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}", "http://google.com")
}

func BenchmarkGobwas_Match(b *testing.B) {
	runSingleBenchmark(b, "{https://*.gobwas.com,http://exclude.gobwas.com}", "https://safe.gobwas.com")
}

func BenchmarkGobwas_NoMatch(b *testing.B) {
	runSingleBenchmark(b, "{https://*.gobwas.com,http://exclude.gobwas.com}", "http://safe.gobwas.com")
}

func BenchmarkPrefix_Match(b *testing.B) {
	runSingleBenchmark(b, "abc*", "abcdef")
}

func BenchmarkPrefix_NoMatch(b *testing.B) {
	runSingleBenchmark(b, "abc*", "af")
}

func BenchmarkSuffix_Match(b *testing.B) {
	runSingleBenchmark(b, "*def", "abcdef")
}

func BenchmarkSuffix_NoMatch(b *testing.B) {
	runSingleBenchmark(b, "*def", "af")
}

func BenchmarkPrefixSuffix_Match(b *testing.B) {
	runSingleBenchmark(b, "ab*ef", "abcdef")
}

func BenchmarkPrefixSuffix_NoMatch(b *testing.B) {
	runSingleBenchmark(b, "ab*ef", "af")
}

func BenchmarkSuperWildcard_Match(b *testing.B) {
	runSingleBenchmark(b, "a**z", "a/b/c/d/z")
}

func BenchmarkSingleChar_Match(b *testing.B) {
	runSingleBenchmark(b, "a?c", "abc")
}

func BenchmarkRange_Match(b *testing.B) {
	runSingleBenchmark(b, "[0-9][0-9][0-9]", "123")
}

func runSingleBenchmark(b *testing.B, pattern, fixture string) {
	g := golb.Compile(pattern)
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		g.Match(fixture)
	}
}

func BenchmarkAllPatterns(b *testing.B) {
	compiled := make([]*golb.Glob, len(benchmarkTable))
	for i, tc := range benchmarkTable {
		compiled[i] = golb.Compile(tc.pattern)
	}

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		for j, tc := range benchmarkTable {
			result := compiled[j].Match(tc.fixture)
			if result != tc.shouldMatch {
				b.Errorf("pattern %q with fixture %q: expected %v, got %v",
					tc.pattern, tc.fixture, tc.shouldMatch, result)
			}
		}
	}
}

func BenchmarkCompilation(b *testing.B) {
	patterns := []string{
		"[a-z][!a-x]*cat*[h][!b]*eyes*",
		"https://*.google.*",
		"{https://*.google.*,*yandex.*,*yahoo.*,*mail.ru}",
		"abc*",
		"*def",
		"ab*ef",
		"a**z",
		"a?c",
		"[0-9][0-9][0-9]",
		"[!abc]def",
	}

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		for _, pattern := range patterns {
			golb.Compile(pattern)
		}
	}
}

func BenchmarkConvenienceMatch(b *testing.B) {
	for i := 0; b.Loop(); i++ {
		for _, tc := range benchmarkTable {
			golb.Match(tc.pattern, tc.fixture)
		}
	}
}

func BenchmarkMemoryAlloc(b *testing.B) {
	pattern := "[a-z][!a-x]*cat*[h][!b]*eyes*"
	fixture := "my cat has very bright eyes"

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		g := golb.Compile(pattern)
		g.Match(fixture)
	}
}

func BenchmarkGenerateTable(b *testing.B) {
	b.Skip("This is for generating comparison table - run manually")

	for _, tc := range benchmarkTable {
		g := golb.Compile(tc.pattern)

		b.Run(tc.name, func(b *testing.B) {
			for i := 0; b.Loop(); i++ {
				g.Match(tc.fixture)
			}
		})
	}
}

func BenchmarkFileMatching(b *testing.B) {
	patterns := []*golb.Glob{
		golb.Compile("*.go"),
		golb.Compile("**/*.go"),
		golb.Compile("*.{js,ts,jsx,tsx}"),
		golb.Compile("test/**/*_test.go"),
		golb.Compile("src/**/components/*.{js,jsx}"),
		golb.Compile("node_modules/**"),
	}

	files := []string{
		"main.go",
		"src/components/Button.jsx",
		"test/unit/parser_test.go",
		"node_modules/react/index.js",
		"docs/README.md",
		"build/dist/app.js",
		"src/utils/helper.ts",
		"test/integration/api_test.go",
	}

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		for _, pattern := range patterns {
			for _, file := range files {
				pattern.Match(file)
			}
		}
	}
}
