.PHONY: build-windows build-macos clean

build-windows:
	export CGO_ENABLED=1 && \
	export CC=x86_64-w64-mingw32-gcc && \
	export GOOS=windows && \
	export GOARCH=amd64 && \
	go build -ldflags -H=windowsgui -o b64-geogb.exe ./cmd

build-macos:
	export CG0_ENABLED=1 && \
	export GOOS=darwin && \
	export GOARCH=arm64 && \
	go build -ldflags="-s -w" -o b64-geogb ./cmd

clean:
	rm -f b64-geogb.exe b64-geogb
