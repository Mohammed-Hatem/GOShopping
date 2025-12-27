-- Sample seed data for bookstore schema

-- 1. Publishers
INSERT INTO PUBLISHER (Publisher_ID, Name, Address, Phone) VALUES
    (1, 'O''Reilly Media', '1005 Gravenstein Highway North, Sebastopol, CA', '+1-800-998-9938'),
    (2, 'Penguin Random House', '1745 Broadway, New York, NY', '+1-212-782-9000'),
    (3, 'HarperCollins', '195 Broadway, New York, NY', '+1-212-207-7000')
ON CONFLICT (Publisher_ID) DO NOTHING;

-- 2. Authors
INSERT INTO AUTHOR (Name, ISBN) VALUES
    ('Robert C. Martin', '9780132350884'),
    ('Martin Fowler', '9780201485677'),
    ('J.K. Rowling', '9780747532743'),
    ('George R.R. Martin', '9780553103540')
ON CONFLICT (Name, ISBN) DO NOTHING;

-- 3. Books
INSERT INTO BOOK (ISBN, Title, Publication_Year, Selling_Price, Category, Author_Name, Stock_Quantity, Threshold, Publisher_ID) VALUES
    ('9780132350884', 'Clean Code', 2008, 39.99, 'Science', 'Robert C. Martin', 50, 10, 1),
    ('9780201485677', 'Refactoring', 1999, 44.99, 'Science', 'Martin Fowler', 40, 10, 2),
    ('9780747532743', 'Harry Potter and the Philosopher''s Stone', 1997, 19.99, 'Art', 'J.K. Rowling', 100, 20, 2),
    ('9780553103540', 'A Game of Thrones', 1996, 24.99, 'History', 'George R.R. Martin', 8, 15, 3),
    ('9780321765723', 'The Art of Computer Programming', 1968, 79.99, 'Science', 'Donald Knuth', 25, 5, 1),
    ('9780134685991', 'Effective Java', 2018, 49.99, 'Science', 'Joshua Bloch', 30, 8, 1),
    ('9780262033848', 'Introduction to Algorithms', 2009, 84.99, 'Science', 'Thomas Cormen', 15, 10, 2),
    ('9780062315007', 'Sapiens: A Brief History', 2015, 18.99, 'History', 'Yuval Noah Harari', 60, 12, 3),
    ('9780143126569', 'The Geography of Thought', 2003, 16.99, 'Geography', 'Richard Nisbett', 35, 7, 2),
    ('9780061234006', 'The Power of Myth', 1988, 14.99, 'Religion', 'Joseph Campbell', 45, 9, 3)
ON CONFLICT (ISBN) DO NOTHING;

-- 4. Customers
INSERT INTO CUSTOMER (Username, Password, First_Name, Last_Name, Email, Phone_Number, Shipping_Address) VALUES
    ('alice', 'hashed_password_1', 'Alice', 'Johnson', 'alice@example.com', '+1-555-0001', '123 Main St, Springfield'),
    ('bob',   'hashed_password_2', 'Bob',   'Smith',   'bob@example.com',   '+1-555-0002', '456 Oak Ave, Springfield')
ON CONFLICT (Username) DO NOTHING;

-- 5. Administrators
INSERT INTO ADMINISTRATOR (Admin_ID, Username, Password) VALUES
    (1, 'admin1', 'admin_pass_1'),
    (2, 'admin2', 'admin_pass_2')
ON CONFLICT (Admin_ID) DO NOTHING;

-- 6. Shopping carts
INSERT INTO SHOPPING_CART (Cart_ID, Username) VALUES
    (1, 'alice'),
    (2, 'bob')
ON CONFLICT (Username) DO NOTHING;

-- 7. Cart items
INSERT INTO CART_ITEM (Cart_ID, ISBN, Quantity) VALUES
    (1, '9780132350884', 1),
    (1, '9780747532743', 2),
    (2, '9780553103540', 1)
ON CONFLICT (Cart_ID, ISBN) DO NOTHING;

-- 8. Sales orders
INSERT INTO SALES_ORDER (Order_ID, Order_Date, Total_Amount, Credit_Card_No, Expiry_Date, Username) VALUES
    (1, CURRENT_DATE - INTERVAL '2 days', 79.98, '4111111111111111', CURRENT_DATE + INTERVAL '1 year', 'alice'),
    (2, CURRENT_DATE - INTERVAL '1 day', 24.99, '4000123412341234', CURRENT_DATE + INTERVAL '2 years', 'bob')
ON CONFLICT (Order_ID) DO NOTHING;

-- 9. Order items
INSERT INTO ORDER_ITEM (Order_ID, ISBN, Quantity, Unit_Price) VALUES
    (1, '9780132350884', 1, 39.99),
    (1, '9780747532743', 2, 19.99),
    (2, '9780553103540', 1, 24.99)
ON CONFLICT (Order_ID, ISBN) DO NOTHING;

-- 10. Publisher replenishment orders
INSERT INTO PUBLISHER_ORDER (Rep_Order_ID, Order_Date, Status, Quantity, ISBN, Admin_ID) VALUES
    (1, CURRENT_DATE - INTERVAL '5 days', 'Pending', 100, '9780747532743', 1),
    (2, CURRENT_DATE - INTERVAL '3 days', 'Confirmed', 50,  '9780132350884', 2),
    (3, CURRENT_DATE - INTERVAL '1 day', 'Pending', 20, '9780553103540', 1)
ON CONFLICT (Rep_Order_ID) DO NOTHING;

-- Ensure SERIAL sequences are ahead of seeded IDs
SELECT setval(pg_get_serial_sequence('publisher', 'publisher_id'), COALESCE((SELECT MAX(publisher_id) FROM publisher), 1));
SELECT setval(pg_get_serial_sequence('administrator', 'admin_id'), COALESCE((SELECT MAX(admin_id) FROM administrator), 1));
SELECT setval(pg_get_serial_sequence('shopping_cart', 'cart_id'), COALESCE((SELECT MAX(cart_id) FROM shopping_cart), 1));
SELECT setval(pg_get_serial_sequence('sales_order', 'order_id'), COALESCE((SELECT MAX(order_id) FROM sales_order), 1));
SELECT setval(pg_get_serial_sequence('publisher_order', 'rep_order_id'), COALESCE((SELECT MAX(rep_order_id) FROM publisher_order), 1));
