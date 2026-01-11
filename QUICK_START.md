# Quick Start Guide

## One-Command Installation

Run this command to install everything:

```bash
sudo bash install-dependencies.sh
```

Then install project dependencies:

```bash
# Install frontend dependencies
cd frontend && npm install && cd ..

# Install backend dependencies  
go mod download
```

## Running the Application

### Terminal 1 - Backend
```bash
go run main.go
```

### Terminal 2 - Frontend
```bash
cd frontend
npm run dev
```

### Open Browser
Navigate to: `http://localhost:5173`

## Prerequisites Checklist

- [ ] Node.js installed (v16+)
- [ ] npm installed
- [ ] Go installed (v1.18+)
- [ ] MongoDB installed and running
- [ ] `.env` file configured with MongoDB connection

## Common Commands

```bash
# Check versions
node --version
npm --version
go version

# Install frontend deps
cd frontend && npm install

# Install backend deps
go mod download

# Run backend
go run main.go

# Run frontend
cd frontend && npm run dev

# Build frontend for production
cd frontend && npm run build
```
