
build:
	go build

tar: build
	rm -rf hooktftp*.tar.gz
	tools/mktar

go-test:
	tools/go-test-all.sh

acceptance-test:
	tools/create-fixtures.sh
	tools/acceptance.sh

test: build go-test acceptance-test
