# AuthKeeper Usage Guide

Complete guide to using AuthKeeper for managing OAuth2/OIDC credentials.

## Table of Contents

1. [First Time Setup](#first-time-setup)
2. [Managing Clients](#managing-clients)
3. [Issuing Tokens](#issuing-tokens)
4. [Security Best Practices](#security-best-practices)
5. [Troubleshooting](#troubleshooting)

## First Time Setup

When you run AuthKeeper for the first time, it will create a vault in `~/.authkeeper/vault.enc`.

### 1. Create Your Master Password

**First-time experience:**

```bash
./authkeeper add
```

You'll see:
```
🔐 Create New Vault

First time setup - creating encrypted vault
You're creating a new vault. Please choose a strong master password.
⚠ This password encrypts all your credentials - don't forget it!

Enter master password: ••••••••••••
Confirm master password: ••••••••••••
✓ Master password set successfully!
```

The master password requirements:
- **Minimum 8 characters**
- Choose a strong, unique password
- You'll need to enter it twice for confirmation
- Store it in a secure password manager

**Important**: If you lose your master password, there is no way to recover your credentials!

## Managing Clients

### Adding a Client

```bash
./authkeeper add
```

**First time (vault creation):**
1. **Create Master Password** (minimum 8 characters)
   - Enter password
   - Confirm password
   - Password must match

**Subsequent times:**
1. **Enter Master Password**
   - Enter your existing password to unlock vault

**Then for both:**
2. **Enter Client Details**
   - **Name**: Friendly name for the client (e.g., "Production Auth Server")
   - **Client ID**: Your OAuth2 client identifier
   - **Client Secret**: Your OAuth2 client secret (hidden as you type)
   - **Token URL**: The OAuth2 token endpoint
   - **Scopes**: Space-separated list of scopes (optional)

3. **Confirm Details**
   - Review all information
   - Press y to save or n to cancel

**Example:**
```
Master Password: ••••••••••••
Client Name: Production API
Client ID: prod-client-abc123
Client Secret: ••••••••••••
Token URL: https://auth.example.com/oauth/token
Scopes: read write admin
```

### Listing Clients

```bash
./authkeeper list
```

Shows all stored clients with:
- Name
- Client ID
- Token URL
- Scopes
- Creation timestamp

### Deleting a Client

```bash
./authkeeper delete
```

**Safety features:**
- Select client from list
- Confirmation required (press Y)
- Cancel anytime with N or Esc

## Issuing Tokens

### Get an Access Token

```bash
./authkeeper token
```

**Process:**
1. Enter master password
2. Select client from list
3. Wait for token fetch (with spinner animation)
4. Copy the access token

**Token Information Displayed:**
- Access token (full JWT or opaque token)
- Token type (usually "Bearer")
- Expires in (seconds)
- Granted scopes

### Using the Token

The token can be used in HTTP Authorization headers:

```bash
# Example with curl
TOKEN="<your-access-token>"
curl -H "Authorization: Bearer $TOKEN" https://api.example.com/resource
```

### Token Expiration

- Tokens are not cached - they're fetched fresh each time
- Check the "Expires In" field to know how long the token is valid
- Re-run `authkeeper token` when the token expires

## Security Best Practices

### Master Password

✅ **DO:**
- Use a unique, strong password
- Store it in a password manager
- Never share it
- Change it periodically

❌ **DON'T:**
- Use the same password as other services
- Write it down in plain text
- Share your vault file without changing the password first

### Vault File

Location: `~/.authkeeper/vault.enc`

**Permissions:**
- Owner read/write only (0600)
- Never commit to version control
- Back up encrypted (the file is already encrypted)

**Backup:**
```bash
# Safe to backup the encrypted vault
cp ~/.authkeeper/vault.enc ~/backups/vault.enc.backup

# To restore
cp ~/backups/vault.enc.backup ~/.authkeeper/vault.enc
```

### Client Secrets

- Never log or print client secrets
- Rotate them regularly at the OAuth2 provider
- Update AuthKeeper after rotation
- Delete old clients you no longer use

### Tokens

- Tokens are displayed in plain text - be careful where you run the command
- Don't commit tokens to git
- Don't share tokens
- Use the minimum scopes necessary

## Troubleshooting

### "Failed to decrypt vault (wrong password?)"

**Cause:** Incorrect master password

**Solution:**
- Try your password again carefully
- Check if Caps Lock is on
- If you've forgotten the password, you'll need to delete the vault and start over

### "Client with name already exists"

**Cause:** Trying to add a client with a duplicate name

**Solution:**
- Choose a different name
- Or delete the existing client first
- Or list clients to see what names are already used

### "Token request failed with status 401"

**Cause:** Invalid client credentials

**Solution:**
- Verify client ID and secret are correct
- Check if credentials were rotated at the provider
- Update the client in AuthKeeper
- Verify the token URL is correct

### "Token request failed with status 400"

**Cause:** Bad request - possibly invalid scopes or grant type

**Solution:**
- Check if the scopes are supported by the provider
- Verify the token URL is correct
- Ensure the provider supports client_credentials flow

### "No clients found in vault"

**Cause:** Vault is empty

**Solution:**
- Add a client using `authkeeper add`

### Connection Refused / Timeout

**Cause:** Cannot reach the token endpoint

**Solution:**
- Check network connectivity
- Verify the token URL
- Check if firewall is blocking the request
- Try accessing the URL in a browser

## Advanced Usage

### Multiple Vaults

Currently, AuthKeeper uses a single vault. To use multiple vaults:

```bash
# Backup current vault
mv ~/.authkeeper/vault.enc ~/.authkeeper/vault-prod.enc

# Use for different environment
mv ~/.authkeeper/vault-staging.enc ~/.authkeeper/vault.enc
./authkeeper token

# Switch back
mv ~/.authkeeper/vault.enc ~/.authkeeper/vault-staging.enc
mv ~/.authkeeper/vault-prod.enc ~/.authkeeper/vault.enc
```

### Scripting

For automation, you can parse the token output:

```bash
# This is not recommended for production - tokens will be in shell history!
# Use for development/testing only

./authkeeper token | grep "Access Token" | cut -d' ' -f3
```

**Better approach:** Use service accounts with limited scopes for automation.

### Testing with Mock Server

For development and testing:

```bash
# Terminal 1: Start mock server
go run examples/mock-server/main.go

# Terminal 2: Test AuthKeeper
./authkeeper add
# Add mock server (http://localhost:8080/oauth/token)

./authkeeper token
# Get test token
```

## Tips and Tricks

1. **Organize Client Names**
   - Use prefixes: `prod-api`, `staging-api`, `dev-api`
   - Include environment: `Auth0-Production`, `Auth0-Staging`

2. **Keep Scopes Minimal**
   - Request only the scopes you need
   - Different clients for different scope sets

3. **Regular Maintenance**
   - Review clients monthly
   - Delete unused clients
   - Rotate credentials regularly

4. **Quick Navigation**
   - Use arrow keys or vim-style (j/k) for navigation
   - Tab for next field
   - Shift+Tab for previous field

5. **Keyboard Shortcuts**
   - `Ctrl+C` or `Esc`: Quick exit
   - `Enter`: Confirm/Continue
   - `q`: Quit (where applicable)

## Getting Help

```bash
# Show help
./authkeeper --help

# Show version
./authkeeper --version

# Command-specific help
./authkeeper add --help
./authkeeper token --help
```

## Support

- 📖 [README](README.md) - Quick start guide
- 🤝 [CONTRIBUTING](CONTRIBUTING.md) - Development guide
- 💡 [Examples](examples/README.md) - Example configurations
- 🐛 Issues - Report bugs or request features on GitHub
