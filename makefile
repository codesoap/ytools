.PHONY: install uninstall

INSTALLDIR ?= /usr/local

install:
	GOBIN="${INSTALLDIR}/bin" go install ./...
	mkdir -p "${INSTALLDIR}/man/man7"
	install -m 644 "man/ytools.7" "${INSTALLDIR}/man/man7"

uninstall:
	rm -f "${INSTALLDIR}/bin/ytools-search" \
		"${INSTALLDIR}/bin/ytools-pick" \
		"${INSTALLDIR}/bin/ytools-info" \
		"${INSTALLDIR}/bin/ytools-recommend" \
		"${INSTALLDIR}/man/man7/ytools.7"
