package tml

import (
	"fmt"
	"sort"
	"strings"
)

// Env is the named view of the top-level frame the tests and the
// examples read.
type Env map[string]Value

// String renders the bindings one per line, in name order — the
// stable format the tests and the examples check against.
func (env Env) String() string {
	names := make([]string, 0, len(env))
	for name := range env {
		names = append(names, name)
	}

	sort.Strings(names)

	lines := make([]string, 0, len(env))
	for _, name := range names {
		lines = append(lines, fmt.Sprintf("%s = %s", name, env[name]))
	}

	return strings.Join(lines, "\n")
}
