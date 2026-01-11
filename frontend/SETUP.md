# Frontend Setup Instructions

## Quick Start

1. **Install Node.js** (if not already installed)
   - Download from https://nodejs.org/ (v16 or higher recommended)
   - Verify installation: `node --version` and `npm --version`

2. **Navigate to frontend directory**
   ```bash
   cd frontend
   ```

3. **Install dependencies**
   ```bash
   npm install
   ```

4. **Start the development server**
   ```bash
   npm run dev
   ```

5. **Open your browser**
   - The app will be available at `http://localhost:5173`

## Important Notes

### Backend Connection
- Make sure your Go backend server is running on `http://localhost:3000`
- The frontend is configured to proxy API requests to the backend
- If your backend runs on a different port, update `vite.config.ts`

### Environment Variables
- Create a `.env` file in the `frontend` directory if you need to change the API URL:
  ```
  VITE_API_URL=http://localhost:3000
  ```

### First Time Setup
1. Start the backend server first
2. Then start the frontend development server
3. Register a new user account or login with existing credentials

### Troubleshooting

**Port already in use:**
- Change the port in `vite.config.ts` or use: `npm run dev -- --port 5174`

**API connection errors:**
- Verify backend is running on port 3000
- Check CORS settings in backend if needed
- Verify API endpoints match between frontend and backend

**Build errors:**
- Delete `node_modules` and `package-lock.json`
- Run `npm install` again
- Check Node.js version (should be v16+)

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## Project Structure

```
frontend/
├── src/
│   ├── components/     # Reusable components (Navbar, etc.)
│   ├── pages/          # Page components (Home, Login, Cart, etc.)
│   ├── context/        # React Context (Auth, Cart)
│   ├── services/       # API service layer
│   ├── types/          # TypeScript definitions
│   └── App.tsx         # Main app component
├── public/             # Static files
└── package.json        # Dependencies
```

## Features

✅ User authentication (Login/Register)
✅ Product browsing and search
✅ Shopping cart
✅ Checkout with coupons
✅ Order management
✅ Product reviews
✅ User profile

Enjoy your Amazon clone! 🛍️
