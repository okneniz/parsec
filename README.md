# Parsec

![CI](https://github.com/okneniz/parsec/actions/workflows/ci.yml/badge.svg) ![Downloads](https://img.shields.io/github/downloads/okneniz/parsec/total) ![Contributors](https://img.shields.io/github/contributors/okneniz/parsec?color=dark-green) ![Forks](https://img.shields.io/github/forks/okneniz/parsec?style=social) ![Stargazers](https://img.shields.io/github/stars/okneniz/parsec?style=social) ![Issues](https://img.shields.io/github/issues/okneniz/parsec) ![License](https://img.shields.io/github/license/okneniz/parsec) 

Golang parser combinator library inspired by [haskell parsec](https://hackage.haskell.org/package/parsec).

> But what is parser combinator?

> In the parse combinatorial framework, a "parser" is a function that takes some semistructured input and produces some structured output, and "combinator" is a function that allows combining / composing things. So "parser combinators" is a way of expressing a system where you write a lot of small parsing functions and compose then together.


## Installation

Requires Go 1.26 or newer.

```bash
go get github.com/okneniz/parsec
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/okneniz/parsec/strings"
)

func main() {
	// parse an unsigned integer, allowing surrounding whitespaces
	parser := strings.Padded(
		strings.Try(strings.Space("whitespace")),
		strings.Unsigned[int](),
	)

	result, err := strings.ParseString(" 42 ", parser)
	fmt.Println(result, err) // 42 <nil>
}
```

Small parsing functions are combined into bigger ones: `Padded` wraps `Unsigned`
with optional whitespace skipping, `Try` makes backtracking possible, and
`ParseString` feeds the input text through the resulting combinator.

## Documentation

[GoDoc documentation](https://pkg.go.dev/github.com/okneniz/parsec)

### Backtracking

Greedy combinators (`Satisfy`, `Eq`, `Range`, `OneOf` and others) consume one input item even when they fail. Combinators that try alternatives (`Choice`, `Or`) or repeat parsing in a loop (`Many` and similar) do not restore the buffer position on their own. So whenever a combinator can fail inside such a branch, wrap it in `Try` — it is the backtracking primitive that rewinds the buffer on failure:

```go
choice := strings.Choice(
	"expected number or parens",
	strings.Try(parseNumber()),
	strings.Try(strings.Parens(parseExpr)),
)
```

### Examples

- text
  - [json](https://github.com/okneniz/parsec/tree/master/examples/strings/json)
  - [timestamps](https://github.com/okneniz/parsec/tree/master/examples/strings/timestamps)
  - [credit cards](https://github.com/okneniz/parsec/tree/master/examples/strings/cards)
  - [math expressions](https://github.com/okneniz/parsec/tree/master/examples/strings/math)
- binary
  - [message pack](https://github.com/okneniz/parsec/tree/master/examples/bytes/message_pack)
  - [png](https://github.com/okneniz/parsec/tree/master/examples/bytes/png)

## Development

```bash
make test    # run all tests
make lint    # format code and run golangci-lint
make fmt     # format code only
```

The linter version is pinned in `go.mod` as a [Go tool dependency](https://go.dev/ref/mod#go-tool), so `make lint` uses exactly the same version locally and in CI. To bump it:

```bash
make install-linter # or: go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@<version>
```

## Roadmap

See the [open issues](https://github.com/okneniz/parsec/issues) for a list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.
* If you have suggestions for adding or removing projects, feel free to [open an issue](https://github.com/okneniz/parsec/issues/new) to discuss it, or directly create a pull request after you edit the *README.md* file with necessary changes.
* Please make sure you check your spelling and grammar.
* Create individual PR for each suggestion.

### Creating A Pull Request

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
