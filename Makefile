.PHONY: swremo swremo-client
bin:
	mkdir -p bin

swremo: bin
	GOOS=linux GOARCH=arm GOARM=6 go build -o bin/swremo ./cmd/swremo
