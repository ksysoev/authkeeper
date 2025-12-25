#!/bin/bash

# Test first-time vault creation
echo "Testing first-time vault creation..."
echo ""

# Start mock server in background
echo "Starting mock OAuth2 server..."
go run examples/mock-server/main.go > /tmp/mock-server.log 2>&1 &
MOCK_PID=$!
sleep 2

# Cleanup function
cleanup() {
    echo ""
    echo "Cleaning up..."
    kill $MOCK_PID 2>/dev/null
    rm -f ~/.authkeeper/vault.enc
    echo "Done!"
}
trap cleanup EXIT

# Test add command - this will show the new vault creation flow
echo "Run: ./authkeeper add"
echo "This should prompt for password confirmation on first run"
echo ""

# Note: This is a manual test - user needs to interact
# For automated testing, we'd need expect or similar

