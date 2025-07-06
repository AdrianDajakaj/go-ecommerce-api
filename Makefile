# Go development tasks

.PHONY: lint test format pre-commit

# Lint code with golangci-lint
lint:
	@echo "📝 Running golangci-lint..."
	@golangci-lint run

# Format code
format:
	@echo "🎨 Formatting code..."
	@gofmt -w .
	@goimports -w .

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./... -v

# Pre-commit checks
pre-commit: format lint test
	@echo "✅ All pre-commit checks passed!"

# Install tools
install-tools:
	@echo "🔧 Installing development tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
