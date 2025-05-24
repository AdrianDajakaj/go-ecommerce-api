# E-Commerce API Documentation

## Overview

This project implements a Clean Architecture approach for an e-commerce system written in Go. It uses the Echo framework for HTTP routing and GORM as the ORM layer. The application operates on an SQLite database.

---

## Requirements

- Go 1.20+
- SQLite3 (wbudowany, nie wymaga osobnej instalacji)

## Getting Started

### 1. Download repository 

```bash
git clone https://github.com/AdrianDajakaj/go-ecommerce-api.git
cd go-ecommerce-api
```

### 2. Install dependencies 

```bash
go mod tidy
```

### 3. Run server

```bash
go run cmd/server.go
```

Server will be available at: `http://localhost:8080`

## Basic configuration

Basic configuration is in `cmd/server.go`:

```go
dsn := "ecommerce.db" // database name
```

To change the `server port`:

```go
e.Start(":8080") // change used port
```

## Architecture

```
.
├── README.md
├── cmd
│   └── server.go
├── ecommerce.db
├── go.mod
├── go.sum
└── internal
    ├── domain
    │   ├── model
    │   │   ├── address.go
    │   │   ├── cart.go
    │   │   ├── cart_item.go
    │   │   ├── category.go
    │   │   ├── order.go
    │   │   ├── order_item.go
    │   │   ├── product.go
    │   │   └── user.go
    │   └── repository
    │       ├── address_repository.go
    │       ├── cart_item_repository.go
    │       ├── cart_repository.go
    │       ├── category_repository.go
    │       ├── order_item_repository.go
    │       ├── order_repository.go
    │       ├── product_repository.go
    │       └── user_repository.go
    ├── infrastructure
    │   └── persistence
    │       ├── repository
    │       │   ├── address_repository.go
    │       │   ├── cart_item_repository.go
    │       │   ├── cart_repository.go
    │       │   ├── category_repository.go
    │       │   ├── order_item_repository.go
    │       │   ├── order_repository.go
    │       │   ├── product_repository.go
    │       │   └── user_repository.go
    │       ├── scope
    │       │   ├── cart_scope.go
    │       │   ├── category_scope.go
    │       │   ├── order_scope.go
    │       │   ├── product_scope.go
    │       │   └── user_scope.go
    │       └── sqlite
    │           └── gorm_db.go
    ├── interface
    │   └── http
    │       ├── handler
    │       │   ├── cart_handler.go
    │       │   ├── category_handler.go
    │       │   ├── order_handler.go
    │       │   ├── product_handler.go
    │       │   └── user_handler.go
    │       └── router.go
    └── usecase
        ├── cart_usecase.go
        ├── category_usecase.go
        ├── order_usecase.go
        ├── product_usecase.go
        └── user_usecase.go
```

---

## Features

* Clean Architecture structure
* CRUD operations for Users, Products, Categories, Carts, Orders
* Cart management (Add, Update, Remove items)
* Order placement and status tracking
* GORM scopes for dynamic filtering via query parameters

---

## Base URL

```
http://localhost:8080
```

---

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
  "name": "Electronics"
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

---

## Endpoint Patterns

Each model exposes standardized RESTful endpoints:

### Users

* `POST    /users/register` – Register new user
* `POST    /users/login` – Authenticate user
* `GET     /users` – Get all users
* `GET     /users/{id}` – Get user by ID
* `GET     /users/search?...` – Search users using query parameters (with scopes)
* `PUT     /users/{id}` – Update user
* `DELETE  /users/{id}` – Delete user

### Categories

* `POST    /categories` – Create category
* `GET     /categories` – Get all categories
* `GET     /categories/{id}` – Get category by ID
* `GET     /categories/search?...` – Search categories using query parameters (with scopes)
* `PUT     /categories/{id}` – Update category
* `DELETE  /categories/{id}` – Delete category

### Products

* `POST    /products` – Create product
* `GET     /products` – Get all products
* `GET     /products/{id}` – Get product by ID
* `GET     /products/search?...` – Search products using query parameters (with scopes)
* `PUT     /products/{id}` – Update product
* `DELETE  /products/{id}` – Delete product

### Carts

