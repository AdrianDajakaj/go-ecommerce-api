# E-Commerce API Documentation

## Overview
This project implements a Clean Architecture approach for an e-commerce system written in Go. It uses the Echo framework for HTTP routing, GORM as the ORM layer, and SQLite as the database engine.

## Requirements

- Go 1.20+
- SQLite3 (bundled, no separate installation required)

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/AdrianDajakaj/go-ecommerce-api.git
cd go-ecommerce-api
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run the server

```bash
go run cmd/server.go
```

Server will be available at: `http://localhost:8080`

## Configuration

By default, the SQLite file is named `ecommerce.db` (in the project root).

- To change the database file, edit cmd/server.go:

```go
dsn := "ecommerce.db" 
```

- To change the HTTP port, adjust:

```go
e.Start(":8080")
```

## Authentication & Authorization

This API is protected by JWT and role-based access control:

1. Register & Login

- `POST /users/register`
  - Creates a new user with role `"user"`.
  - Extra fields like `"role"` in the JSON payload are ignored.

- POST `/users/login`
  - Expects JSON:
  ```json
  { "email": "...", "password": "..." }
  ```
  - Returns a JSON Web Token (`token`) on success.

2. Roles

- Each user has a `role` field: `"user"` or `"admin"`.
- By default, newly registered users get `"user"`.
- To grant admin privileges, manually update the `role` in the SQLite database:
```bash
sqlite3 ecommerce.db <<SQL
UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';
.exit
SQL
```
- Endpoints requiring admin privileges will check the token’s `role`.

3. JWT Middleware

- Protected routes expect an `Authorization` header:
```bash
Authorization: <JWT_TOKEN>
```
- If the token is missing, invalid, or expired, the API returns `401 Unauthorized`.
- If a user tries to access an admin-only endpoint without `"admin"` role, the API returns `403 Forbidden`.

## Data Models & JSON Samples

### User

```json
{
  "email": "john@example.com",
  "password": "secure123",
  "name": "John",
  "surname": "Doe",
  "address": {
    "country": "USA",
    "city": "New York",
    "postcode": "10001",
    "street": "Broadway",
    "number": "1"
  }
}
```

### Product

```json
{
  "name": "Phone",
  "description": "Smartphone",
  "price": 299.99,
  "stock": 100,
  "is_active": true,
  "category_id": 1,
  "images": [
    { "url": "https://example.com/images/phone-front.jpg" },
    { "url": "https://example.com/images/phone-back.jpg" }
  ]
}
```

### Category

```json
{
  "name": "Electronics",
  "parent_id": 1
}
```

### Cart Add Item

```json
{
  "product_id": 1,
  "quantity": 2
}
```

### Order Creation

```json
{
  "payment_method": "CARD",
  "shipping_address_id": 1
}
```

## Endpoint Patterns

### Users

| Method | Path              | Protected? | Roles Allowed    | Description                                     |
| ------ | ----------------- | ---------- | ---------------- | ----------------------------------------------- |
| POST   | `/users/register` | No         | —                | Register new user (`role` defaults to `"user"`) |
| POST   | `/users/login`    | No         | —                | Login and receive JWT                           |
| GET    | `/users`          | Yes (JWT)  | `admin`          | Get all users                                   |
| GET    | `/users/{id}`     | Yes (JWT)  | `admin` or owner | Get user by ID                                  |
| GET    | `/users/search?…` | Yes (JWT)  | `admin`          | Search users with query parameters              |
| PUT    | `/users/{id}`     | Yes (JWT)  | `admin` or owner | Update user profile                             |
| DELETE | `/users/{id}`     | Yes (JWT)  | `admin` or owner | Delete user                                     |

### Catehories

