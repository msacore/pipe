help:
	@echo "pipe  - Go module to easily work with pipelines"
	@echo
	@echo "help  - Show this information"
	@echo "doc   - Serve Go HTML documentation"
	@echo "test  - Run Go tests"
.PHONY: help

doc:
	@echo "Navigate to http://localhost:6060/pkg/github.com/msacore/pipe"
	~/go/bin/godoc -http=:6060
.PHONY: doc

test:
	go test ./...
.PHONY: test