# Amazon Clone Frontend

A modern, full-featured e-commerce frontend built with React, TypeScript, and Tailwind CSS, designed to look and feel like Amazon.

## Features

- 🛍️ **Product Browsing**: Browse and search products with filtering options
- 🔐 **Authentication**: User registration and login system
- 🛒 **Shopping Cart**: Add, update, and remove items from cart
- 💳 **Checkout**: Complete checkout process with coupon support
- 📦 **Order Management**: View, cancel, and return orders
- ⭐ **Product Reviews**: View and submit product reviews
- 👤 **User Profile**: Manage account information

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Router** - Navigation
- **Tailwind CSS** - Styling
- **Axios** - HTTP client
- **Context API** - State management

## Getting Started

### Prerequisites

- Node.js (v16 or higher)
- npm or yarn

### Installation

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Create a `.env` file in the frontend directory (optional):
```env
VITE_API_URL=http://localhost:3000
```

4. Start the development server:
```bash
npm run dev
```

The application will be available at `http://localhost:5173`

### Building for Production

```bash
npm run build
```

The production build will be in the `dist` directory.

## Project Structure

```
frontend/
├── src/
│   ├── components/      # Reusable UI components
│   ├── pages/          # Page components
│   ├── context/        # React Context providers
│   ├── services/       # API service layer
│   ├── types/          # TypeScript type definitions
│   ├── utils/          # Utility functions
│   ├── App.tsx         # Main app component
│   ├── main.tsx        # Entry point
│   └── index.css      # Global styles
├── public/             # Static assets
├── package.json        # Dependencies
├── tsconfig.json       # TypeScript config
├── vite.config.ts      # Vite config
└── tailwind.config.js  # Tailwind CSS config
```

## API Integration

The frontend connects to the backend API running on `http://localhost:3000` by default. Make sure your backend server is running before starting the frontend.

### API Endpoints Used

- **Auth**: `/auth/login`, `/auth/register/user`
- **Products**: `/auth/products/true/page{page}`, `/auth/products/get/:id`, `/auth/products/filter`
- **Cart**: `/auth/cart/add`, `/auth/cart/view`, `/auth/cart/update`, `/auth/cart/delete`
- **Orders**: `/auth/order/checkout`, `/auth/order/view`, `/auth/order/cancel`, `/auth/order/return`
- **Reviews**: `/auth/review/add`, `/auth/review/view/:product_id`
- **Coupons**: `/auth/coupons/view`

## Features in Detail

### Authentication
- User registration with validation
- Login with role-based access (customer, seller, admin)
- JWT token-based authentication
- Protected routes

### Product Browsing
- Paginated product listings
- Search functionality
- Filter by category, brand, and price range
- Product detail pages with reviews

### Shopping Cart
- Add products to cart
- Update quantities
- Remove items
- Real-time cart count in navbar

### Checkout
- Shipping address input
- Payment method selection
- Coupon code application
- Order placement

### Order Management
- View order history
- Cancel pending orders
- Return delivered orders
- Order status tracking

## Styling

The application uses Tailwind CSS with custom Amazon-inspired colors:
- Amazon Orange: `#FF9900`
- Amazon Dark: `#131921`
- Amazon Light: `#232F3E`
- Amazon Yellow: `#FEBE69`

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build

## Notes

- The frontend expects the backend to be running on port 3000
- JWT tokens are stored in localStorage
- User IDs are extracted from JWT tokens (simplified implementation)
- In production, implement proper token validation and user ID management

## Contributing

This is a full-stack e-commerce application. When making changes:

1. Ensure backend API endpoints match frontend expectations
2. Update TypeScript types if backend structures change
3. Test all user flows (browse, cart, checkout, orders)

## License

This project is part of an Amazon clone application.
