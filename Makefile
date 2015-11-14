terse: ./util.go ./server.go ./cmd/terse/terse.go
	go build ./cmd/terse

run:
	go run ./cmd/terse/terse.go

clean:
	rm -f ./terse


.PHONY: run clean
