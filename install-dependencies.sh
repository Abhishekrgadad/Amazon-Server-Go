#!/bin/bash

echo "=========================================="
echo "Installing Dependencies for Amazon Clone"
echo "=========================================="
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo "This script requires sudo privileges."
    echo "Please run: sudo bash install-dependencies.sh"
    exit 1
fi

echo "Step 1: Updating package list..."
apt update

echo ""
echo "Step 2: Installing Node.js and npm..."
apt install -y nodejs npm

echo ""
echo "Step 3: Installing Go..."
apt install -y golang-go

echo ""
echo "Step 4: Verifying installations..."
echo "Node.js version:"
node --version
echo "npm version:"
npm --version
echo "Go version:"
go version

echo ""
echo "=========================================="
echo "Installation Complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Install frontend dependencies: cd frontend && npm install"
echo "2. Install Go backend dependencies: go mod download"
echo ""
