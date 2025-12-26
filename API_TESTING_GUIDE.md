# Bookstore Order Processing System - API Testing Guide

## 🚀 Quick Start

### Prerequisites
1. **Server Running**: Make sure the API server is running on `http://localhost:3001`
2. **Database Setup**: Ensure PostgreSQL is running and migrations have been applied
3. **Postman Installed**: Download and install Postman from [postman.com](https://www.postman.com/)

### Starting the Server
```bash
cd /home/mohammed/Desktop/DB\ project/bookstore-project
export JWT_SECRET=your-secret-key-here
go run cmd/api/main.go
```

## 📚 Available Endpoints

### 🔐 Authentication Endpoints

#### 1. Customer Login
- **URL**: `POST /api/auth/login`
- **Body**:
```json
{
    "username": "alice",
    "password": "hashed_password_1",
    "role": "customer"
}
```
- **Response**:
```json
{
    "username": "alice",
    "role": "customer",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 2. Admin Login
- **URL**: `POST /api/auth/login`
- **Body**:
```json
{
    "username": "admin1",
    "password": "admin_pass_1",
    "role": "admin"
}
```
- **Response**:
```json
{
    "username": "admin1",
    "role": "admin",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 3. Customer Signup
- **URL**: `POST /api/auth/signup`
- **Body**:
```json
{
    "username": "newcustomer",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "shipping_address": "123 Main St, City, State 12345"
}
```

### 📖 Public Book Search Endpoints (No Authentication Required)

#### 1. Get All Books
- **URL**: `GET /api/books`
- **Example**: `curl http://localhost:3001/api/books`

#### 2. Get Book by ISBN
- **URL**: `GET /api/books/{isbn}`
- **Example**: `curl http://localhost:3001/api/books/9780132350884`

#### 3. Get Books by Category
- **URL**: `GET /api/books/category/{category}`
- **Valid Categories**: `Science`, `Art`, `Religion`, `History`, `Geography`
- **Example**: `curl http://localhost:3001/api/books/category/Science`

#### 4. Get Books by Author
- **URL**: `GET /api/books/author/{author_name}`
- **Example**: `curl "http://localhost:3001/api/books/author/Robert%20C.%20Martin"`

#### 5. Get Books by Title
- **URL**: `GET /api/books/title/{title}`
- **Example**: `curl "http://localhost:3001/api/books/title/Clean%20Code"`

#### 6. Get Books by Publication Year
- **URL**: `GET /api/books/pubyear/{year}`
- **Example**: `curl http://localhost:3001/api/books/pubyear/2008`

#### 7. Get Books by Publisher
- **URL**: `GET /api/books/publisher/{publisher_id}`
- **Example**: `curl http://localhost:3001/api/books/publisher/1`

### 🔒 Admin-Only Endpoints (Require Admin Token)

#### Book Management

##### 1. Add New Book
- **URL**: `POST /api/admin/books`
- **Headers**: `Authorization: Bearer {admin_token}`
- **Body**:
```json
{
    "isbn": "9780123456789",
    "title": "Test Book for API",
    "publication_year": 2024,
    "selling_price": 29.99,
    "category": "Science",
    "author_name": "Test Author",
    "stock_quantity": 15,
    "threshold": 5,
    "publisher_id": 1
}
```

##### 2. Update Existing Book
- **URL**: `PUT /api/admin/books/{isbn}`
- **Headers**: `Authorization: Bearer {admin_token}`
- **Body**:
```json
{
    "title": "Updated Test Book",
    "selling_price": 34.99,
    "category": "Art",
    "publication_year": 2024
}
```

##### 3. Update Book Stock Quantity
- **URL**: `PATCH /api/admin/books/{isbn}/stock`
- **Headers**: `Authorization: Bearer {admin_token}`
- **Body**:
```json
{
    "quantity": 12
}
```

#### Publisher Order Management

##### 1. Get All Publisher Orders
- **URL**: `GET /api/admin/publisher-orders`
- **Headers**: `Authorization: Bearer {admin_token}`

##### 2. Get Pending Publisher Orders
- **URL**: `GET /api/admin/publisher-orders/pending`
- **Headers**: `Authorization: Bearer {admin_token}`

##### 3. Get Publisher Order by ID
- **URL**: `GET /api/admin/publisher-orders/{order_id}`
- **Headers**: `Authorization: Bearer {admin_token}`

##### 4. Place Publisher Order
- **URL**: `POST /api/admin/publisher-orders`
- **Headers**: `Authorization: Bearer {admin_token}`
- **Body**:
```json
{
    "isbn": "9780132350884",
    "quantity": 25,
    "admin_id": 1
}
```

##### 5. Confirm Publisher Order
- **URL**: `PUT /api/admin/publisher-orders/{order_id}/confirm`
- **Headers**: `Authorization: Bearer {admin_token}`
- **Note**: This will automatically add the ordered quantity to the book's stock

### 👤 Customer Endpoints (Require Customer Token)

#### 1. Get Customer Profile
- **URL**: `GET /api/customer/profile`
- **Headers**: `Authorization: Bearer {customer_token}`

#### 2. Update Customer Profile
- **URL**: `PATCH /api/customer/profile`
- **Headers**: `Authorization: Bearer {customer_token}`
- **Body**:
```json
{
    "first_name": "Alice Updated",
    "email": "alice.updated@example.com"
}
```

## 🎯 Testing Scenarios

### Scenario 1: Basic Book Search Flow
1. Get all books: `GET /api/books`
2. Search by category: `GET /api/books/category/Science`
3. Get specific book: `GET /api/books/9780132350884`
4. Search by author: `GET /api/books/author/Robert%20C.%20Martin`

### Scenario 2: Admin Book Management
1. Login as admin: `POST /api/auth/login` (role: "admin")
2. Add new book: `POST /api/admin/books`
3. Update book details: `PUT /api/admin/books/{isbn}`
4. Update stock: `PATCH /api/admin/books/{isbn}/stock`

### Scenario 3: Automatic Ordering Trigger
1. Login as admin
2. Find a book with stock > threshold (e.g., "A Game of Thrones" has 8 stock, threshold 15)
3. Update stock to below threshold: `PATCH /api/admin/books/9780553103540/stock` with `{"quantity": 3}`
4. Check pending orders: `GET /api/admin/publisher-orders/pending` - should show new auto-created order
5. Confirm the order: `PUT /api/admin/publisher-orders/{order_id}/confirm`
6. Check book stock again - should be increased by ordered quantity

### Scenario 4: Database Trigger Validation
1. Try to set negative stock: `PATCH /api/admin/books/{isbn}/stock` with `{"quantity": -5}`
2. Should receive 400 error with message about negative stock
3. Try to add book with invalid category: Should receive 400 error
4. Try to add book with negative threshold: Should receive 400 error

### Scenario 5: Customer Authentication Flow
1. Login as customer: `POST /api/auth/login` (role: "customer")
2. Get profile: `GET /api/customer/profile`
3. Update profile: `PATCH /api/customer/profile`
4. Try to access admin endpoint without admin token: Should receive 401 error

## 🧪 Postman Testing Guide

### Import Collection
1. Open Postman
2. Click "Import" in the top left
3. Select "File" tab
4. Choose `Bookstore-API-Collection.postman_collection.json`
5. Click "Import"

### Collection Variables
The collection includes these variables:
- `baseUrl`: `http://localhost:3001/api` (API base URL)
- `adminToken`: Stores admin authentication token
- `customerToken`: Stores customer authentication token

### Test Sequence
1. **Start with Authentication**
   - Run "Admin Login" - sets adminToken
   - Run "Customer Login" - sets customerToken

2. **Test Public Endpoints** (No auth required)
   - "Get All Books"
   - "Get Book by ISBN"
   - "Get Books by Category"
   - "Get Books by Author"

3. **Test Admin Endpoints** (Requires adminToken)
   - "Add New Book"
   - "Update Book Stock (Triggers Auto-Order)"
   - "Get Publisher Orders"
   - "Confirm Publisher Order"

4. **Test Database Triggers**
   - "Test Negative Stock Prevention" - Should fail with 400 error
   - "Test Auto-Order Trigger" - Should create automatic order

### Expected Results

#### ✅ Success Scenarios
- All GET requests return 200 with proper data
- Admin login returns 200 with token
- Book creation returns 201 with book details
- Stock updates return 200 with confirmation

#### ❌ Expected Failures
- Wrong credentials return 401
- Invalid categories return 400
- Negative stock returns 400
- Unauthorized access to admin endpoints returns 401

## 🔍 Testing Database Triggers

### Trigger 1: Prevent Negative Stock
```bash
# This should fail with 400 error
curl -X PATCH http://localhost:3001/api/admin/books/9780132350884/stock \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{"quantity": -5}'
```

### Trigger 2: Automatic Ordering
```bash
# This should trigger automatic order creation
curl -X PATCH http://localhost:3001/api/admin/books/9780553103540/stock \
  -H "Authorization: Bearer {admin_token}" \
  -H "Content-Type: application/json" \
  -d '{"quantity": 3}'  # Below threshold of 15

# Check pending orders
curl -X GET http://localhost:3001/api/admin/publisher-orders/pending \
  -H "Authorization: Bearer {admin_token}"
```

### Trigger 3: Stock Update on Order Confirmation
```bash
# Confirm an order
curl -X PUT http://localhost:3001/api/admin/publisher-orders/3/confirm \
  -H "Authorization: Bearer {admin_token}"

# Check updated stock
curl -X GET http://localhost:3001/api/books/9780553103540
```

## 🐛 Common Issues & Solutions

### Issue: "Address already in use"
- **Solution**: Kill existing process: `pkill -f "go run cmd/api/main.go"`

### Issue: Admin login fails
- **Solution**: Ensure JWT_SECRET environment variable is set: `export JWT_SECRET=your-secret`

### Issue: Database connection fails
- **Solution**: Check PostgreSQL is running and environment variables are correct

### Issue: Token expires
- **Solution**: Login again to get new token (tokens expire after 24 hours)

## 📊 Test Data Reference

### Sample Books from Seed Data
- ISBN: `9780132350884` - Clean Code (Science, 50 stock, threshold 10)
- ISBN: `9780553103540` - A Game of Thrones (History, 8 stock, threshold 15)
- ISBN: `9780747532743` - Harry Potter (Art, 100 stock, threshold 20)

### Admin Users
- Username: `admin1`, Password: `admin_pass_1`
- Username: `admin2`, Password: `admin_pass_2`

### Customer Users
- Username: `alice`, Password: `hashed_password_1`
- Username: `bob`, Password: `hashed_password_2`

## ✅ Complete Test Checklist

- [ ] Server starts without errors
- [ ] Customer login works
- [ ] Admin login works
- [ ] All public book search endpoints work
- [ ] Admin can add new book
- [ ] Admin can update existing book
- [ ] Admin can update book stock
- [ ] Negative stock is prevented (400 error)
- [ ] Automatic order is created when stock drops below threshold
- [ ] Admin can view publisher orders
- [ ] Admin can confirm publisher orders
- [ ] Stock is updated when order is confirmed
- [ ] Customer profile endpoints work
- [ ] Unauthorized access is properly blocked
- [ ] Database constraints work (category validation, etc.)

## 🎯 Success Criteria

Your implementation is complete when:
1. ✅ All CRUD operations for books work correctly
2. ✅ Database triggers prevent invalid operations
3. ✅ Automatic ordering triggers work
4. ✅ Authentication and authorization work correctly
5. ✅ All endpoints return appropriate HTTP status codes
6. ✅ Error messages are clear and helpful
7. ✅ Stock management is accurate and consistent

Happy Testing! 🚀