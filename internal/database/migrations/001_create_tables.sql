-- Initial schema for bookstore
-- Note: DROP statements removed to prevent data loss on server restart
-- If you need to reset the database, do it manually via psql

CREATE TABLE IF NOT EXISTS publisher (
    publisher_id SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    address      TEXT,
    phone        VARCHAR(20)
);
-- 2. author Table
CREATE TABLE IF NOT EXISTS author (
    name VARCHAR(255),
    isbn VARCHAR(20),
    PRIMARY KEY (name, isbn)
    CONSTRAINT fk_author_book
        FOREIGN KEY (isbn) REFERENCES book(isbn)
);

-- 3. Book Table
CREATE TABLE IF NOT EXISTS book (
    isbn             VARCHAR(20) PRIMARY KEY,
    title            VARCHAR(255) NOT NULL,
    publication_year INT,
    selling_price    DECIMAL(10, 2),
    category         VARCHAR(100),
    author_name       VARCHAR(255) NOT NULL,
    stock_quantity   INT DEFAULT 0,
    threshold        INT DEFAULT 10,
    publisher_id     INT,
    CONSTRAINT fk_book_publisher
        FOREIGN KEY (publisher_id) REFERENCES publisher(publisher_id)
);

-- 4. Customer Table
CREATE TABLE IF NOT EXISTS customer (
    username         VARCHAR(50) PRIMARY KEY,
    password         VARCHAR(255) NOT NULL,
    first_name       VARCHAR(100),
    last_name        VARCHAR(100),
    email            VARCHAR(255) UNIQUE,
    phone_number     VARCHAR(20),
    shipping_address TEXT
);

-- 5. Administrator Table
CREATE TABLE IF NOT EXISTS administrator (
    admin_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- 6. Shopping Cart Table
CREATE TABLE IF NOT EXISTS shopping_cart (
    cart_id  SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    CONSTRAINT fk_cart_customer
        FOREIGN KEY (username) REFERENCES customer(username)
);

-- 7. Cart Item Table
CREATE TABLE IF NOT EXISTS cart_item (
    cart_id INT,
    isbn    VARCHAR(20),
    quantity INT NOT NULL DEFAULT 1,
    PRIMARY KEY (cart_id, isbn),
    CONSTRAINT fk_cart_item_cart
        FOREIGN KEY (cart_id) REFERENCES shopping_cart(cart_id),
    CONSTRAINT fk_cart_item_book
        FOREIGN KEY (isbn) REFERENCES book(isbn)
);

-- 8. Sales Order Table
CREATE TABLE IF NOT EXISTS sales_order (
    order_id      SERIAL PRIMARY KEY,
    order_date    DATE NOT NULL,
    total_amount  DECIMAL(10, 2),
    credit_card_no VARCHAR(20),
    expiry_date   DATE,
    username      VARCHAR(50),
    CONSTRAINT fk_sales_order_customer
        FOREIGN KEY (username) REFERENCES customer(username)
);

-- 9. Order Item Table
CREATE TABLE IF NOT EXISTS order_item (
    order_id  INT,
    isbn      VARCHAR(20),
    quantity  INT NOT NULL,
    unit_price DECIMAL(10, 2),
    PRIMARY KEY (order_id, isbn),
    CONSTRAINT fk_order_item_order
        FOREIGN KEY (order_id) REFERENCES sales_order(order_id),
    CONSTRAINT fk_order_item_book
        FOREIGN KEY (isbn) REFERENCES book(isbn)
);

-- 10. Publisher Order Table
-- This table tracks orders made by the publisher for books
CREATE TABLE IF NOT EXISTS publisher_order (
    rep_order_id SERIAL PRIMARY KEY,
    order_date   DATE NOT NULL,
    status       VARCHAR(50),
    quantity     INT,
    isbn         VARCHAR(20),
    admin_id     INT,
    CONSTRAINT fk_publisher_order_book
        FOREIGN KEY (isbn) REFERENCES book(isbn),
    CONSTRAINT fk_publisher_order_admin
        FOREIGN KEY (admin_id) REFERENCES administrator(admin_id)
);
