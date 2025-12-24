-- Sample seed data for bookstore schema

-- 1. Publishers
INSERT INTO PUBLISHER (Publisher_ID, Name, Address, Phone) VALUES
    (1, 'O''Reilly Media', '1005 Gravenstein Highway North, Sebastopol, CA', '+1-800-998-9938'),
    (2, 'Penguin Random House', '1745 Broadway, New York, NY', '+1-212-782-9000'),
    (3, 'HarperCollins', '195 Broadway, New York, NY', '+1-212-207-7000');

-- 2. Authors
INSERT INTO AUTHOR (Name) VALUES
    ('Robert C. Martin'),
    ('Martin Fowler'),
    ('J.K. Rowling'),
    ('George R.R. Martin');

-- 3. Books
INSERT INTO BOOK (ISBN, Title, Publication_Year, Selling_Price, Category, Author_Name, Stock_Quantity, Threshold, Publisher_ID) VALUES
    ('9780132350884', 'Clean Code', 2008, 39.99, 'Programming', 'Robert C. Martin', 50, 10, 1),
    ('9780201485677', 'Refactoring', 1999, 44.99, 'Programming', 'Martin Fowler', 40, 10, 2),
    ('9780747532743', 'Harry Potter and the Philosopher''s Stone', 1997, 19.99, 'Fantasy', 'J.K. Rowling', 100, 20, 2),
    ('9780553103540', 'A Game of Thrones', 1996, 24.99, 'Fantasy', 'George R.R. Martin', 80, 15, 3);

-- 4. Customers
INSERT INTO CUSTOMER (Username, Password, First_Name, Last_Name, Email, Phone_Number, Shipping_Address) VALUES
    ('alice', 'hashed_password_1', 'Alice', 'Johnson', 'alice@example.com', '+1-555-0001', '123 Main St, Springfield'),
    ('bob',   'hashed_password_2', 'Bob',   'Smith',   'bob@example.com',   '+1-555-0002', '456 Oak Ave, Springfield');

-- 5. Administrators
INSERT INTO ADMINISTRATOR (Admin_ID, Username, Password) VALUES
    (1, 'admin1', 'admin_pass_1'),
    (2, 'admin2', 'admin_pass_2');

-- 6. Shopping carts
INSERT INTO SHOPPING_CART (Cart_ID, Username) VALUES
    (1, 'alice'),
    (2, 'bob');

-- 7. Cart items
INSERT INTO CART_ITEM (Cart_ID, ISBN, Quantity) VALUES
    (1, '9780132350884', 1),
    (1, '9780747532743', 2),
    (2, '9780553103540', 1);

-- 8. Sales orders
INSERT INTO SALES_ORDER (Order_ID, Order_Date, Total_Amount, Credit_Card_No, Expiry_Date, Username) VALUES
    (1, CURRENT_DATE - INTERVAL '2 days', 79.98, '4111111111111111', CURRENT_DATE + INTERVAL '1 year', 'alice'),
    (2, CURRENT_DATE - INTERVAL '1 day', 24.99, '4000123412341234', CURRENT_DATE + INTERVAL '2 years', 'bob');

-- 9. Order items
INSERT INTO ORDER_ITEM (Order_ID, ISBN, Quantity, Unit_Price) VALUES
    (1, '9780132350884', 1, 39.99),
    (1, '9780747532743', 2, 19.99),
    (2, '9780553103540', 1, 24.99);

-- 10. Publisher replenishment orders
INSERT INTO PUBLISHER_ORDER (Rep_Order_ID, Order_Date, Status, Quantity, ISBN, Admin_ID) VALUES
    (1, CURRENT_DATE - INTERVAL '5 days', 'Pending', 100, '9780747532743', 1),
    (2, CURRENT_DATE - INTERVAL '3 days', 'Completed', 50,  '9780132350884', 2);
