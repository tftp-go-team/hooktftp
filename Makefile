
build:
	go install

go-test:
	tools/go-test-all.sh

acceptance-test:
	tools/create-fixtures.sh
	tools/acceptance.sh

test: build go-test acceptance-test

