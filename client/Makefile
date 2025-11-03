.PHONY: test test-unit test-e2e test-coverage test-verbose test-short help

COVERAGE_FILE = coverage.out
COVERAGE_HTML = coverage.html

test: test-unit

test-unit: format
	@go test ./...

test-coverage: format
	@go test -short -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

clean:
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@rm -f *.test
	@find . -name "*.test" -type f -delete

format:
	@go fmt ./...

test-race: format
	@go test -short -race ./...
