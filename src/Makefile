prefix = /usr/local
exec_prefix = $(prefix)
bindir = $(exec_prefix)/bin

INSTALL = install
INSTALL_PROGRAM = $(INSTALL)
INSTALL_DATA = $(INSTALL) -m 644

build: deps
	go build

deps:
	go get

tar: build
	rm -rf hooktftp*.tar.gz
	tools/mktar

go-test:
	tools/go-test-all.sh

acceptance-test:
	tools/create-fixtures.sh
	tools/acceptance.sh

test: build go-test acceptance-test

install:
	mkdir -p $(DESTDIR)$(bindir)
	$(INSTALL_PROGRAM) -t $(DESTDIR)$(bindir) hooktftp

