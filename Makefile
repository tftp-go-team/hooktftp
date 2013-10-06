
build:
	go install

go-test:
	cd regexptransform/
	go test
	cd ..

	cd tftp
	go test
	cd ..

	cd config
	go test
	cd ..

	go test

acceptance-test:
	dyntftp -config test_config.json -port 1234 -root . &
	tools/create-fixtures.sh
	tools/acceptance.sh
	killall -v -9 dyntftp

test: build unit-test go-test