* `GET     /cart/{user_id}` – Get cart by user
* `GET     /cart/search?...` – Search carts with filters (scopes)
* `POST    /cart/{user_id}/add` – Add product to cart
* `PUT     /cart/item/{item_id}` – Update cart item quantity
* `DELETE  /cart/item/{item_id}` – Remove item
* `DELETE  /cart/{user_id}/clear` – Clear entire cart

### Orders

* `POST    /users/{user_id}/orders` – Create order from cart
* `GET     /orders` – Get all orders
* `GET     /orders/{id}` – Get order by ID
* `GET     /users/{user_id}/orders` – Get all User's orders
* `GET     /orders/search?...` – Search orders (via scopes)
* `PUT     /orders/{id}/status` – Update status
* `PUT     /orders/{id}/cancel` – Cancel order

---

## Scopes (Filtering with Query Parameters)

### User Scopes

* `email=john@example.com`
* `name=John`
* `surname=Doe`
* `country=USA`
* `city=New York`

### Product Scopes

* `name=phone`
* `category_id=1`
* `is_active=true`
* `price_min=100&price_max=500`

### Category Scopes

* `name=electronics`
* `created_after=2023-01-01T00:00:00Z`
* `created_before=2023-12-31T23:59:59Z`
* `min_products=3`
* `with_products=true`

### Order Scopes

* `user_id=1`
* `status=PAID`
* `created_after=2023-01-01T00:00:00Z`
* `total_min=100&total_max=1000`

### Cart Scopes

* `user_id=1`
* `total_min=100&total_max=500`
* `created_before=2024-01-01T00:00:00Z`

---

## cURL Examples With Responses

### 🧑 User – cURL Examples with Responses

1. 🔐 Register a New User

- **Request**

```bash
curl -X POST -i http://localhost:8080/users/register \
  -H "Content-Type: application/json" \
  -d '{
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
  }'
```

- **Response**

<pre> <code>HTTP/1.1 201 Created
Content-Type: application/json
Date: Sat, 24 May 2025 21:02:58 GMT
Content-Length: 390 </code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:02:58.197805299+02:00",
  "updated_at": "2025-05-24T23:02:58.197805299+02:00",
  "email": "john@example.com",
  "name": "John",
  "surname": "Doe",
  "address_id": 1,
  "address": {
    "id": 1,
    "created_at": "2025-05-24T23:02:58.139428936+02:00",
    "updated_at": "2025-05-24T23:02:58.139428936+02:00",
    "country": "USA",
    "city": "New York",
    "postcode": "10001",
    "street": "Broadway",
    "number": "1"
  }
}
```

2. 🔓 Login

- **Request**

```bash
curl -X POST -i http://localhost:8080/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "secure123"
  }'
```

- **Response**

<pre> <code>HTTP/1.1 200 OK Content-Type: application/json Date: Sat, 24 May 2025 21:03:54 GMT Content-Length: 390 </code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:02:58.197805299+02:00",
  "updated_at": "2025-05-24T23:02:58.197805299+02:00",
  "email": "john@example.com",
  "name": "John",
  "surname": "Doe",
  "address_id": 1,
  "address": {
    "id": 1,
    "created_at": "2025-05-24T23:02:58.139428936+02:00",
    "updated_at": "2025-05-24T23:02:58.139428936+02:00",
    "country": "USA",
    "city": "New York",
    "postcode": "10001",
    "street": "Broadway",
    "number": "1"
  }
}
```

3. 📄 Get All Users

- **Request**

```bash
curl -i http://localhost:8080/users
```

- **Response**