| Method | Path                             | Protected? | Roles Allowed | Description                             |
| ------ | -------------------------------- | ---------- | ------------- | --------------------------------------- |
| GET    | `/categories`                    | No         | —             | Get all categories                      |
| GET    | `/categories/{id}`               | No         | —             | Get category by ID                      |
| GET    | `/categories/{id}/subcategories` | No         | —             | Get subcategories of a category         |
| GET    | `/categories/search?…`           | No         | —             | Search categories with query parameters |
| POST   | `/categories`                    | Yes (JWT)  | `admin`       | Create new category                     |
| PUT    | `/categories/{id}`               | Yes (JWT)  | `admin`       | Update category                         |
| DELETE | `/categories/{id}`               | Yes (JWT)  | `admin`       | Delete category                         |

### Products

| Method | Path                 | Protected? | Roles Allowed | Description                           |
| ------ | -------------------- | ---------- | ------------- | ------------------------------------- |
| GET    | `/products`          | No         | —             | Get all products                      |
| GET    | `/products/{id}`     | No         | —             | Get product by ID                     |
| GET    | `/products/search?…` | No         | —             | Search products with query parameters |
| POST   | `/products`          | Yes (JWT)  | `admin`       | Create new product                    |
| PUT    | `/products/{id}`     | Yes (JWT)  | `admin`       | Update product                        |
| DELETE | `/products/{id}`     | Yes (JWT)  | `admin`       | Delete product                        |

### Carts

All `/cart` endpoints require JWT.
- Regular users see/modify only their own cart items.
- Admin can also filter/search all carts.

| Method | Path                   | Protected? | Roles Allowed     | Description                                       |
| ------ | ---------------------- | ---------- | ----------------- | ------------------------------------------------- |
| GET    | `/cart`                | Yes (JWT)  | `user` or `admin` | Get authenticated user's cart                     |
| POST   | `/cart/add`            | Yes (JWT)  | `user` or `admin` | Add product to authenticated user's cart          |
| PUT    | `/cart/item/{item_id}` | Yes (JWT)  | `user` or `admin` | Update quantity of a cart item (owner/admin only) |
| DELETE | `/cart/item/{item_id}` | Yes (JWT)  | `user` or `admin` | Remove a cart item (owner/admin only)             |
| DELETE | `/cart/clear`          | Yes (JWT)  | `user` or `admin` | Clear authenticated user's cart                   |
| GET    | `/cart/search?…`       | Yes (JWT)  | `user` or `admin` | Search carts: admin sees all; user sees own only  |

### Orders

