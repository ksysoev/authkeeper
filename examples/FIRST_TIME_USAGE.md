# First-Time Usage Example

This document shows the complete first-time experience with AuthKeeper.

## Prerequisites

Start the mock OAuth2 server:

```bash
go run examples/mock-server/main.go
```

Output:
```
🔐 Mock OAuth2 Server starting on :8080
Token endpoint: http://localhost:8080/oauth/token
Test credentials:
  Client ID: test-client-id
  Client Secret: test-client-secret
```

## First Run: Adding Your First Client

Run the add command:

```bash
./authkeeper add
```

### Step 1: Vault Creation

Since this is the first time, you'll be prompted to create a vault:

```
🔐 Add OIDC Client

First time setup - creating encrypted vault

🔐 Create New Vault

You're creating a new vault. Please choose a strong master password.
⚠ This password encrypts all your credentials - don't forget it!

Enter master password: ••••••••••••
Confirm master password: ••••••••••••
```

If passwords don't match:
```
✗ Error: Passwords do not match. Please try again.

Enter master password: ••••••••••••
Confirm master password: ••••••••••••
```

If password is too short:
```
✗ Error: Password must be at least 8 characters long

Enter master password: ••••••••••••
Confirm master password: ••••••••••••
```

When successful:
```
✓ Master password set successfully!
```

### Step 2: Add Client Details

Now enter your first client:

```
Enter client credentials

Client Name: Mock Server
Client ID: test-client-id
Client Secret: ••••••••••••••••••
Token URL: http://localhost:8080/oauth/token
Scopes (optional, space-separated): read write

Review client details:
┌────────────────────────────────────────────────────────┐
│ Name:         Mock Server                              │
│ Client ID:    test-client-id                           │
│ Client Secret: ••••••••••••••••••                      │
│ Token URL:    http://localhost:8080/oauth/token       │
│ Scopes:       read, write                              │
└────────────────────────────────────────────────────────┘

Save this client? (y/n): y
Saving to encrypted vault... 
✓ Client added successfully!
```

## Subsequent Runs

On subsequent runs, you'll only be asked for your password once:

```bash
./authkeeper add
```

```
🔐 Add OIDC Client

Enter master password to unlock vault
Master Password: ••••••••••••

Enter client credentials
...
```

## Testing Other Commands

### List Clients

```bash
./authkeeper list
```

```
📋 OIDC Clients

Enter master password to unlock vault
Master Password: ••••••••••••

Found 1 client(s)

1. Mock Server
┌────────────────────────────────────────────────────────┐
│ Client ID:  test-client-id                             │
│ Token URL:  http://localhost:8080/oauth/token         │
│ Scopes:     read, write                                │
│ Created:    2024-12-25 10:30:15                        │
└────────────────────────────────────────────────────────┘
```

### Get Token

```bash
./authkeeper token
```

```
🎫 Issue Access Token

Enter master password to unlock vault
Master Password: ••••••••••••

Select OIDC client:

1. Mock Server

Enter number: 1

Fetching access token... 
✓ Token issued successfully!

Client: Mock Server

┌────────────────────────────────────────────────────────┐
│ Access Token:                                          │
│                                                        │
│ mock_token_test-client-id_12345                       │
│                                                        │
│ Token Type: Bearer                                     │
│ Expires In: 3600 seconds                               │
│ Scope: read write                                      │
└────────────────────────────────────────────────────────┘

💡 Tip: Copy the access token to use in your API requests
```

### Delete Client

```bash
./authkeeper delete
```

```
🗑️  Delete OIDC Client

Enter master password to unlock vault
Master Password: ••••••••••••

Select client to delete:

1. Mock Server

Enter number: 1

⚠ Are you sure you want to delete 'Mock Server'?
This action cannot be undone.

Delete this client? (y/n): y
Deleting client... 
✓ Client deleted successfully!
```

## Before Vault Creation

If you try to use list/token/delete before creating the vault:

```bash
./authkeeper list
```

```
📋 OIDC Clients

⚠ Vault not found
Use 'authkeeper add' to create vault and add your first client
```

## Tips

1. **Strong Password**: Use at least 12 characters with mixed case, numbers, and symbols
2. **Password Manager**: Store your master password in a secure password manager
3. **Backup**: Backup `~/.authkeeper/vault.enc` (it's already encrypted)
4. **First Command**: Always use `authkeeper add` first to create the vault
5. **Wrong Password**: If you get "failed to decrypt vault", you entered the wrong password
6. **Lost Password**: There's no password recovery - you'll need to start over
