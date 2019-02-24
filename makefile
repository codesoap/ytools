.PHONY: all

ROOT = github.com/codesoap/ytools

all: bin/ytools-audio bin/ytools-search

bin/ytools-audio: src/cmd/ytools-audio/main.go
	go build -o bin/ytools-audio ${ROOT}/src/cmd/ytools-audio

bin/ytools-search: src/cmd/ytools-search/main.go
	go build -o bin/ytools-search ${ROOT}/src/cmd/ytools-search
