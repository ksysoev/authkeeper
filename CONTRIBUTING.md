# Contributing to AuthKeeper

Thank you for considering contributing to AuthKeeper! Here are some guidelines to help you get started.

## Development Setup

1. **Clone the repository**
```bash
git clone https://github.com/ksysoev/authkeeper.git
cd authkeeper
```

2. **Install dependencies**
```bash
go mod tidy
```

3. **Build the project**
```bash
make build
```

4. **Run tests**
```bash
make test
```

## Project Structure

```
authkeeper/
├── cmd/authkeeper/          # Main application entry point
├── pkg/
│   ├── cmd/                 # Cobra command implementations
│   ├── vault/               # Vault storage and encryption
│   ├── oauth/               # OAuth2 client implementation
│   ├── ui/                  # Simple terminal UI utilities
│   └── crypto/              # Encryption primitives
├── examples/
│   └── mock-server/         # Mock OAuth2 server
├── Makefile                 # Build tasks
└── README.md
```

## Coding Guidelines

### Go Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused

### UI Guidelines
- Use color helpers from `ui` package
- Provide clear prompts and error messages
- Always mask sensitive inputs (passwords, secrets)
- Show progress indicators for slow operations

### Commit Messages
- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, etc.)
- Keep the first line under 72 characters

Example:
```
Add support for refresh token flow

- Implement refresh token storage in vault
- Add new command for token refresh
- Update token model to handle refresh tokens
```

## Adding New Features

### Adding a New OAuth2 Flow

1. Update `pkg/oauth/client.go` with new flow method
2. Update command in `pkg/cmd/commands.go`
3. Use `ui` package for user interaction
4. Update README.md with usage instructions
5. Add tests

### Adding a New Command

1. Create command function in `pkg/cmd/commands.go`
2. Add command to root command in `InitCommand()`
3. Use `ui` package helpers for prompts and output
4. Document in README.md

### Updating UI

All UI utilities are in `pkg/ui/ui.go`. Add new helper functions there for consistent user interaction.

## Testing

### Unit Tests
```bash
make test
```

### Manual Testing
Use the mock server for testing:
```bash
# Terminal 1
go run examples/mock-server/main.go

# Terminal 2
./authkeeper add
# Add mock server credentials

./authkeeper token
# Test token issuance
```

### Testing Checklist
- [ ] All commands work correctly
- [ ] Error cases are handled
- [ ] TUI navigation works smoothly
- [ ] Animations display correctly
- [ ] Vault encryption/decryption works
- [ ] OAuth2 flows succeed
- [ ] Help text is accurate

## Pull Request Process

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/my-new-feature
   ```
3. **Make your changes**
4. **Add tests** for new functionality
5. **Run tests and linting**
   ```bash
   make test
   make lint
   ```
6. **Commit your changes**
7. **Push to your fork**
8. **Open a Pull Request**

### PR Guidelines
- Describe what your PR does
- Reference any related issues
- Include screenshots for UI changes
- Ensure CI passes
- Keep PRs focused on a single feature/fix

## Code Review

All submissions require review. We'll look for:
- Code quality and style
- Test coverage
- Documentation updates
- Backward compatibility

## Questions?

Feel free to open an issue for:
- Bug reports
- Feature requests
- Questions about the codebase
- Documentation improvements

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
