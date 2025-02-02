GOCMD=go
BINARY_NAME=gaze
OUT=dist
CMD=cmd/gaze/main.go

build:
	${GOCMD} build -ldflags "-s -w" -v ${CMD}
build-all: build-macos-amd build-macos-arm build-windows build-linux
build-all-amd: build-macos-amd build-windows build-linux
build-macos-amd:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o ${OUT}/macos_amd/$(BINARY_NAME) -v ${CMD}
build-macos-arm:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o ${OUT}/macos_arm/$(BINARY_NAME) -v ${CMD}
build-windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${OUT}/windows/$(BINARY_NAME).exe -v ${CMD}
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${OUT}/linux/$(BINARY_NAME) -v ${CMD}
ut:
	${GOCMD} test github.com/wtetsu/gaze/pkg/...
e2e:
	${GOCMD} build -ldflags "-s -w" -o test/e2e/ -v ${CMD}
	cd test/e2e && sh test_all.sh
cov:
	${GOCMD} test -coverprofile=coverage.txt -covermode=atomic github.com/wtetsu/gaze/pkg/...
clean:
	${GOCMD} clean ${CMD}
	rm -f ${OUT}/macos-amd/$(BINARY_NAME)
	rm -f ${OUT}/macos-arm/$(BINARY_NAME)
	rm -f ${OUT}/windows/$(BINARY_NAME)
	rm -f ${OUT}/linux/$(BINARY_NAME)
check-cross-compile:
	@echo "Checking cross-compiled binaries..."
	@file ${OUT}/linux/gaze       | grep -c "ELF 64-bit LSB executable, x86-64" | grep -q "1" || (echo "Error: Linux binary is not correctly built"     && exit 1)
	@file ${OUT}/macos_amd/gaze   | grep -c "Mach-O 64-bit x86_64 executable"   | grep -q "1" || (echo "Error: macOS amd binary is not correctly built" && exit 1)
	@file ${OUT}/macos_arm/gaze   | grep -c "Mach-O 64-bit arm64 executable"    | grep -q "1" || (echo "Error: macOS arm binary is not correctly built" && exit 1)
	@file ${OUT}/windows/gaze.exe | grep -c "x86-64, for MS Windows"            | grep -q "1" || (echo "Error: Windows binary is not correctly built"   && exit 1)
	@echo "Cross-compilation check passed!"