<pre> <code>HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:07:31 GMT
Content-Length: 392 </code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:02:58.197805299+02:00",
    "updated_at": "2025-05-24T23:02:58.197805299+02:00",
    "email": "john@example.com",
    "name": "John",
    "surname": "Doe",
    "address_id": 1,
    "address": {
      "id": 1,
      "created_at": "2025-05-24T23:02:58.139428936+02:00",
      "updated_at": "2025-05-24T23:02:58.139428936+02:00",
      "country": "USA",
      "city": "New York",
      "postcode": "10001",
      "street": "Broadway",
      "number": "1"
    }
  }
]
```

4. 🔍 Get User by ID

- **Request**

```bash
curl http://localhost:8080/users/1
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:08:29 GMT
Content-Length: 390
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:02:58.197805299+02:00",
  "updated_at": "2025-05-24T23:02:58.197805299+02:00",
  "email": "john@example.com",
  "name": "John",
  "surname": "Doe",
  "address_id": 1,
  "address": {
    "id": 1,
    "created_at": "2025-05-24T23:02:58.139428936+02:00",
    "updated_at": "2025-05-24T23:02:58.139428936+02:00",
    "country": "USA",
    "city": "New York",
    "postcode": "10001",
    "street": "Broadway",
    "number": "1"
  }
}
```
5. 🧭 Search Users (With Scopes)

- **Request**

```bash
curl -i "http://localhost:8080/users/search?country=USA&city=New%20York"
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:10:27 GMT
Content-Length: 392
</code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:02:58.197805299+02:00",
    "updated_at": "2025-05-24T23:02:58.197805299+02:00",
    "email": "john@example.com",
    "name": "John",
    "surname": "Doe",
    "address_id": 1,
    "address": {
      "id": 1,
      "created_at": "2025-05-24T23:02:58.139428936+02:00",
      "updated_at": "2025-05-24T23:02:58.139428936+02:00",
      "country": "USA",
      "city": "New York",
      "postcode": "10001",
      "street": "Broadway",
      "number": "1"
    }
  }
]
```

6. 📝 Update User

- **Request**

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "name": "Johnathan",
    "surname": "Doe"
  }'

```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:12:49 GMT
Content-Length: 328
</code> </pre>

```json
{
  "id": 1,
  "created_at": "0001-01-01T00:00:00Z",
  "updated_at": "2025-05-24T23:11:32.979252277+02:00",
  "email": "john.doe@example.com",
  "name": "Johnathan",
  "surname": "Doe",
  "address_id": 0,
  "address": {
    "id": 0,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z",
    "country": "",
    "city": "",
    "postcode": "",
    "street": "",
    "number": ""
  }
}
```

7. ❌ Delete User

- **Request**

```bash
curl -X -i DELETE http://localhost:8080/users/1
```

- **Response**

<pre> <code>
HTTP/1.1 204 No Content
</code> </pre>

### 📂 Category – cURL Examples with Responses

1. ➕ Create Category

- **Request**

```bash
curl -X POST -i http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics"
  }'
```

- **Response**

<pre> <code>
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sat, 24 May 2025 21:13:39 GMT
Content-Length: 132
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:13:39.187374675+02:00",
  "updated_at": "2025-05-24T23:13:39.187374675+02:00",
  "name": "Electronics"
}
```

2. 📄 Get All Categories

- **Request**

```bash
curl -i http://localhost:8080/categories
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:15:03 GMT
Content-Length: 134
</code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:13:39.187374675+02:00",
    "updated_at": "2025-05-24T23:13:39.187374675+02:00",
    "name": "Electronics"
  }
]
```

3. 🔍 Get Category by ID

- **Request**

```bash
curl -i http://localhost:8080/categories/1
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:15:55 GMT
Content-Length: 132
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:13:39.187374675+02:00",
  "updated_at": "2025-05-24T23:13:39.187374675+02:00",
  "name": "Electronics"
}
```

4. 🧭 Search Categories (Scopes)

- **Request**

```bash
curl -i "http://localhost:8080/categories/search?name=Electronics&with_products=true"
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:16:38 GMT
Content-Length: 134
</code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:13:39.187374675+02:00",
    "updated_at": "2025-05-24T23:13:39.187374675+02:00",
    "name": "Electronics"
  }
]
```

5. 📝 Update Category

- **Request**

```bash
curl -X PUT http://localhost:8080/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Consumer Electronics"
  }'
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:17:36 GMT
Content-Length: 126
</code> </pre>

```json
{
  "id": 1,
  "created_at": "0001-01-01T00:00:00Z",
  "updated_at": "2025-05-24T23:17:36.640609459+02:00",
  "name": "Consumer Electronics"
}
```

6. ❌ Delete Category

- **Request**

```bash
curl -X DELETE http://localhost:8080/categories/1
```

- **Response**

<pre> <code>
HTTP/1.1 204 No Content
</code> </pre>

### 🛍️ Product – cURL Examples with Responses

1. ➕ Create Product

- **Request**

