# Amazon Clone Backend (Go + Fiber + MongoDB)

This project is a backend for an Amazon-like e-commerce platform, built using Go, Fiber, and MongoDB.

## Project Structure

```

main.go
    - Main function where the program starsts 
config/
    - Database configurations

errors/
    - defined errors to reuse

modules/
  auth/
    - CRUD operations for Users, Sellers and Admin
    - Registeration Routes
    - Login Routes
    - Added Pagination

  cart/
    - CRUD operations for Cart.
    - View Cart Records.
    - Clear Cart 
   
  coupons/
    - CRUD operations for Coupons.
    - View Coupons 
    - Coupons Validation and Expiry.
   
  order/
    - CRUD operations for Order product.
    - Cancel and Return Simulation.
    - Order Status after Checkout.
    
  product/
    - CRUD operations on Products.
    - List Products with pagination.
    - List Products which are in-stocks and out-off stocks.
    - Filter Products based on Category, Price, Brand etc.

  review/
    - CRUD operations on Review. 
    - Listing Average Review and Comments.
    
  websocket/
    - Functions to Trigger Notification alerts.
    - Alerts after order placed and order cancelled.

router/
    - Group of particular module routes.
    ...
```

## Getting Started

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```
2. **Set up MongoDB** and configure your connection string in `config/database.go`.
3. **Run the server**:
   ```bash
   go run main.go
   ```

## Features
- Authentication and Validation
- Pagination and limit
- Structured Response
- Product, Cart, Order, Coupon, Review, and WebSocket modules
- Modular code structure for scalability

Feel free to contribute or open issues for improvements!
