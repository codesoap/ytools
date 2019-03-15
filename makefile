.PHONY: all install uninstall clean

PREFIX = /usr/local
MANPREFIX = /usr/local/man
ROOT = github.com/codesoap/ytools

all: bin/ytools-search bin/ytools-pick bin/ytools-info bin/ytools-recommend bin/ytools-comments

install: all
	mkdir -p "${DESTDIR}${PREFIX}/bin"
	install -m 755 "bin/ytools-search" "${DESTDIR}${PREFIX}/bin"
	install -m 755 "bin/ytools-pick" "${DESTDIR}${PREFIX}/bin"
	install -m 755 "bin/ytools-info" "${DESTDIR}${PREFIX}/bin"
	install -m 755 "bin/ytools-recommend" "${DESTDIR}${PREFIX}/bin"
	install -m 755 "bin/ytools-comments" "${DESTDIR}${PREFIX}/bin"
	mkdir -p "${DESTDIR}${MANPREFIX}/man7"
	install -m 644 "man/ytools.7" "${DESTDIR}${MANPREFIX}/man7"

uninstall:
	rm -f "${DESTDIR}${PREFIX}/bin/ytools-search" \
		"${DESTDIR}${PREFIX}/bin/ytools-pick" \
		"${DESTDIR}${PREFIX}/bin/ytools-info" \
		"${DESTDIR}${PREFIX}/bin/ytools-recommend" \
		"${DESTDIR}${PREFIX}/bin/ytools-comments" \
		"${DESTDIR}${MANPREFIX}/man7/ytools.7"

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

bin/ytools-comments: src/cmd/ytools-comments/main.go src/ytools/common.go
	go build -o bin/ytools-comments ${ROOT}/src/cmd/ytools-comments