```bash
curl -X POST -i http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Phone",
    "description": "Smartphone with AMOLED display",
    "price": 299.99,
    "stock": 100,
    "is_active": true,
    "category_id": 1,
    "images": [
      { "url": "https://example.com/images/phone-front.jpg" },
      { "url": "https://example.com/images/phone-back.jpg" }
    ]
  }'
```

- **Response**

<pre> <code>
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sat, 24 May 2025 21:18:58 GMT
Content-Length: 734
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:18:58.601529769+02:00",
  "updated_at": "2025-05-24T23:18:58.601529769+02:00",
  "name": "Phone",
  "description": "Smartphone with AMOLED display",
  "price": 299.99,
  "stock": 100,
  "is_active": true,
  "category_id": 1,
  "category": {
    "id": 1,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "2025-05-24T23:17:36.640609459+02:00",
    "name": "Consumer Electronics"
  },
  "images": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:18:58.602234769+02:00",
      "updated_at": "2025-05-24T23:18:58.602234769+02:00",
      "url": "https://example.com/images/phone-front.jpg",
      "product_id": 1
    },
    {
      "id": 2,
      "created_at": "2025-05-24T23:18:58.602234769+02:00",
      "updated_at": "2025-05-24T23:18:58.602234769+02:00",
      "url": "https://example.com/images/phone-back.jpg",
      "product_id": 1
    }
  ]
}
```

2. 📄 Get All Products

- **Request**

```bash
curl -i http://localhost:8080/products
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:20:06 GMT
Content-Length: 736
</code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:18:58.601529769+02:00",
    "updated_at": "2025-05-24T23:18:58.601529769+02:00",
    "name": "Phone",
    "description": "Smartphone with AMOLED display",
    "price": 299.99,
    "stock": 100,
    "is_active": true,
    "category_id": 1,
    "category": {
      "id": 1,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "2025-05-24T23:17:36.640609459+02:00",
      "name": "Consumer Electronics"
    },
    "images": [
      {
        "id": 1,
        "created_at": "2025-05-24T23:18:58.602234769+02:00",
        "updated_at": "2025-05-24T23:18:58.602234769+02:00",
        "url": "https://example.com/images/phone-front.jpg",
        "product_id": 1
      },
      {
        "id": 2,
        "created_at": "2025-05-24T23:18:58.602234769+02:00",
        "updated_at": "2025-05-24T23:18:58.602234769+02:00",
        "url": "https://example.com/images/phone-back.jpg",
        "product_id": 1
      }
    ]
  }
]
```

3. 🔍 Get Product by ID

- **Request**

```bash
curl -i http://localhost:8080/products/1
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:21:01 GMT
Content-Length: 734
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:18:58.601529769+02:00",
  "updated_at": "2025-05-24T23:18:58.601529769+02:00",
  "name": "Phone",
  "description": "Smartphone with AMOLED display",
  "price": 299.99,
  "stock": 100,
  "is_active": true,
  "category_id": 1,
  "category": {
    "id": 1,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "2025-05-24T23:17:36.640609459+02:00",
    "name": "Consumer Electronics"
  },
  "images": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:18:58.602234769+02:00",
      "updated_at": "2025-05-24T23:18:58.602234769+02:00",
      "url": "https://example.com/images/phone-front.jpg",
      "product_id": 1
    },
    {
      "id": 2,
      "created_at": "2025-05-24T23:18:58.602234769+02:00",
      "updated_at": "2025-05-24T23:18:58.602234769+02:00",
      "url": "https://example.com/images/phone-back.jpg",
      "product_id": 1
    }
  ]
}
```

4. 🧭 Search Products (Scopes)

- **Request**

```bash
curl -i "http://localhost:8080/products/search?category_id=1&is_active=true&price_min=100&price_max=500&name=phone"
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:21:45 GMT
Content-Length: 736
</code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:18:58.601529769+02:00",
    "updated_at": "2025-05-24T23:18:58.601529769+02:00",
    "name": "Phone",
    "description": "Smartphone with AMOLED display",
    "price": 299.99,
    "stock": 100,
    "is_active": true,
    "category_id": 1,
    "category": {
      "id": 1,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "2025-05-24T23:17:36.640609459+02:00",
      "name": "Consumer Electronics"
    },
    "images": [
      {
        "id": 1,
        "created_at": "2025-05-24T23:18:58.602234769+02:00",
        "updated_at": "2025-05-24T23:18:58.602234769+02:00",
        "url": "https://example.com/images/phone-front.jpg",
        "product_id": 1
      },
      {
        "id": 2,
        "created_at": "2025-05-24T23:18:58.602234769+02:00",
        "updated_at": "2025-05-24T23:18:58.602234769+02:00",
        "url": "https://example.com/images/phone-back.jpg",
        "product_id": 1
      }
    ]
  }
]
```

5. 📝 Update Product

- **Request**

```bash
curl -X PUT -i http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Phone",
    "description": "Now with 5G support",
    "price": 349.99,
    "stock": 80,
    "is_active": true,
    "category_id": 1
  }'
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:22:36 GMT
Content-Length: 715
</code> </pre>

