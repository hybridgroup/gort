VERSION := $(shell sed -n 3p version.go | cut -d' ' -f4)

.PHONY: release

release:
	@cd commands && go-bindata -pkg="commands" support/... && cd ..
	@goxc -d=./build -bc="linux, windows, darwin" -pv=$(VERSION)
