VERSION := $(shell sed -n 3p version.go | cut -d' ' -f4)

.PHONY: assets build homebrew

assets:
	@cd commands && go-bindata -pkg="commands" support/... && cd ..

build:
	GOBIN="$(GOPATH)/bin" goxc -d=./build -bc="linux, windows, darwin" -pv=$(VERSION)

release: assets build

homebrew:
	openssl sha256 < build/$(VERSION)/gort_$(VERSION)_darwin_amd64.zip
	openssl sha256 < build/$(VERSION)/gort_$(VERSION)_darwin_386.zip

deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/laher/goxc