```json
{
  "id": 1,
  "created_at": "0001-01-01T00:00:00Z",
  "updated_at": "2025-05-24T23:22:36.440839439+02:00",
  "name": "Updated Phone",
  "description": "Now with 5G support",
  "price": 349.99,
  "stock": 80,
  "is_active": true,
  "category_id": 1,
  "category": {
    "id": 1,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "2025-05-24T23:17:36.640609459+02:00",
    "name": "Consumer Electronics"
  },
  "images": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:18:58.602234769+02:00",
      "updated_at": "2025-05-24T23:18:58.602234769+02:00",
      "url": "https://example.com/images/phone-front.jpg",
      "product_id": 1
    },
    {
      "id": 2,
      "created_at": "2025-05-24T23:18:58.602234769+02:00",
      "updated_at": "2025-05-24T23:18:58.602234769+02:00",
      "url": "https://example.com/images/phone-back.jpg",
      "product_id": 1
    }
  ]
}
```

6. ❌ Delete Product

- **Request**

```bash
curl -X DELETE -i http://localhost:8080/products/1
```

- **Response**

<pre> <code>
HTTP/1.1 204 No Content
</code> </pre>

### 🛒 Cart – cURL Examples with Responses

1. ➕ Add Product to User's Cart

- **Request**

```bash
curl -X POST -i http://localhost:8080/cart/1/add \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:24:31 GMT
Content-Length: 675
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:24:31.046637875+02:00",
  "updated_at": "2025-05-24T23:24:31.0587211+02:00",
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:24:31.052739562+02:00",
      "updated_at": "2025-05-24T23:24:31.052739562+02:00",
      "cart_id": 1,
      "product_id": 1,
      "product": {
        "id": 1,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "2025-05-24T23:22:36.440839439+02:00",
        "name": "Updated Phone",
        "description": "Now with 5G support",
        "price": 349.99,
        "stock": 80,
        "is_active": true,
        "category_id": 1,
        "category": {
          "id": 0,
          "created_at": "0001-01-01T00:00:00Z",
          "updated_at": "0001-01-01T00:00:00Z",
          "name": ""
        },
        "images": null
      },
      "quantity": 2,
      "unit_price": 349.99,
      "subtotal": 699.98
    }
  ],
  "total": 699.98
}
```

2. 📄 Get Cart by User ID

- **Request**

```bash
curl -i http://localhost:8080/cart/1
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:25:34 GMT
Content-Length: 675
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:24:31.046637875+02:00",
  "updated_at": "2025-05-24T23:24:31.0587211+02:00",
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:24:31.052739562+02:00",
      "updated_at": "2025-05-24T23:24:31.052739562+02:00",
      "cart_id": 1,
      "product_id": 1,
      "product": {
        "id": 1,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "2025-05-24T23:22:36.440839439+02:00",
        "name": "Updated Phone",
        "description": "Now with 5G support",
        "price": 349.99,
        "stock": 80,
        "is_active": true,
        "category_id": 1,
        "category": {
          "id": 0,
          "created_at": "0001-01-01T00:00:00Z",
          "updated_at": "0001-01-01T00:00:00Z",
          "name": ""
        },
        "images": null
      },
      "quantity": 2,
      "unit_price": 349.99,
      "subtotal": 699.98
    }
  ],
  "total": 699.98
}
```

3. 🔄 Update Cart Item Quantity

- **Request**

