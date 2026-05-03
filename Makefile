ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=set
    NUL=nul
else
    SET=export
    NUL=/dev/null
endif

GO=go
NAME:=$(notdir $(CURDIR))
VERSION:=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)
GOOPT:=-ldflags "-s -w -X github.com/hymkor/jegan.version=$(VERSION)"
EXE:=$(shell $(GO) env GOEXE)

build:
	$(GO) fmt ./...
	$(SET) "CGO_ENABLED=0" && $(GO) build $(GOOPT)
