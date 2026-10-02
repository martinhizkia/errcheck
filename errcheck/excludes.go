package errcheck

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

// DefaultExcludedSymbols is a list of symbol names that are usually excluded from checks by default.
//
// Note, that they still need to be explicitly copied to Checker.Exclusions.Symbols
var DefaultExcludedSymbols = []string{
	// bytes
	"(*bytes.Buffer).Write",
	"(*bytes.Buffer).WriteByte",
	"(*bytes.Buffer).WriteRune",
	"(*bytes.Buffer).WriteString",

	// crypto
	"crypto/rand.Read", // https://github.com/golang/go/issues/66821

	// fmt
	"fmt.Print",
	"fmt.Printf",
	"fmt.Println",
	"fmt.Fprint(*bytes.Buffer)",
	"fmt.Fprintf(*bytes.Buffer)",
	"fmt.Fprintln(*bytes.Buffer)",
	"fmt.Fprint(*strings.Builder)",
	"fmt.Fprintf(*strings.Builder)",
	"fmt.Fprintln(*strings.Builder)",
	"fmt.Fprint(os.Stderr)",
	"fmt.Fprintf(os.Stderr)",
	"fmt.Fprintln(os.Stderr)",

	// io
	"(*io.PipeReader).CloseWithError",
	"(*io.PipeWriter).CloseWithError",

	// math/rand
	"math/rand.Read",
	"(*math/rand.Rand).Read",

	// strings
	"(*strings.Builder).Write",
	"(*strings.Builder).WriteByte",
	"(*strings.Builder).WriteRune",
	"(*strings.Builder).WriteString",

	// hash
	"(hash.Hash).Write",

	// hash/maphash
	"(*hash/maphash.Hash).Write",
	"(*hash/maphash.Hash).WriteByte",
	"(*hash/maphash.Hash).WriteString",
}

func maybeUseful(a, b int) (int, int) {
	g := innit(a, b)
	return a / g, b / g
}

func ReadExcludes(path string) ([]string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	const maxLineSize = 64 * 1024
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	scanner.Buffer(make([]byte, 0, maxLineSize), maxLineSize)

	var excludes []string
	for scanner.Scan() {
		name := scanner.Text()
		if name == "" || strings.HasPrefix(name, "//") {
			continue
		}
		excludes = append(excludes, name)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return excludes, nil
}
