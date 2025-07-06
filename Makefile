# Go development tasks

.PHONY: lint test format pre-commit lint-format-only

# Lint only formatting (for commits with compilation errors)
lint-format-only:
	@echo "🎨 Checking formatting..."
	@gofmt -d . | tee /tmp/gofmt.out
	@if [ -s /tmp/gofmt.out ]; then echo "❌ Code not formatted!"; exit 1; fi
	@echo "✅ Code is properly formatted"

# Full lint with golangci-lint (when code compiles)
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

# Pre-commit checks (lenient for development)
pre-commit: format lint-format-only
	@echo "✅ Pre-commit checks passed!"

# Full pre-commit checks (when ready for production)
pre-commit-full: format lint test
	@echo "✅ All pre-commit checks passed!"

# Install tools
install-tools:
	@echo "🔧 Installing development tools..."
	@go install golang.org/x/tools/cmd/goimports@latest
