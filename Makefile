GOCMD=go
BINARY_NAME=gaze
OUT=dist
CMD=cmd/gaze/main.go

build:
	go build -v ${CMD}
build-all: build build-darwin build-windows build-linux
build-darwin:
	GOOS=darwin GOARCH=amd64 go build  -ldflags "-s -w" -o ${OUT}/darwin/$(BINARY_NAME) -v ${CMD}
build-windows:
	GOOS=windows GOARCH=amd64 go build  -ldflags "-s -w" -o ${OUT}/windows/$(BINARY_NAME).exe -v ${CMD}
build-linux:
	GOOS=linux GOARCH=amd64 go build  -ldflags "-s -w" -o ${OUT}/linux/$(BINARY_NAME) -v ${CMD}
ut:
	go test github.com/wtetsu/gaze/pkg/...
cov:
	go test -coverprofile=coverage.txt -covermode=atomic github.com/wtetsu/gaze/pkg/...
clean:
	go clean ${CMD}
	rm -f ${OUT}/darwin/$(BINARY_NAME)
	rm -f ${OUT}/windows/$(BINARY_NAME)
	rm -f ${OUT}/linux/$(BINARY_NAME)
