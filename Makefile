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
	docker run --rm -ti -v $(pwd):/go/src/github.com/tftp-go-team/hooktftp -w /app --name hooktftp hooktftp-shell bash
