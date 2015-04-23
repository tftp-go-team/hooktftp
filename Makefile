


.PHONY: build
build:
	$(MAKE) -C src

.PHONY: test
test: build
	$(MAKE) -C test

.PHONY: clean
clean:
	$(MAKE) -C src clean
	$(MAKE) -C test clean

