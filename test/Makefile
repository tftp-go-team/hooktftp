# Author: Raphael Medaer
# Date:   2015-05-04
# About:  This Makefile contains procedure for end-to-end tests.
#         This file doesn't contains unit tests.

.PHONY: all
all: check_shell check_files check_web exit

.PHONY: hooktftp
hooktftp:
	../hooktftp -v config.yml & echo "$$!" > hooktftp.pid 
	@echo "Waiting for hooktftp server initialization ..."
	sleep 3

.PHONY: httpd
httpd:
	python -m SimpleHTTPServer 8000 & echo "$$!" > httpd.pid
	@echo "Waiting for HTTP server initialization ..."
	sleep 3 

.PHONY: exit
exit: exit_hooktftp exit_httpd

.PHONY: exit_hooktftp
exit_hooktftp: hooktftp.pid
	-sh -c "kill `cat hooktftp.pid`"
	-rm hooktftp.pid

.PHONY: exit_httpd
exit_httpd: httpd.pid
	-sh -c "kill `cat httpd.pid`"
	-rm httpd.pid

.PHONY: check_shell
check_shell: fixtures shell shell/small shell/medium shell/large shell/limit
	diff fixtures shell

.PHONY: check_files
check_files: fixtures files files/small files/medium files/large files/limit
	diff fixtures files

.PHONY: check_web
check_web: fixtures web web/small web/medium web/large web/limit
	diff fixtures web

fixtures:
	mkdir fixtures

fixtures/small: fixtures
	dd if=/dev/urandom of=fixtures/small bs=4 count=1

fixtures/medium: fixtures
	dd if=/dev/urandom of=fixtures/medium bs=1048577 count=5

fixtures/large: fixtures
	dd if=/dev/urandom of=fixtures/large bs=1048577 count=50

fixtures/limit: fixtures
	dd if=/dev/urandom of=fixtures/limit bs=512 count=65535

shell: 
	mkdir shell

shell/small: shell fixtures/small hooktftp
	atftp --get --remote-file shell/small --local-file shell/small localhost 4242

shell/medium: shell fixtures/medium hooktftp
	atftp --get --remote-file shell/medium --local-file shell/medium localhost 4242

shell/large: shell fixtures/large hooktftp
	atftp --option "blksize 1468" --get --remote-file shell/large --local-file shell/large localhost 4242

shell/limit: shell fixtures/limit hooktftp
	atftp --option "blksize 512" --get --remote-file shell/limit --local-file shell/limit localhost 4242

files:
	mkdir files

files/small: files fixtures/small hooktftp
	atftp --get --remote-file files/small --local-file files/small localhost 4242

files/medium: files fixtures/medium hooktftp
	atftp --get --remote-file files/medium --local-file files/medium localhost 4242

files/large: files fixtures/large hooktftp
	atftp --option "blksize 1468" --get --remote-file files/large --local-file files/large localhost 4242

files/limit: files fixtures/limit hooktftp
	atftp --option "blksize 512" --get --remote-file files/limit --local-file files/limit localhost 4242

web:
	mkdir web

web/small: web fixtures/small hooktftp httpd
	atftp --get --remote-file web/small --local-file web/small localhost 4242

web/medium: web fixtures/medium hooktftp httpd
	atftp --get --remote-file web/medium --local-file web/medium localhost 4242

web/large: web fixtures/large hooktftp httpd
	atftp --option "blksize 1468" --get --remote-file web/large --local-file web/large localhost 4242

web/limit: web fixtures/limit hooktftp httpd
	atftp --option "blksize 512" --get --remote-file web/limit --local-file web/limit localhost 4242


clean: 
	$(RM) -rf fixtures
	$(RM) -rf files
	$(RM) -rf shell
	$(RM) -rf web
