TARGET = hlybot
USERNAME = bevisy

.PHONY: \
	build \
	clean \
	install \
	image \
	push \
	clean_image

build:
	CGO_ENABLED=0 go build -o ${TARGET}

clean:
	rm -f ${TARGET}

install:
	cp ${TARGET} $GOPATH/bin/

image: build clean_image
	docker build -t docker.io/${USERNAME}/${TARGET}:latest -f Dockerfile .

push:
	docker push docker.io/${USERNAME}/${TARGET}:latest

clean_image:
	docker rmi ${USERNAME}/${TARGET}:latest
