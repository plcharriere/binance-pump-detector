build:
	go build -o bin/binance-pump-detector src/*.go

run:
	go run -race src/*.go

clean:
	rm -rf bin

all: build