.PHONY: all install clean

ROOT = github.com/codesoap/ytools

all: bin/ytools-search bin/ytools-pick bin/ytools-info bin/ytools-recommend

install: all
	cp "bin/ytools-search" "${HOME}/bin"
	cp "bin/ytools-pick" "${HOME}/bin"
	cp "bin/ytools-info" "${HOME}/bin"
	cp "bin/ytools-recommend" "${HOME}/bin"

clean:
	rm -rf bin

bin/ytools-search: src/cmd/ytools-search/main.go src/ytools/common.go
	go build -o bin/ytools-search ${ROOT}/src/cmd/ytools-search

bin/ytools-pick: src/cmd/ytools-pick/main.go src/ytools/common.go
	go build -o bin/ytools-pick ${ROOT}/src/cmd/ytools-pick

bin/ytools-info: src/cmd/ytools-info/main.go src/ytools/common.go
	go build -o bin/ytools-info ${ROOT}/src/cmd/ytools-info

bin/ytools-recommend: src/cmd/ytools-recommend/main.go src/ytools/common.go
	go build -o bin/ytools-recommend ${ROOT}/src/cmd/ytools-recommend
