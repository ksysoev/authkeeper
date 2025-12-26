# Testing Documentation

This project uses **testify** for assertions and **mockery** for generating mocks. The test suite provides comprehensive coverage of all packages while maintaining code readability.

## Test Coverage Summary

| Package | Coverage | Description |
|---------|----------|-------------|
| `pkg/cmd` | 73.3% | Command-line interface commands |
| `pkg/core` | 67.2% | Core business logic and service layer |
| `pkg/prov` | 96.0% | OAuth2 provider implementation |
| `pkg/repo` | 83.9% | Encrypted vault repository |
| `pkg/ui` | 30.5% | CLI user interface (limited due to interactive nature) |
| **Total** | **47.7%** | Overall project coverage |

## Running Tests

### Run all tests
```bash
make test
```

### Run tests with coverage report
```bash
make test-coverage
```

This will:
1. Run all tests with race detector
2. Generate `coverage.out` file
3. Generate `coverage.html` for browser viewing
4. Display total coverage percentage

### View coverage in browser
```bash
open coverage.html
```

## Test Structure

### Unit Tests

Each package has comprehensive unit tests:

- **`pkg/core/service_test.go`**: Tests for core business logic
  - Service creation
  - Client management (add, get, list, delete)
  - Token issuance
  - Repository initialization checks
  - Error handling scenarios

- **`pkg/prov/oauth_test.go`**: Tests for OAuth2 provider
  - Token request with client credentials
  - HTTP request validation
  - Error handling (4xx, 5xx responses)
  - Context cancellation
  - Network errors
  - Invalid JSON responses

- **`pkg/repo/vault_test.go`**: Tests for encrypted vault
  - Save and retrieve clients
  - Encryption/decryption
  - Duplicate name handling
  - Password validation
  - Empty file handling
  - Multiple operations

- **`pkg/cmd/app_test.go`**: Tests for CLI commands
  - Command creation
  - Command structure validation

- **`pkg/ui/cli_test.go`**: Tests for UI components
  - Service integration
  - Print functions
  - Error handling

## Mocking

### Mock Configuration

Mocks are configured in `.mockery.yaml` and generated automatically using mockery:

```yaml
packages:
  github.com/ksysoev/authkeeper/pkg/core:
    interfaces:
      Repository:
      Provider:
  github.com/ksysoev/authkeeper/pkg/ui:
    interfaces:
      CoreService:
```

### Generate Mocks

To regenerate mocks after interface changes:

```bash
make generate-mocks
```

### Mock Files

Generated mocks are placed in the same package as the interface:
- `pkg/core/repository_mock.go`
- `pkg/core/provider_mock.go`
- `pkg/ui/core_service_mock.go`

## Test Patterns

### Table-Driven Tests

All tests use table-driven patterns for better organization:

```go
tests := []struct {
    name        string
    input       InputType
    setupMock   func(*MockType)
    expected    ExpectedType
    expectedErr string
}{
    {
        name: "successful case",
        // test definition
    },
    // more test cases
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test execution
    })
}
```

### Mock Expectations

Using testify's mock package with expecter pattern:

```go
repo := NewMockRepository(t)
repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil)
```

### Assertions

Using testify's assert package:

```go
assert.NoError(t, err)
assert.Equal(t, expected, actual)
assert.Contains(t, str, substring)
```

## Coverage Goals

- **Core business logic** (pkg/core, pkg/prov, pkg/repo): Target 80%+ coverage
- **UI components** (pkg/ui): Lower coverage due to interactive nature
- **Commands** (pkg/cmd): Focus on structure validation

## Continuous Integration

Tests run automatically on:
- Pull requests
- Push to main branch
- Manual workflow dispatch

CI configuration includes:
- Race detector enabled
- Multiple Go versions tested
- Coverage reports generated

## Best Practices

1. **Test Independence**: Each test is independent and can run in any order
2. **Clear Naming**: Test names describe what they test
3. **Focused Tests**: Each test verifies one specific behavior
4. **Mock Isolation**: Mocks are created per test to avoid interference
5. **Error Testing**: Both success and failure paths are tested
6. **Readable**: Tests are easy to understand and maintain

## Adding New Tests

When adding new functionality:

1. Write tests first (TDD approach recommended)
2. Update mocks if interfaces change: `make generate-mocks`
3. Run tests: `make test`
4. Check coverage: `make test-coverage`
5. Ensure coverage doesn't drop significantly

## Troubleshooting

### Mock Generation Issues

If mocks are not generating correctly:
```bash
go clean -modcache
go mod download
make generate-mocks
```

### Race Condition Warnings

Race detector warnings indicate potential concurrency issues. Fix them before merging.

### Coverage Drops

If coverage drops unexpectedly:
1. Check which lines are not covered: `open coverage.html`
2. Add tests for uncovered paths
3. Consider if untested code is actually needed
