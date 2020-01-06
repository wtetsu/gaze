GOCMD=go
BINARY_NAME=gaze
OUT=dist
CMD=cmd/gaze/main.go

all: test build
build-all: build build-osx build-windows build-linux
build:
	go build -o ${OUT}/$(BINARY_NAME) -v ${CMD}
build-osx:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ${OUT}/osx/$(BINARY_NAME) -v ${CMD}
build-windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o ${OUT}/windows/$(BINARY_NAME).exe -v ${CMD}
build-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ${OUT}/linux/$(BINARY_NAME) -v ${CMD}
test:
	go test -v ./...
clean:
	go clean ${CMD}
	rm -f ${OUT}/osx/$(BINARY_NAME)
	rm -f ${OUT}/windows/$(BINARY_NAME)
	rm -f ${OUT}/linux/$(BINARY_NAME)
