GOCMD=go
BINARY_NAME=gaze
OUT=dist
CMD=cmd/gaze/main.go
VERSION := $(shell cat cmd/gaze/version)
PLATFORMS = macos_amd macos_arm windows linux

build:
	${GOCMD} build -ldflags "-s -w" -v ${CMD}
build-all: build-macos-amd build-macos-arm build-windows build-linux
build-all-amd: build-macos-amd build-windows build-linux
build-macos-amd:
	GOOS=darwin GOARCH=amd64  ${GOCMD} build -ldflags "-s -w" -o ${OUT}/gaze_macos_amd_${VERSION}/${BINARY_NAME}   -v ${CMD}
build-macos-arm:
	GOOS=darwin GOARCH=arm64  ${GOCMD} build -ldflags "-s -w" -o ${OUT}/gaze_macos_arm_${VERSION}/${BINARY_NAME}   -v ${CMD}
build-windows:
	GOOS=windows GOARCH=amd64 ${GOCMD} build -ldflags "-s -w" -o ${OUT}/gaze_windows_${VERSION}/${BINARY_NAME}.exe -v ${CMD}
build-linux:
	GOOS=linux GOARCH=amd64   ${GOCMD} build -ldflags "-s -w" -o ${OUT}/gaze_linux_${VERSION}/${BINARY_NAME}       -v ${CMD}
ut:
	${GOCMD} test github.com/wtetsu/gaze/pkg/...
e2e:
	${GOCMD} build -ldflags "-s -w" -o test/e2e/ -v ${CMD}
	cd test/e2e && sh test_all.sh
cov:
	${GOCMD} test -coverprofile=coverage.txt -covermode=atomic github.com/wtetsu/gaze/pkg/...
clean:
	${GOCMD} clean ${CMD}
	@for plat in $(PLATFORMS); do \
		rm -rf ${OUT}/gaze_$${plat}_${VERSION}; \
		rm -f  ${OUT}/gaze_$${plat}_${VERSION}.zip; \
	done
check-cross-compile:
	@echo "Checking cross-compiled binaries..."
	@file ${OUT}/gaze_macos_amd_${VERSION}/${BINARY_NAME}   | grep -c "Mach-O 64-bit x86_64 executable"   | grep -q "1" || (echo "Error: macOS amd binary is not correctly built" && exit 1)
	@file ${OUT}/gaze_macos_arm_${VERSION}/${BINARY_NAME}   | grep -c "Mach-O 64-bit arm64 executable"    | grep -q "1" || (echo "Error: macOS arm binary is not correctly built" && exit 1)
	@file ${OUT}/gaze_windows_${VERSION}/${BINARY_NAME}.exe | grep -c "x86-64, for MS Windows"            | grep -q "1" || (echo "Error: Windows binary is not correctly built"   && exit 1)
	@file ${OUT}/gaze_linux_${VERSION}/${BINARY_NAME}       | grep -c "ELF 64-bit LSB executable, x86-64" | grep -q "1" || (echo "Error: Linux binary is not correctly built"     && exit 1)
	@echo "Cross-compilation check passed!"
package: clean build-all check-cross-compile
	@for plat in $(PLATFORMS); do \
		cp LICENSE README.md ${OUT}/gaze_$${plat}_${VERSION}; \
		(cd ${OUT} && zip -r gaze_$${plat}_${VERSION}.zip ./gaze_$${plat}_${VERSION}); \
	done
