# AuthKeeper

A secure CLI secret manager and OAuth2/OIDC client with a beautiful TUI built with Bubble Tea.

## Features

- 🔐 **Encrypted Vault** - AES-256-GCM encryption with PBKDF2 key derivation
- 🎯 **OAuth2/OIDC Support** - Client credentials flow implementation
- 🎨 **Beautiful TUI** - Minimalistic and fancy interface with animations
- 🔑 **Secure Storage** - All credentials encrypted at rest
- 🎫 **Easy Token Issuance** - Quick access token generation
- ⌨️ **Intuitive Navigation** - Keyboard-driven interface

## Installation

```bash
# Clone the repository
git clone https://github.com/ksysoev/authkeeper.git
cd authkeeper

# Build the binary
make build

# Or install directly
make install
```

## Quick Start

### 1. Add your first OIDC client

```bash
./authkeeper add
```

You'll be prompted to:
1. Enter your master password (creates vault if first time)
2. Enter client details:
   - Client Name (e.g., "My Auth Server")
   - Client ID
   - Client Secret
   - Token URL
   - Scopes (optional)

### 2. Issue an access token

```bash
./authkeeper token
```

Select your client from the list and get an access token instantly!

### 3. List all clients

```bash
./authkeeper list
```

### 4. Delete a client

```bash
./authkeeper delete
```

## Demo with Mock Server

Try AuthKeeper with a local mock OAuth2 server:

```bash
# Terminal 1: Start mock OAuth2 server
go run examples/mock-server/main.go

# Terminal 2: Add the mock server as a client
./authkeeper add
# Use these values:
# - Name: Mock Server
# - Client ID: test-client-id
# - Client Secret: test-client-secret
# - Token URL: http://localhost:8080/oauth/token
# - Scopes: read write

# Get an access token
./authkeeper token
```

## Commands

| Command | Description |
|---------|-------------|
| `authkeeper add` | Add a new OIDC client to vault |
| `authkeeper token` | Issue access token for a client |
| `authkeeper list` | List all stored clients |
| `authkeeper delete` | Delete a client from vault |
| `authkeeper --help` | Show help information |

## Keyboard Navigation

- **Enter** - Confirm/Continue
- **Tab/↓** - Next field/item
- **Shift+Tab/↑** - Previous field/item
- **Esc** - Cancel/Quit
- **Y/N** - Confirm/Cancel (delete operations)
- **q** - Quit (where applicable)

## Architecture

```
authkeeper/
├── cmd/authkeeper/          # Main application entry
├── internal/
│   ├── cmd/                 # Cobra commands
│   ├── vault/               # Encrypted vault management
│   ├── oauth/               # OAuth2 client
│   └── tui/                 # Bubble Tea UI models
├── pkg/
│   └── crypto/              # Encryption utilities
└── examples/
    └── mock-server/         # Mock OAuth2 server for testing
```

## Security

### Encryption
- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Derivation**: PBKDF2 with 100,000 iterations
- **Salt**: 32 bytes random salt per vault
- **Nonce**: Random nonce per encryption operation

### Best Practices
- Master password is never stored on disk
- Vault file permissions: 0600 (read/write owner only)
- All sensitive data encrypted at rest
- Memory is cleared after use where possible

### Vault Location
Default vault location: `~/.authkeeper/vault.enc`

## Development

```bash
# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Build
make build

# Clean
make clean
```

## Tech Stack

- **Framework**: [Cobra](https://github.com/spf13/cobra) - CLI framework
- **TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- **Components**: [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- **Crypto**: golang.org/x/crypto - Encryption primitives

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details

## Roadmap

- [ ] Authorization code flow support
- [ ] Device code flow support
- [ ] Token refresh mechanism
- [ ] Export/import vault
- [ ] Multiple vault support
- [ ] Token caching
- [ ] Shell completion scripts

## Screenshots

*Note: Since this is a TUI application, it's best experienced in your terminal!*

The interface features:
- 🎨 Purple/violet color scheme
- ⚡ Smooth animations during operations
- 📦 Rounded border boxes
- 🔄 Spinner animations for async operations
- ✨ Focused input highlighting
