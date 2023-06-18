default: test

test:
	go test -v -coverprofile=coverage.out ./...

benchmark:
	# go test -v -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out -count=3 -run=^# ./hash-map/...
	go test -v -bench=. -count=3 -run=^# ./...

coverage:
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
