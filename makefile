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

bin/ytools-search: cmd/ytools-search/main.go common.go
	go build -o bin/ytools-search ${ROOT}/cmd/ytools-search

bin/ytools-pick: cmd/ytools-pick/main.go common.go
	go build -o bin/ytools-pick ${ROOT}/cmd/ytools-pick

bin/ytools-info: cmd/ytools-info/main.go common.go
	go build -o bin/ytools-info ${ROOT}/cmd/ytools-info

bin/ytools-recommend: cmd/ytools-recommend/main.go common.go
	go build -o bin/ytools-recommend ${ROOT}/cmd/ytools-recommend

bin/ytools-comments: cmd/ytools-comments/main.go common.go
	go build -o bin/ytools-comments ${ROOT}/cmd/ytools-comments
