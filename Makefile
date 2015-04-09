
generate:
	@go generate ./itest

./itest/decls_list.go:
	@go install
	@go generate ./itest

test: ./itest/decls_list.go
	@go test ./itest -cover

bench: ./itest/decls_list.go
	@go test ./itest -bench=.

clean:
	$(RM) ./itest/decls_list.go