build:
	go build -o bin/binance-pump-detector src/*.go
	cp config.ini.example bin/config.ini

run:
	go run -race src/*.go

clean:
	rm -rf bin

all: build