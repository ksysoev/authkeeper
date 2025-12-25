# AuthKeeper

A secure CLI secret manager and OAuth2/OIDC client with a simple terminal interface.

## Features

- 🔐 **Encrypted Vault** - AES-256-GCM encryption with PBKDF2 key derivation
- 🎯 **OAuth2/OIDC Support** - Client credentials flow implementation
- 🎨 **Simple Interface** - Clean terminal UI with colored output
- 🔑 **Secure Storage** - All credentials encrypted at rest
- 🎫 **Easy Token Issuance** - Quick access token generation
- ⌨️ **Intuitive Prompts** - Straightforward command-line interaction

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

### First Time Setup

When you run `authkeeper add` for the first time, you'll be prompted to create a new vault:

```bash
./authkeeper add
```

**First-time flow:**
1. Create a strong master password (minimum 8 characters)
2. Confirm the password by entering it again
3. Enter your first client details

**Important:** Your master password encrypts all credentials. If you forget it, there's no way to recover your data!

### 1. Add your first OIDC client

After creating the vault (or on subsequent runs):

```bash
./authkeeper add
```

You'll be prompted to:
1. Enter your master password
2. Enter client details:
   - Client Name (e.g., "My Auth Server")
   - Client ID
   - Client Secret (hidden as you type)
   - Token URL
   - Scopes (optional, space-separated)
3. Review and confirm

### 2. Issue an access token

```bash
./authkeeper token
```

Enter your password, select a client from the numbered list, and get your token!

### 3. List all clients

```bash
./authkeeper list
```

### 4. Delete a client

```bash
./authkeeper delete
```

Select a client and confirm deletion.

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

## User Interface

The interface uses colored terminal output for clarity:
- 🟣 **Magenta** - Titles and headings
- 🔵 **Cyan** - Interactive prompts and selections
- 🟢 **Green** - Success messages
- 🔴 **Red** - Error messages
- 🟡 **Yellow** - Warnings
- ⚪ **Gray** - Muted/informational text

All interactions are simple text prompts - just type your answers and press Enter.

## Architecture

```
authkeeper/
├── cmd/authkeeper/          # Main application entry
├── pkg/
│   ├── cmd/                 # Cobra commands
│   ├── vault/               # Encrypted vault management
│   ├── oauth/               # OAuth2 client
│   ├── ui/                  # Terminal UI utilities
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
- **Terminal I/O**: golang.org/x/term - Password input handling
- **Crypto**: golang.org/x/crypto - Encryption primitives
- **OAuth2**: golang.org/x/oauth2 - OAuth2 support

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

## Interface Style

The CLI uses a simple, straightforward text interface:
- Clear prompts for user input
- Colored output for better readability
- Password fields masked with bullets (•)
- Boxed output for structured data
- Confirmation prompts for destructive operations
- Spinner indicators for operations that take time