All `/orders` endpoints require JWT.
- `GetOrder`, `CancelOrder` check ownership or admin.
- `UpdateStatus` only for admin.
- `Search` for users always filters to their own orders (ignores `user_id`)`; admin can search all.

| Method | Path                  | Protected? | Roles Allowed      | Description                                           |
| ------ | --------------------- | ---------- | ------------------ | ----------------------------------------------------- |
| POST   | `/orders`             | Yes (JWT)  | `user` or `admin`  | Create order from authenticated user's cart           |
| GET    | `/orders/{id}`        | Yes (JWT)  | `owner` or `admin` | Get order by ID (owner or admin only)                 |
| GET    | `/orders`             | Yes (JWT)  | `admin`            | Get all orders                                        |
| GET    | `/orders/user`        | Yes (JWT)  | `user` or `admin`  | Get authenticated user's orders (admin sees only own) |
| PUT    | `/orders/{id}/status` | Yes (JWT)  | `admin`            | Update order status                                   |
| PUT    | `/orders/{id}/cancel` | Yes (JWT)  | `owner` or `admin` | Cancel order (owner or admin; owner only if pending)  |
| GET    | `/orders/search?…`    | Yes (JWT)  | `user` or `admin`  | Search orders: admin sees all; user sees own only     |

## Scopes (Filtering via Query Parameters)

These scopes apply to `search` endpoints:

### User Scopes
- `email=<value>` — exact match
- `name=<value>` — contains
- `surname=<value>` — contains
- `country=<value>` — exact
- `city=<value>` — exact

### Product Scopes
- `name=<value>` — contains
- `category_id=<id>` — exact
- `is_active=<true|false>` — exact
- `price_min=<n>&price_max=<m>` — range
- `with_category=true` — eager-load Category object

### Category Scopes
- `name=<value>` — contains
- `created_after=<RFC3339 timestamp>` — ≥ date
- `created_before=<RFC3339 timestamp>` — ≤ date
- `min_products=<n>` — minimum number of products
- `parent_id=<id>` — exact
- `with_products=true` — eager-load Products array
- `with_subcategories=true` — eager-load Subcategories array

### Order Scopes
- `user_id=<id>` — exact (ignored for regular users)
- `status=<value>` — exact (e.g., PENDING, PAID, CANCELLED)
- `created_after=<RFC3339 timestamp>` — ≥ date
- `total_min=<n>&total_max=<m>` — range

### Cart Scopes
- `user_id=<id>` — exact (ignored for regular users)
- `total_min=<n>&total_max=<m>` — range
- `created_before=<RFC3339 timestamp>` — ≤ date

## cURL Examples (with JWT & Roles)

Below are example flows demonstrating registration, role assignment, and protected endpoint usage. Replace `localhost:8080` if your server runs on a different host/port.

### 1. Register & Login
  1. Register two users (both default to `role = "user"`)

  ```bash
  # Register `admin@example.com`
  curl -s -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Password123!",
    "name": "Jan",
    "surname": "Nowak",
    "address": {
      "street": "Ul. Przykładowa 1",
      "city": "Warszawa",
      "zip": "00-001",
      "country": "PL"
    }
  }' | jq

  # Register `user@example.com`
  curl -s -X POST http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!",
    "name": "Anna",
    "surname": "Kowalska",
    "address": {
      "street": "Ul. Przykładowa 2",
      "city": "Kraków",
      "zip": "31-002",
      "country": "PL"
    }
  }' | jq
  ```

  2. Manually grant admin role in SQLite

  ```bash
  sqlite3 ecommerce.db <<SQL
  UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';
  .exit
  SQL
  ```

  3. Login both users and capture tokens

  ```bash
  # Login as admin
  export ADMIN_TOKEN=$(
    curl -s -X POST http://localhost:8080/users/login \
      -H "Content-Type: application/json" \
      -d '{"email":"admin@example.com","password":"Password123!"}' \
      | jq -r .token
  )

  # Login as regular user
  export USER_TOKEN=$(
    curl -s -X POST http://localhost:8080/users/login \
      -H "Content-Type: application/json" \
      -d '{"email":"user@example.com","password":"Password123!"}' \
      | jq -r .token
  )

  echo "ADMIN_TOKEN:  $ADMIN_TOKEN"
  echo "USER_TOKEN:   $USER_TOKEN"
  ```

### 2. User Endpoints

1. `GET /users/{id}`
- Admin can fetch any user
- Regular user can fetch only their own ID (others → 403)
- No token → 401

```bash
# Assume `admin@example.com` has id=1, `user@example.com` has id=2

# Admin fetches user id=2  → 200
curl -s -X GET http://localhost:8080/users/2 \
  -H "Authorization: $ADMIN_TOKEN" \
  | jq

# Regular user fetches admin (id=1) → 403
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/users/1 \
  -H "Authorization: $USER_TOKEN"

# Regular user fetches own profile (id=2) → 200
curl -s -X GET http://localhost:8080/users/2 \
  -H "Authorization: $USER_TOKEN" \
  | jq
```

2. GET `/users`
- Admin → 200 + list of all users
- Regular user → 403
- No token → 401

```bash
# Admin fetches all users → 200
curl -s -X GET http://localhost:8080/users \
  -H "Authorization: $ADMIN_TOKEN" \
  | jq

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/users \
  -H "Authorization: $USER_TOKEN"

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/users
```

3. GET `/users/search?…`
- Admin can search (e.g. ?email=user@example.com)
- Regular user → 403
- No token → 401

```bash
# Admin searches by email → 200
curl -s -X GET 'http://localhost:8080/users/search?email=user@example.com' \
  -H "Authorization: $ADMIN_TOKEN" \
  | jq

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X GET 'http://localhost:8080/users/search?name=Anna' \
  -H "Authorization: $USER_TOKEN"

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X GET 'http://localhost:8080/users/search?name=Anna'
```

4. PUT `/users/{id}`
- Admin can update any user
- Regular user can update only their own profile
- No token → 401

```bash
# Admin updates user id=2  → 200
curl -s -X PUT http://localhost:8080/users/2 \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Anna-Edytowana", 
    "surname": "Kowalska-Nowa", 
    "email": "user@example.com" 
  }' | jq

