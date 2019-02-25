.PHONY: all clean

ROOT = github.com/codesoap/ytools

all: bin/ytools-pick bin/ytools-search

clean:
	rm -r bin

bin/ytools-pick: src/cmd/ytools-pick/main.go
	go build -o bin/ytools-pick ${ROOT}/src/cmd/ytools-pick

bin/ytools-search: src/cmd/ytools-search/main.go
	go build -o bin/ytools-search ${ROOT}/src/cmd/ytools-search