```bash
curl -X PUT -i http://localhost:8080/cart/item/1 \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 3
  }'
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:26:54 GMT
Content-Length: 678
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:24:31.046637875+02:00",
  "updated_at": "2025-05-24T23:26:54.97415734+02:00",
  "user_id": 1,
  "items": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:24:31.052739562+02:00",
      "updated_at": "2025-05-24T23:26:54.967645549+02:00",
      "cart_id": 1,
      "product_id": 1,
      "product": {
        "id": 1,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "2025-05-24T23:22:36.440839439+02:00",
        "name": "Updated Phone",
        "description": "Now with 5G support",
        "price": 349.99,
        "stock": 80,
        "is_active": true,
        "category_id": 1,
        "category": {
          "id": 0,
          "created_at": "0001-01-01T00:00:00Z",
          "updated_at": "0001-01-01T00:00:00Z",
          "name": ""
        },
        "images": null
      },
      "quantity": 3,
      "unit_price": 349.99,
      "subtotal": 1049.97
    }
  ],
  "total": 1049.97
}
```

4. ❌ Remove Cart Item

- **Request**

```bash
curl -X DELETE -i http://localhost:8080/carts/item/1
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:47:52 GMT
Content-Length: 133
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:24:31.046637875+02:00",
  "updated_at": "2025-05-24T23:47:52.626251683+02:00",
  "user_id": 1,
  "total": 0
}
```
5. 🧹 Clear Entire Cart

- **Request**

```bash
curl -X DELETE -i http://localhost:8080/cart/1/clear
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:48:57 GMT
Content-Length: 133
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:24:31.046637875+02:00",
  "updated_at": "2025-05-24T23:48:57.263471076+02:00",
  "user_id": 1,
  "total": 0
}
```

6. 🧭 Search Carts (Scopes)

- **Request**

```bash
curl -i "http://localhost:8080/cart/search?user_id=1&total_max=1100"
```

- **Response**


<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:29:24 GMT
Content-Length: 680
</code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:24:31.046637875+02:00",
    "updated_at": "2025-05-24T23:26:54.97415734+02:00",
    "user_id": 1,
    "items": [
      {
        "id": 1,
        "created_at": "2025-05-24T23:24:31.052739562+02:00",
        "updated_at": "2025-05-24T23:26:54.967645549+02:00",
        "cart_id": 1,
        "product_id": 1,
        "product": {
          "id": 1,
          "created_at": "0001-01-01T00:00:00Z",
          "updated_at": "2025-05-24T23:22:36.440839439+02:00",
          "name": "Updated Phone",
          "description": "Now with 5G support",
          "price": 349.99,
          "stock": 80,
          "is_active": true,
          "category_id": 1,
          "category": {
            "id": 0,
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z",
            "name": ""
          },
          "images": null
        },
        "quantity": 3,
        "unit_price": 349.99,
        "subtotal": 1049.97
      }
    ],
    "total": 1049.97
  }
]
```

### 📦 Order – cURL Examples with Responses

1. ➕ Create Order from Cart

- **Request**

```bash
curl -X POST -i http://localhost:8080/users/1/orders   -H "Content-Type: application/json"
  -d '{
    "payment_method": "CARD",
    "shipping_address_id": 1
  }'
```

- **Response**
<pre> <code>
HTTP/1.1 201 Created
Content-Type: application/json
Date: Sat, 24 May 2025 21:31:58 GMT
Content-Length: 880
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:31:58.794202304+02:00",
  "updated_at": "2025-05-24T23:31:58.794202304+02:00",
  "user_id": 1,
  "user": {
    "id": 0,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z",
    "email": "",
    "name": "",
    "surname": "",
    "address_id": 0,
    "address": {
      "id": 0,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z",
      "country": "",
      "city": "",
      "postcode": "",
      "street": "",
      "number": ""
    }
  },
  "status": "PENDING",
  "shipping_address_id": 1,
  "shipping_address": {
    "id": 0,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z",
    "country": "",
    "city": "",
    "postcode": "",
    "street": "",
    "number": ""
  },
  "payment_method": "CARD",
  "items": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:31:58.794656522+02:00",
      "updated_at": "2025-05-24T23:31:58.794656522+02:00",
      "order_id": 1,
      "product_id": 1,
      "name": "Updated Phone",
      "unit_price": 349.99,
      "quantity": 3,
      "subtotal": 1049.97
    }
  ],
  "total": 1049.97
}
```

2. 📄 Get All Orders

- **Request**

