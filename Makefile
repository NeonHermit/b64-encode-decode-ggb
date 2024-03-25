.PHONY: build-windows clean

build-windows:
	export CGO_ENABLED=1 && \
	export CC=x86_64-w64-mingw32-gcc && \
	export GOOS=windows && \
	export GOARCH=amd64 && \
	go build -ldflags -H=windowsgui -o b64-geogb.exe ./cmd

clean:
	rm -f b64-geogb.exe