# Regular user tries to update admin (id=1) → 403
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/users/1 \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Zly","surname":"Zmiana"}'

# Regular user updates own profile (id=2) → 200
curl -s -X PUT http://localhost:8080/users/2 \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Anna-Nowa",
    "surname": "Kowalska-Nowa",
    "email": "user@example.com"
  }' | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/users/2 \
  -H "Content-Type: application/json" \
  -d '{"name":"Anna","surname":"X"}'
```

5. DELETE `/users/{id}`
- Admin can delete any user
- Regular user can delete only their own account
- No token → 401

```bash
# Regular user tries to delete admin (id=1) → 403
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/users/1 \
  -H "Authorization: $USER_TOKEN"

# Regular user deletes own account (id=2) → 204
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/users/2 \
  -H "Authorization: $USER_TOKEN"

# Admin deletes a user (id=2) → 204
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/users/2 \
  -H "Authorization: $ADMIN_TOKEN"

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/users/1
```

### 3. Category Endpoints

1. Public GETs

```bash
# Get all categories → 200
curl -s -X GET http://localhost:8080/categories | jq

# Get category by id=1 → 200 or 404
curl -s -X GET http://localhost:8080/categories/1 | jq

# Get subcategories (id=1) → 200 or 404
curl -s -X GET http://localhost:8080/categories/1/subcategories | jq

# Search categories (e.g. name=Electronics) → 200
curl -s -X GET 'http://localhost:8080/categories/search?name=Electronics' | jq
```

2. Protected POST/PUT/DELETE (admin only)

```bash
# Admin creates a new category → 201
curl -s -X POST http://localhost:8080/categories \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Elektronika",
    "description": "Sprzęt elektroniczny"
  }' | jq

# Assume new category id = 1

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/categories \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Zabawki","description":"Dla dzieci"}'

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Zabawki","description":"Dla dzieci"}'

# Admin updates category id=1 → 200
curl -s -X PUT http://localhost:8080/categories/1 \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Elektronika RTV",
    "description": "Telewizory, głośniki, itp."
  }' | jq

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/categories/1 \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Błędna","description":"Zła"}'

# Admin deletes category id=1 → 204
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/categories/1 \
  -H "Authorization: $ADMIN_TOKEN"

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/categories/1 \
  -H "Authorization: $USER_TOKEN"

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/categories/1
```

### 4. Product Endpoints

1. Public GETs

```bash
# Get all products → 200
curl -s -X GET http://localhost:8080/products | jq

# Get product by id=1 → 200 or 404
curl -s -X GET http://localhost:8080/products/1 | jq

# Search products (e.g. category_id=1, price range) → 200
curl -s -X GET 'http://localhost:8080/products/search?category_id=1&price_min=100&price_max=500&name=phone' | jq
```

2. Protected POST/PUT/DELETE (admin only)

```bash
# Admin creates a product → 201
curl -s -X POST http://localhost:8080/products \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop ASUS",
    "description": "Laptop do pracy",
    "price": 3500,
    "category_id": 1,
    "stock": 10
  }' | jq

# Assume new product id = 1

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/products \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"ZłyProdukt","description":"Brak","price":1,"category_id":1,"stock":5}'

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"ZłyProdukt","description":"Brak","price":1,"category_id":1,"stock":5}'

# Admin updates product id=1 → 200
curl -s -X PUT http://localhost:8080/products/1 \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop ASUS ROG",
    "description": "Gamingowy laptop",
    "price": 4500,
    "category_id": 1,
    "stock": 7
  }' | jq

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/products/1 \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"X","description":"Y","price":1,"category_id":1,"stock":5}'