```bash
curl -i http://localhost:8080/orders
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:35:06 GMT
Content-Length: 983
</code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:31:58.794202304+02:00",
    "updated_at": "2025-05-24T23:31:58.794202304+02:00",
    "user_id": 1,
    "user": {
      "id": 1,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "2025-05-24T23:12:49.09327422+02:00",
      "email": "john.doe@example.com",
      "name": "Johnathan",
      "surname": "Doe",
      "address_id": 0,
      "address": {
        "id": 0,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z",
        "country": "",
        "city": "",
        "postcode": "",
        "street": "",
        "number": ""
      }
    },
    "status": "PENDING",
    "shipping_address_id": 1,
    "shipping_address": {
      "id": 1,
      "created_at": "2025-05-24T23:02:58.139428936+02:00",
      "updated_at": "2025-05-24T23:02:58.139428936+02:00",
      "country": "USA",
      "city": "New York",
      "postcode": "10001",
      "street": "Broadway",
      "number": "1"
    },
    "payment_method": "CARD",
    "items": [
      {
        "id": 1,
        "created_at": "2025-05-24T23:31:58.794656522+02:00",
        "updated_at": "2025-05-24T23:31:58.794656522+02:00",
        "order_id": 1,
        "product_id": 1,
        "name": "Updated Phone",
        "unit_price": 349.99,
        "quantity": 3,
        "subtotal": 1049.97
      }
    ],
    "total": 1049.97
  }
]
```

3. 🔍 Get Order by ID

- **Request**

```bash
curl -i http://localhost:8080/orders/1
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:36:31 GMT
Content-Length: 981
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:31:58.794202304+02:00",
  "updated_at": "2025-05-24T23:31:58.794202304+02:00",
  "user_id": 1,
  "user": {
    "id": 1,
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "2025-05-24T23:12:49.09327422+02:00",
    "email": "john.doe@example.com",
    "name": "Johnathan",
    "surname": "Doe",
    "address_id": 0,
    "address": {
      "id": 0,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z",
      "country": "",
      "city": "",
      "postcode": "",
      "street": "",
      "number": ""
    }
  },
  "status": "PENDING",
  "shipping_address_id": 1,
  "shipping_address": {
    "id": 1,
    "created_at": "2025-05-24T23:02:58.139428936+02:00",
    "updated_at": "2025-05-24T23:02:58.139428936+02:00",
    "country": "USA",
    "city": "New York",
    "postcode": "10001",
    "street": "Broadway",
    "number": "1"
  },
  "payment_method": "CARD",
  "items": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:31:58.794656522+02:00",
      "updated_at": "2025-05-24T23:31:58.794656522+02:00",
      "order_id": 1,
      "product_id": 1,
      "name": "Updated Phone",
      "unit_price": 349.99,
      "quantity": 3,
      "subtotal": 1049.97
    }
  ],
  "total": 1049.97
}
```

4. 🧭 Search Orders (Scopes)

- **Request**

