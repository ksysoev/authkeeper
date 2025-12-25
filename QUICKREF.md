# AuthKeeper Quick Reference

## Commands

| Command | Description | Shortcut |
|---------|-------------|----------|
| `authkeeper add` | Add new OIDC client | - |
| `authkeeper token` | Get access token | - |
| `authkeeper list` | List all clients | - |
| `authkeeper delete` | Delete a client | - |
| `authkeeper --help` | Show help | `-h` |
| `authkeeper --version` | Show version | `-v` |

## Keyboard Shortcuts

### Navigation
- `↓` / `j` - Move down
- `↑` / `k` - Move up
- `Tab` - Next field
- `Shift+Tab` - Previous field
- `Enter` - Confirm/Select
- `Esc` - Cancel/Back
- `Ctrl+C` - Force quit
- `q` - Quit (where available)

### Delete Confirmation
- `Y` - Confirm delete
- `N` - Cancel delete

## File Locations

| Item | Location | Permissions |
|------|----------|-------------|
| Vault | `~/.authkeeper/vault.enc` | 0600 |
| Directory | `~/.authkeeper/` | 0700 |

## Security Features

- 🔐 AES-256-GCM encryption
- 🔑 PBKDF2 key derivation (100,000 iterations)
- 🧂 32-byte random salt
- 🎲 Random nonce per operation
- 🔒 Password never stored on disk

## Token Response Fields

```
✓ Token issued successfully!
Client: Production API

┌─────────────────────────────────────┐
│ Access Token:                       │
│                                     │
│ eyJhbGciOiJSUzI1NiIsInR5cCI6...   │
│                                     │
│ Token Type: Bearer                  │
│ Expires In: 3600 seconds           │
│ Scope: read write                   │
└─────────────────────────────────────┘
```

## OAuth2 Client Credentials Flow

```
┌──────────────┐         ┌──────────────┐
│  AuthKeeper  │────────▶│ Token Server │
│              │  POST   │              │
│              │◀────────│              │
└──────────────┘  Token  └──────────────┘
```

**Request:**
```http
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials
&client_id=...
&client_secret=...
&scope=...
```

**Response:**
```json
{
  "access_token": "...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "scope": "read write"
}
```

## Common Token URLs

### Auth0
```
https://YOUR_DOMAIN.auth0.com/oauth/token
```

### Keycloak
```
https://keycloak.example.com/realms/REALM_NAME/protocol/openid-connect/token
```

### Okta
```
https://YOUR_DOMAIN.okta.com/oauth2/default/v1/token
```

### Google Cloud
```
https://oauth2.googleapis.com/token
```

### Azure AD
```
https://login.microsoftonline.com/TENANT_ID/oauth2/v2.0/token
```

## Using Access Tokens

### cURL
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     https://api.example.com/resource
```

### HTTPie
```bash
http GET https://api.example.com/resource \
     Authorization:"Bearer YOUR_TOKEN"
```

### JavaScript (fetch)
```javascript
fetch('https://api.example.com/resource', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
```

### Python (requests)
```python
import requests

headers = {'Authorization': f'Bearer {token}'}
response = requests.get('https://api.example.com/resource', headers=headers)
```

## Color Scheme

- 🟣 Primary: `#7C3AED` (Purple)
- 🟪 Secondary: `#A78BFA` (Light Purple)
- 🟢 Success: `#10B981` (Green)
- 🔴 Error: `#EF4444` (Red)
- ⚪ Muted: `#6B7280` (Gray)

## Status Indicators

- `⠋` - Loading/Processing (animated spinner)
- `✓` - Success
- `✗` - Error
- `▶` - Selected item
- `💡` - Tip/Information
- `⚠️` - Warning

## Error Codes

| Status | Description | Solution |
|--------|-------------|----------|
| 400 | Bad Request | Check scopes and grant type |
| 401 | Unauthorized | Verify client credentials |
| 403 | Forbidden | Check client permissions |
| 404 | Not Found | Verify token URL |
| 500 | Server Error | Check provider status |

## Backup Strategy

```bash
# Backup vault
cp ~/.authkeeper/vault.enc ~/backup/vault-$(date +%Y%m%d).enc

# Restore vault
cp ~/backup/vault-20241225.enc ~/.authkeeper/vault.enc
```

## Environment Variables

Currently, AuthKeeper doesn't use environment variables. All configuration is interactive through the TUI.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error |

## Tips

1. ⌨️ **Use Tab** to navigate between fields quickly
2. 🔍 **Check the help text** at the bottom of each screen
3. 🎨 **Focused fields** are highlighted in purple
4. ⏰ **Watch the spinner** during async operations
5. 💾 **Backup your vault** regularly
6. 🔐 **Strong password** = secure vault
7. 📝 **Organize names** for easy selection
8. 🔄 **Rotate credentials** periodically

---

**Version:** dev  
**License:** MIT  
**Repository:** github.com/ksysoev/authkeeper