# Admin deletes product id=1 → 204
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/products/1 \
  -H "Authorization: $ADMIN_TOKEN"

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/products/1 \
  -H "Authorization: $USER_TOKEN"

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/products/1
```

### 5. Cart Endpoints (require JWT)

Endpoints under `/cart` always require a valid JWT.

1. GET `/cart`
- Regular user → their own cart (200)
- Admin → their own (empty) cart (200)
- No token → 401

```bash
# Regular user fetches their cart → 200
curl -s -X GET http://localhost:8080/cart \
  -H "Authorization: $USER_TOKEN" \
  | jq

# Admin fetches their cart → 200
curl -s -X GET http://localhost:8080/cart \
  -H "Authorization: $ADMIN_TOKEN" \
  | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/cart
```
2. POST `/cart/add`
- Add an item to the authenticated user’s cart. Requires JWT.

```bash
# Regular user adds product id=1 to cart → 200
curl -s -X POST http://localhost:8080/cart/add \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 2}' | jq

# Admin adds product id=1 to their cart → 200
curl -s -X POST http://localhost:8080/cart/add \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 1}' | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/cart/add \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "quantity": 1}'
```

3. PUT `/cart/item/{item_id}`
- Update quantity of a specific cart item.
- Owner or admin may update; others → 403.
- No token → 401.

```bash
# Assume cart item id=1 belongs to user id=2

# Regular user updates their item (id=1) → 200
curl -s -X PUT http://localhost:8080/cart/item/1 \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 5}' | jq

# Regular user tries to update admin’s item (id=2) → 403
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/cart/item/2 \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 2}'

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/cart/item/1 \
  -H "Content-Type: application/json" \
  -d '{"quantity": 2}'
```

4. DELETE `/cart/item/{item_id}`
- Remove a cart item.
- Owner or admin allowed; others → 403.
- No token → 401.

```bash
# Regular user removes their item (id=1) → 204
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/cart/item/1 \
  -H "Authorization: $USER_TOKEN"

# Regular user tries to remove admin’s item (id=2) → 403
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/cart/item/2 \
  -H "Authorization: $USER_TOKEN"

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/cart/item/1
```

5. DELETE `/cart/clear`
- Clear authenticated user’s cart.
- No token → 401.

```bash
# Regular user clears cart → 204
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/cart/clear \
  -H "Authorization: $USER_TOKEN"

# Admin clears cart (likely empty) → 204
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/cart/clear \
  -H "Authorization: $ADMIN_TOKEN"

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE http://localhost:8080/cart/clear
```

6. GET `/cart/search?…`
- Admin can filter any user’s carts (e.g. ?user_id=2).
- Regular user always sees only their own (ignores user_id).
- No token → 401.

```bash
# Admin searches all carts (user_id=2, total_max=1000) → 200
curl -s -X GET 'http://localhost:8080/cart/search?user_id=2&total_max=1000' \
  -H "Authorization: $ADMIN_TOKEN" | jq

# Regular user tries to filter by user_id=1 → sees only their own
curl -s -X GET 'http://localhost:8080/cart/search?user_id=1' \
  -H "Authorization: $USER_TOKEN" | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X GET 'http://localhost:8080/cart/search?user_id=2'
```

### 6. Order Endpoints (require JWT)
All `/orders` endpoints require a valid JWT.

Precondition
- Each user’s cart must contain items to create an order.
- Assume user id=2 has a cart with some items; admin id=1 similarly.

1. POST `/orders`
- Create an order from the authenticated user’s cart.
- No token → 401.

```bash
# Regular user creates an order → 201
curl -s -X POST http://localhost:8080/orders \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "CREDIT_CARD",
    "shipping_address_id": 1
  }' | jq

# Admin creates an order (from their cart) → 201
curl -s -X POST http://localhost:8080/orders \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "CREDIT_CARD",
    "shipping_address_id": 1
  }' | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "payment_method": "CREDIT_CARD",
    "shipping_address_id": 1
  }'
