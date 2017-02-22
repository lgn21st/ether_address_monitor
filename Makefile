build:
	go build -o ether_address_monitor

run: build
	./geth-relay

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ether_address_monitor

clean:
	go clean
