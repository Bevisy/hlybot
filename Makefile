.PHONY: \
	build \
	clean \
	install \
	image \
	push

build:
	go build -o hlyBot

clean:
	rm -f hlyBot

install:
	cp hlyBot $GOPATH/bin/

image:
	docker build -t docker.io/bevisy/hlyBot:latest -f Dockerfile .

push:
	docker push docker.io/bevisy/hlyBot:latest