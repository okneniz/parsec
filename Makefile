default: test

test:
	GOTOOLCHAIN=go1.25.0+auto go test -v -timeout 60s -count=1 -coverprofile=coverage.out ./...
	# go test -v -count 1 -timeout 60s -coverprofile=coverage.out ./...

# Pin the linter version in go.mod as a tool dependency.
# Requires Go >= 1.24 (GOTOOLCHAIN=auto downloads it if needed).
install-linter:
	go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2

# Style conventions beyond gofmt, enforced below with awk:
# no one-line function bodies, a blank line between declarations,
# and no blank line directly below a return.
define stylecheck
/^[ \t]*\/\// { if (c == 0) sep = prev; c = 1; prev = $$0; next }
/^[ \t]*func .*\{.+\}[ \t]*(\/\/.*)?$$/ && !/^[ \t]*func \(\).*\{\}[ \t]*$$/ { printf "%s:%d: one-line function body\n", FILENAME, FNR }
/^func / { if (c == 1) { if (sep !~ /^[ \t]*$$/) printf "%s:%d: missing blank line above this declaration\n", FILENAME, FNR } else { if (prev !~ /^[ \t]*$$/) printf "%s:%d: missing blank line above this declaration\n", FILENAME, FNR } }
prev ~ /^[ \t]*return([ \t(]|$$)/ && /^[ \t]*$$/ { printf "%s:%d: blank line below return\n", FILENAME, FNR }
{ c = 0; prev = $$0 }
endef
export stylecheck

GOFILES := $(shell find . -name '*.go' -not -path './.git/*')

# golangci-lint is pinned in go.mod as a tool dependency,
# no manual installation needed: https://go.dev/ref/mod#go-tool
# CGO_ENABLED=0 makes the tool build independent of a local C toolchain.
lint: fmt
	CGO_ENABLED=0 go tool golangci-lint run ./...
	@violations=$$(printf '%s\n' "$$stylecheck" | awk -f - $(GOFILES)); \
	if [ -n "$$violations" ]; then echo "$$violations"; exit 1; fi

fmt:
	CGO_ENABLED=0 go tool golangci-lint fmt ./...

benchmark:
	# go test -v -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out -count=3 -run=^# ./hash-map/...
	go test -v -bench=. -benchmem -count=3 -run=^# ./...


coverage:
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

build:
	go build ./...

pub:
	GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/okneniz/parsec
