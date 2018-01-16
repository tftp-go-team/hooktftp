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
