terse: ./url.go ./server.go ./cmd/terse/terse.go
	go build ./cmd/terse

clean:
	rm -f ./terse
