
build:
	go install

go-test:
	tools/go-test-all.sh

acceptance-test:
	./dyntftp &
	tools/create-fixtures.sh
	tools/acceptance.sh
	killall -v -9 dyntftp

test: build unit-test go-test