```

2. GET `/orders/{id}`
- Owner can fetch their own (200), other IDs → 403.
- Admin can fetch any (200).
- No token → 401.

```bash
# Regular user fetches own order (id=1) → 200
curl -s -X GET http://localhost:8080/orders/1 \
  -H "Authorization: $USER_TOKEN" | jq

# Regular user tries to fetch admin’s order (id=2) → 403
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/orders/2 \
  -H "Authorization: $USER_TOKEN"

# Admin fetches user’s order (id=1) → 200
curl -s -X GET http://localhost:8080/orders/1 \
  -H "Authorization: $ADMIN_TOKEN" | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/orders/1
```

3. GET `/orders/user`
- Returns a list of the authenticated user’s orders.
- Admin sees only their own via this endpoint.
- No token → 401.

```bash
# Regular user fetches their orders → 200
curl -s -X GET http://localhost:8080/orders/user \
  -H "Authorization: $USER_TOKEN" | jq

# Admin fetches their orders → 200
curl -s -X GET http://localhost:8080/orders/user \
  -H "Authorization: $ADMIN_TOKEN" | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/orders/user
```

4. GET `/orders`
- Admin only → all orders (200).
- Regular user → 403.
- No token → 401.

```bash
# Admin fetches all orders → 200
curl -s -X GET http://localhost:8080/orders \
  -H "Authorization: $ADMIN_TOKEN" | jq

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/orders \
  -H "Authorization: $USER_TOKEN"

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X GET http://localhost:8080/orders
```

5. PUT `/orders/{id}/status`
- Admin only can update status (200).
- Regular user → 403.
- No token → 401.

```bash
# Admin updates order status (id=1 to SHIPPED) → 200
curl -s -X PUT http://localhost:8080/orders/1/status \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "SHIPPED"}' | jq

# Regular user tries → 403
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/orders/1/status \
  -H "Authorization: $USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "CANCELLED"}'

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/orders/1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "SHIPPED"}'
```

6. PUT `/orders/{id}/cancel`
- Owner can cancel their own if still pending (200).
- Admin can cancel any (200).
- Other user → 403.
- No token → 401.

```bash
# Regular user cancels own order (id=1) → 200
curl -s -X PUT http://localhost:8080/orders/1/cancel \
  -H "Authorization: $USER_TOKEN" | jq

# Regular user tries to cancel admin’s order (id=2) → 403
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/orders/2/cancel \
  -H "Authorization: $USER_TOKEN"

# Admin cancels any order (e.g., id=1) → 200
curl -s -X PUT http://localhost:8080/orders/1/cancel \
  -H "Authorization: $ADMIN_TOKEN" | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X PUT http://localhost:8080/orders/1/cancel
```

7. GET `/orders/search?…`
- Admin can filter any (e.g. ?user_id=2&status=PENDING).
- Regular user sees only their own (ignores user_id).
- No token → 401.

```bash
# Admin searches orders by user_id=2 → 200
curl -s -X GET 'http://localhost:8080/orders/search?user_id=2&status=PENDING' \
  -H "Authorization: $ADMIN_TOKEN" | jq

# Regular user tries to search by user_id=1 → sees only their own
curl -s -X GET 'http://localhost:8080/orders/search?user_id=1' \
  -H "Authorization: $USER_TOKEN" | jq

# No token → 401
curl -s -o /dev/null -w "%{http_code}\n" -X GET 'http://localhost:8080/orders/search?user_id=2'
```

## Dependencies

- Echo – HTTP framework
- GORM – ORM for Go
- SQLite – Database engine
- Validator – Input validation

## Additional Notes

1. The database file (`ecommerce.db`) is created automatically on first run.
2. GORM’s AutoMigrate creates tables based on models.
3. Inventory (stock) is decreased when items are added to cart and orders are created.
4. Order status transitions are validated (e.g., shipped orders cannot be canceled).
EOF