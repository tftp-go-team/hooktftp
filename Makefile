.PHONY: build
build:
	$(MAKE) -C src

.PHONY: test
test: build
	$(MAKE) -C src test
	$(MAKE) -C test all

.PHONY: clean
clean:
	$(MAKE) -C src clean
	$(MAKE) -C test clean

.PHONY: gox
gox:
	$(MAKE) -C src gox

shell:
	docker build -t hooktftp-shell .
	docker run --rm -ti -v $(shell pwd):/go/src/github.com/tftp-go-team/hooktftp -w /go/src/github.com/tftp-go-team/hooktftp --name hooktftp hooktftp-shell bash

release:
	docker build -t tftpgoteam/hooktftp:latest .
	docker push tftpgoteam/hooktftp:latest
