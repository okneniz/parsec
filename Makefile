default: test

test:
	GOTOOLCHAIN=go1.25.0+auto go test -v -timeout 60s -count=1 -coverprofile=coverage.out ./...
	# go test -v -count 1 -timeout 60s -coverprofile=coverage.out ./...

# Pin the linter version in go.mod as a tool dependency.
# Requires Go >= 1.24 (GOTOOLCHAIN=auto downloads it if needed).
install-linter:
	go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2

# golangci-lint is pinned in go.mod as a tool dependency,
# no manual installation needed: https://go.dev/ref/mod#go-tool
# CGO_ENABLED=0 makes the tool build independent of a local C toolchain.
lint: fmt
	CGO_ENABLED=0 go tool golangci-lint run ./...

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
