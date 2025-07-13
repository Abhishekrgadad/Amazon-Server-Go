# Amazon Clone Backend (Go + Fiber + MongoDB)

This project is a backend for an Amazon-like e-commerce platform, built using Go, Fiber, and MongoDB.

## Project Structure

```
go.mod
go.sum
index.html
main.go
config/
  database.go
  jwt.go
errors/
  functions.go
All modules follows the Below structure as auth
modules/
  auth/
    - Crud operations for Users, Sellers and Admin(single)
    - Registeration Routes
    - Login Routes
    - Added Pagination
  cart/
   
  coupons/
   
  order/
    
  product/
    
  review/
    
  websocket/
router/
  router.go
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
- Authentication (JWT)
- Product, Cart, Order, Coupon, Review, and WebSocket modules
- Modular code structure for scalability

---

Feel free to contribute or open issues for improvements!
