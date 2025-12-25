# Examples

This directory contains example code and utilities for AuthKeeper.

## Mock OAuth2 Server

A simple mock OAuth2 server for testing AuthKeeper locally.

### Usage

```bash
# Start the mock server
go run examples/mock-server/main.go
```

The server will start on `http://localhost:8080` with the following endpoint:
- Token endpoint: `http://localhost:8080/oauth/token`

### Test Credentials

You can use any client ID and secret for testing. The mock server accepts all credentials.

Suggested test credentials:
- **Client ID**: `test-client-id`
- **Client Secret**: `test-client-secret`
- **Scopes**: `read write`

### Example Flow

1. Start the mock server:
```bash
go run examples/mock-server/main.go
```

2. Add the mock server to AuthKeeper:
```bash
./authkeeper add
```
Enter:
- Master password: (create a new password)
- Name: `Mock Server`
- Client ID: `test-client-id`
- Client Secret: `test-client-secret`
- Token URL: `http://localhost:8080/oauth/token`
- Scopes: `read write`

3. Get an access token:
```bash
./authkeeper token
```
Select "Mock Server" from the list.

4. The mock server will log the token request and AuthKeeper will display the token.

## Real World Examples

### Auth0

```
Name: Auth0 Production
Client ID: <your-auth0-client-id>
Client Secret: <your-auth0-client-secret>
Token URL: https://your-domain.auth0.com/oauth/token
Scopes: read:users write:users
```

### Keycloak

```
Name: Keycloak Dev
Client ID: <your-keycloak-client-id>
Client Secret: <your-keycloak-client-secret>
Token URL: https://keycloak.example.com/realms/master/protocol/openid-connect/token
Scopes: openid profile email
```

### Okta

```
Name: Okta
Client ID: <your-okta-client-id>
Client Secret: <your-okta-client-secret>
Token URL: https://your-domain.okta.com/oauth2/default/v1/token
Scopes: custom.scope
```

### Google Cloud

```
Name: Google Cloud
Client ID: <your-google-client-id>
Client Secret: <your-google-client-secret>
Token URL: https://oauth2.googleapis.com/token
Scopes: https://www.googleapis.com/auth/cloud-platform
```

### Azure AD

```
Name: Azure AD
Client ID: <your-azure-client-id>
Client Secret: <your-azure-client-secret>
Token URL: https://login.microsoftonline.com/<tenant-id>/oauth2/v2.0/token
Scopes: https://graph.microsoft.com/.default
```
