package main

import (
	"fmt"
	"github.com/nxtgo/golb"
)

func main() {
	pattern := golb.Compile("*.go")
	files := []string{"main.go", "test.go", "readme.txt", "config.json"}

	fmt.Println("go files:")
	for _, file := range files {
		if pattern.Match(file) {
			fmt.Printf("  %s v/\n", file)
		}
	}

	examples := []struct {
		pattern string
		text    string
		desc    string
	}{
		{"a*c", "abc", "simple wildcard"},
		{"a?c", "abc", "single character"},
		{"a**c", "a/b/c", "super wildcard (crosses directories)"},
		{"[abc]", "b", "character class"},
		{"[!abc]", "x", "negated character class"},
		{"[a-z]", "m", "range"},
		{"{go,js,py}", "py", "alternatives"},
		{"*.{go,js}", "main.go", "combined patterns"},
		{"**/*.go", "src/main.go", "recursive glob"},
	}

	fmt.Println("\npattern examples:")
	for _, ex := range examples {
		result := golb.Match(ex.pattern, ex.text)
		status := "X"
		if result {
			status = "v/"
		}
		fmt.Printf("  %-15s %-12s %s %s\n", ex.pattern, ex.text, status, ex.desc)
	}

	pathPattern := golb.Compile("*.go", '/')
	fmt.Println("\nwith path separators:")
	paths := []string{"main.go", "dir/main.go", "file.txt"}
	for _, path := range paths {
		if pathPattern.Match(path) {
			fmt.Printf("  %s matches *.go\n", path)
		}
	}

	escaped := golb.QuoteMeta("file[1].txt")
	fmt.Printf("\nescaped pattern: %s\n", escaped)
	fmt.Printf("matches 'file[1].txt': %v\n", golb.Match(escaped, "file[1].txt"))
}
