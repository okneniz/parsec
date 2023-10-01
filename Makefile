default: test

test:
	go test -v -timeout 60s -coverprofile=coverage.out ./...
	# go test -v -count 1 -timeout 60s -coverprofile=coverage.out ./...

intall-linter:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

lint: fmt
	golangci-lint run ./...

fmt:
	gofmt -w -s .

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
