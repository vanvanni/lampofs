.PHONY: test bench test-cov

test:
	gotestsum --format testname ./

bench:
	go test -bench=. ./

test-cov:
	gotestsum --format testname -- -cover ./
