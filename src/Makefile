
.PHONY: build
build: deps hooktftp

hooktftp: tftp/parse_rrq.go tftp/rrq_response.go tftp/server.go tftp/rrq_test.go tftp/consts.go tftp/parse_rrq_test.go regexptransform/regexp_transform_test.go regexptransform/regexp_transform.go drop_privileges.go hooks/file_hook.go hooks/hooks_test.go hooks/shell_hook.go hooks/http_hook.go hooks/compile.go config/parse_yaml.go config/parse_yaml_test.go logger/logger.go hooktftp.go
	go build -a -v -o hooktftp

.PHONY: deps
deps:
	go get


.PHONY: test
test: test_regexptransform test_tftp test_config test_hooks

.PHONY: test_regexptransform
test_regexptransform: regexptransform/regexp_transform.go regexptransform/regexp_transform_test.go
	cd regexptransform; go test

.PHONY: test_tftp
test_tftp: tftp/parse_rrq.go tftp/rrq_response.go tftp/server.go tftp/rrq_test.go tftp/consts.go tftp/parse_rrq_test.go
	cd tftp; go test

.PHONY: test_config
test_config: config/parse_yaml.go config/parse_yaml_test.go 
	cd config; go test

.PHONY: test_hooks
test_hooks: hooks/file_hook.go hooks/hooks_test.go hooks/shell_hook.go hooks/http_hook.go hooks/compile.go
	cd hooks; go test

.PHONY: clean
clean:
	go clean
