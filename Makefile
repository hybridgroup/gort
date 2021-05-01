VERSION := $(shell sed -n 4p version.go | cut -d' ' -f4)

.PHONY: assets build homebrew

assets:
	@cd commands && go-bindata -pkg="commands" support/... && cd ..

build-dir:
	mkdir -p ./build/$(VERSION)
	mkdir -p ./build/$(VERSION)/linux_amd64
	mkdir -p ./build/$(VERSION)/windows
	mkdir -p ./build/$(VERSION)/darwin_amd64
	mkdir -p ./build/$(VERSION)/darwin_arm64

build: build-dir
	GOOS=linux go build -o ./build/$(VERSION)/linux_amd64/gort .
	GOOS=windows go build -o ./build/$(VERSION)/windows/gort.exe .
	GOARCH=amd64 GOOS=darwin go build -o ./build/$(VERSION)/darwin_amd64/gort .
	GOARCH=arm64 GOOS=darwin go build -o ./build/$(VERSION)/darwin_arm64/gort .

package:
	zip ./build/gort-$(VERSION)-linux_amd64.zip ./build/$(VERSION)/linux_amd64/gort
	zip ./build/gort-$(VERSION)-windows.zip ./build/$(VERSION)/windows/gort.exe
	zip ./build/gort-$(VERSION)-darwin_amd64.zip ./build/$(VERSION)/darwin_amd64/gort
	zip ./build/gort-$(VERSION)-darwin_arm64.zip ./build/$(VERSION)/darwin_arm64/gort

release: assets build package

homebrew:
	openssl sha256 < build/gort-$(VERSION)_darwin_amd64.zip
	openssl sha256 < build/gort-$(VERSION)_darwin_386.zip
	openssl sha256 < build/gort-$(VERSION)_darwin_386.zip
