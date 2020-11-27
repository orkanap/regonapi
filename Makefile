.PHONY: test test-unit

# with integration tests
test:
	go test -v -tags=integration

# unit tests only
test-unit:
	go test -v