
.PHONY: build
build: deps hooktftp

hooktftp:
	go build -a -v -o hooktftp

.PHONY: deps
deps:
	go get

.PHONY: clean
clean:
	-rm hooktftp
	-rm src
