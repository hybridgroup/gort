VERSION := $(shell sed -n 3p version.go | cut -d' ' -f4)

.PHONY: release homebrew

release:
	@cd commands && go-bindata -pkg="commands" support/... && cd ..
	@goxc -d=./build -bc="linux, windows, darwin" -pv=$(VERSION)

homebrew:
	openssl sha256 < build/$(VERSION)/gort_$(VERSION)_darwin_amd64.zip
	openssl sha256 < build/$(VERSION)/gort_$(VERSION)_darwin_386.zip
