
build:
	go install

unit-test:
	go test
	cd tftp
	go test
	cd ..

acceptance-test:
	dyntftp -config test_config.json -port 1234 -root . &
	tools/create-fixtures.sh
	tools/acceptance.sh
	killall -v -9 dyntftp

test: build unit-test acceptance-test

