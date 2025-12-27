# AuthKeeper

[![Tests](https://github.com/ksysoev/authkeeper/actions/workflows/ci.yml/badge.svg)](https://github.com/ksysoev/authkeeper/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/ksysoev/authkeeper/graph/badge.svg?token=NBZY1JOHTK)](https://codecov.io/gh/ksysoev/authkeeper)
[![Go Report Card](https://goreportcard.com/badge/github.com/ksysoev/authkeeper)](https://goreportcard.com/report/github.com/ksysoev/authkeeper)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A secure CLI secret manager and OAuth2/OIDC client with a simple terminal interface. Store OAuth2/OIDC credentials in an encrypted vault and quickly issue access tokens using the client credentials flow.

## Features

- 🔐 **Encrypted Vault** - AES-256-GCM encryption with PBKDF2 key derivation
- 🎯 **OAuth2/OIDC Support** - Client credentials flow implementation
- 🎨 **Simple Interface** - Clean terminal UI with colored output
- 🔑 **Secure Storage** - All credentials encrypted at rest
- 🎫 **Easy Token Issuance** - Quick access token generation
- ⌨️ **Intuitive Prompts** - Straightforward command-line interaction

## Installation

### Downloading binaries:

Compiled executables can be downloaded from [here](https://github.com/ksysoev/authkeeper/releases).

### Install from source code:

```bash
go install github.com/ksysoev/authkeeper/cmd/authkeeper@latest
```

### Install with homebrew:

```bash
brew tap ksysoev/authkeeper
brew install authkeeper
```

### Build from source:

```bash
# Clone the repository
git clone https://github.com/ksysoev/authkeeper.git
cd authkeeper

# Build the binary
make build

# Or install directly
make install
```

## Usage

### First Time Setup

When you run `authkeeper add` for the first time, you'll be prompted to create a new vault:

```bash
authkeeper add
```

**First-time flow:**
1. Create a strong master password (minimum 8 characters)
2. Confirm the password by entering it again
3. Enter your first client details

**Important:** Your master password encrypts all credentials. If you forget it, there's no way to recover your data!

### Add an OIDC client

```bash
authkeeper add
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

### Issue an access token

```bash
authkeeper token
```

Enter your password, select a client from the numbered list, and get your token!

### List all clients

```bash
authkeeper list
```

### Delete a client

```bash
authkeeper delete
```

Select a client and confirm deletion.

### Example

Try AuthKeeper with a local mock OAuth2 server:

```bash
# Terminal 1: Start mock OAuth2 server
go run examples/mock-server/main.go

# Terminal 2: Add the mock server as a client
authkeeper add
# Use these values:
# - Name: Mock Server
# - Client ID: test-client-id
# - Client Secret: test-client-secret
# - Token URL: http://localhost:8080/oauth/token
# - Scopes: read write

# Get an access token
authkeeper token
```

## Commands

| Command | Description |
|---------|-------------|
| `authkeeper add` | Add a new OIDC client to vault |
| `authkeeper token` | Issue access token for a client |
| `authkeeper list` | List all stored clients |
| `authkeeper delete` | Delete a client from vault |
| `authkeeper --help` | Show help information |

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

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development

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

## License

AuthKeeper is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
