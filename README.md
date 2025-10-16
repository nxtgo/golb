# golb

glob pattern handling for go.

## example

```go
package main

import (
	"fmt"
	"github.com/nxtgo/golb"
)

func main() {
	pattern := "**/{readme,license}.md"
	files := []string{
		"readme.md",
		"docs/readme.md",
		"license.md",
		"src/license.txt",
	}

	for _, f := range files {
		if golb.Match(pattern, f, '/') {
			fmt.Println("matched:", f)
		}
	}
}
```

```
matched: readme.md
matched: docs/readme.md
matched: license.md
```

# license

under CC0 1.0 (public domain) + ip waiver.
