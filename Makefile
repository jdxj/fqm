# linux, windows, darwin
.PHONY: build
build.%: file := fqm
build.%:
	CGO_ENABLED=0 GOOS=$* GOARCH=amd64 go build -ldflags '-s -w' -o $(file)_$*.out ./*.go
	upx --best $(file)_$*.out