```bash
curl -i "http://localhost:8080/orders/search?user_id=1&status=PENDING&total_min=100"
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:37:57 GMT
Content-Length: 983
</code> </pre>

```json
[
  {
    "id": 1,
    "created_at": "2025-05-24T23:31:58.794202304+02:00",
    "updated_at": "2025-05-24T23:31:58.794202304+02:00",
    "user_id": 1,
    "user": {
      "id": 1,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "2025-05-24T23:12:49.09327422+02:00",
      "email": "john.doe@example.com",
      "name": "Johnathan",
      "surname": "Doe",
      "address_id": 0,
      "address": {
        "id": 0,
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z",
        "country": "",
        "city": "",
        "postcode": "",
        "street": "",
        "number": ""
      }
    },
    "status": "PENDING",
    "shipping_address_id": 1,
    "shipping_address": {
      "id": 1,
      "created_at": "2025-05-24T23:02:58.139428936+02:00",
      "updated_at": "2025-05-24T23:02:58.139428936+02:00",
      "country": "USA",
      "city": "New York",
      "postcode": "10001",
      "street": "Broadway",
      "number": "1"
    },
    "payment_method": "CARD",
    "items": [
      {
        "id": 1,
        "created_at": "2025-05-24T23:31:58.794656522+02:00",
        "updated_at": "2025-05-24T23:31:58.794656522+02:00",
        "order_id": 1,
        "product_id": 1,
        "name": "Updated Phone",
        "unit_price": 349.99,
        "quantity": 3,
        "subtotal": 1049.97
      }
    ],
    "total": 1049.97
  }
]
```

5. 🔄 Update Order Status

- **Request**

```bash
curl -X PUT -i http://localhost:8080/orders/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "PAID"
  }'
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:39:44 GMT
Content-Length: 1041
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:31:58.794202304+02:00",
  "updated_at": "2025-05-24T23:39:44.868815325+02:00",
  "user_id": 1,
  "user": {
    "id": 1,
    "created_at": "2025-05-24T23:39:44.868582325+02:00",
    "updated_at": "2025-05-24T23:12:49.09327422+02:00",
    "email": "john.doe@example.com",
    "name": "Johnathan",
    "surname": "Doe",
    "address_id": 0,
    "address": {
      "id": 0,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z",
      "country": "",
      "city": "",
      "postcode": "",
      "street": "",
      "number": ""
    }
  },
  "status": "PAID",
  "paid_at": "2025-05-24T23:39:44.868546159+02:00",
  "shipping_address_id": 1,
  "shipping_address": {
    "id": 1,
    "created_at": "2025-05-24T23:02:58.139428936+02:00",
    "updated_at": "2025-05-24T23:02:58.139428936+02:00",
    "country": "USA",
    "city": "New York",
    "postcode": "10001",
    "street": "Broadway",
    "number": "1"
  },
  "payment_method": "CARD",
  "items": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:31:58.794656522+02:00",
      "updated_at": "2025-05-24T23:31:58.794656522+02:00",
      "order_id": 1,
      "product_id": 1,
      "name": "Updated Phone",
      "unit_price": 349.99,
      "quantity": 3,
      "subtotal": 1049.97
    }
  ],
  "total": 1049.97
}
```

6. ❌ Cancel Order

- **Request**

```bash
curl -X PUT -i http://localhost:8080/orders/1/cancel
```

- **Response**

<pre> <code>
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sat, 24 May 2025 21:41:51 GMT
Content-Length: 1098
</code> </pre>

```json
{
  "id": 1,
  "created_at": "2025-05-24T23:31:58.794202304+02:00",
  "updated_at": "2025-05-24T23:41:51.000742291+02:00",
  "user_id": 1,
  "user": {
    "id": 1,
    "created_at": "2025-05-24T23:41:51.00048975+02:00",
    "updated_at": "2025-05-24T23:12:49.09327422+02:00",
    "email": "john.doe@example.com",
    "name": "Johnathan",
    "surname": "Doe",
    "address_id": 0,
    "address": {
      "id": 0,
      "created_at": "0001-01-01T00:00:00Z",
      "updated_at": "0001-01-01T00:00:00Z",
      "country": "",
      "city": "",
      "postcode": "",
      "street": "",
      "number": ""
    }
  },
  "status": "CANCELLED",
  "paid_at": "2025-05-24T23:39:44.868546159+02:00",
  "cancelled_at": "2025-05-24T23:41:51.000393251+02:00",
  "shipping_address_id": 1,
  "shipping_address": {
    "id": 1,
    "created_at": "2025-05-24T23:02:58.139428936+02:00",
    "updated_at": "2025-05-24T23:02:58.139428936+02:00",
    "country": "USA",
    "city": "New York",
    "postcode": "10001",
    "street": "Broadway",
    "number": "1"
  },
  "payment_method": "CARD",
  "items": [
    {
      "id": 1,
      "created_at": "2025-05-24T23:31:58.794656522+02:00",
      "updated_at": "2025-05-24T23:31:58.794656522+02:00",
      "order_id": 1,
      "product_id": 1,
      "name": "Updated Phone",
      "unit_price": 349.99,
      "quantity": 3,
      "subtotal": 1049.97
    }
  ],
  "total": 1049.97
}
```

## Dependencies

Main packages:

- Echo - framework HTTP

- GORM - ORM for Go

- SQLite - database engine

- Validator - data validation

## Additional Notes

1. The database is automatically created on first run in the file `ecommerce.db`.

2. The database schema is generated using GORM AutoMigrate.

3. All operations on products/cart update the inventory in real time.

4. Order statuses are validated for consistency (e.g., an order that has already been shipped cannot be canceled).