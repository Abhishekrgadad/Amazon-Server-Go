# Installation Instructions

## Quick Installation

Run the installation script with sudo:

```bash
sudo bash install-dependencies.sh
```

## Manual Installation

If you prefer to install manually, follow these steps:

### 1. Install Node.js and npm (for React Frontend)

```bash
sudo apt update
sudo apt install -y nodejs npm
```

Verify installation:
```bash
node --version
npm --version
```

**Note**: The default Ubuntu repository may have an older version of Node.js. For the latest version, you can use NodeSource:

```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs
```

### 2. Install Go (for Backend)

```bash
sudo apt install -y golang-go
```

Verify installation:
```bash
go version
```

**Note**: For the latest Go version, you can download from https://go.dev/dl/

### 3. Install Frontend Dependencies

```bash
cd frontend
npm install
```

### 4. Install Backend Dependencies

```bash
# From the project root directory
go mod download
```

## Running the Project

### Start Backend Server

```bash
# From project root
go run main.go
```

The backend will run on `http://localhost:3000`

### Start Frontend Server

Open a new terminal:

```bash
cd frontend
npm run dev
```

The frontend will run on `http://localhost:5173`

## Troubleshooting

### Node.js version too old

If you get errors about Node.js version, install a newer version:

```bash
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs
```

### Go version too old

If you need a newer Go version:

1. Download from https://go.dev/dl/
2. Extract: `sudo tar -C /usr/local -xzf go*.tar.gz`
3. Add to PATH: Add `export PATH=$PATH:/usr/local/go/bin` to `~/.bashrc`
4. Reload: `source ~/.bashrc`

### Permission Errors

If you encounter permission errors with npm:

```bash
sudo chown -R $USER:$USER ~/.npm
```

### Port Already in Use

If port 3000 or 5173 is already in use:

**Backend**: Change port in `.env` file or `main.go`
**Frontend**: Change port in `vite.config.ts` or run: `npm run dev -- --port 5174`

## Verification

After installation, verify everything works:

```bash
# Check Node.js
node --version  # Should show v16.x or higher
npm --version   # Should show 8.x or higher

# Check Go
go version      # Should show go1.18 or higher

# Check frontend dependencies
cd frontend
npm list        # Should show all packages installed

# Check backend dependencies
go list -m all  # Should show all Go modules
```

## Next Steps

1. Make sure MongoDB is installed and running (required for backend)
2. Create a `.env` file in the project root with your MongoDB connection string
3. Start the backend server
4. Start the frontend server
5. Open `http://localhost:5173` in your browser

## Need Help?

- Check the `README.md` files in both frontend and root directories
- Review `SETUP.md` in the frontend directory
- Check backend logs for errors
- Check browser console for frontend errors
