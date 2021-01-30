.PHONY: build

build:
	go build -o hlyBot

clean:
	rm -f hlyBot

install:
	cp hlyBot $GOPATH/bin/
