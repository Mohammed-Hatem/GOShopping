-- 1. Publisher Table
CREATE TABLE PUBLISHER (
    Publisher_ID INT PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Address TEXT,
    Phone VARCHAR(20)
);

-- 2. Author Table
CREATE TABLE AUTHOR (
    Author_ID INT PRIMARY KEY,
    Name VARCHAR(255) NOT NULL
);

-- 3. Book Table
CREATE TABLE BOOK (
    ISBN VARCHAR(20) PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    Publication_Year INT,
    Selling_Price DECIMAL(10, 2),
    Category VARCHAR(100),
    Stock_Quantity INT DEFAULT 0,
    Threshold INT DEFAULT 10,
    Publisher_ID INT,
    FOREIGN KEY (Publisher_ID) REFERENCES PUBLISHER(Publisher_ID)
);

-- 4. Book_Author Junction Table (Many-to-Many relationship)
CREATE TABLE BOOK_AUTHOR (
    ISBN VARCHAR(20),
    Author_ID INT,
    PRIMARY KEY (ISBN, Author_ID),
    FOREIGN KEY (ISBN) REFERENCES BOOK(ISBN),
    FOREIGN KEY (Author_ID) REFERENCES AUTHOR(Author_ID)
);

-- 5. Customer Table
CREATE TABLE CUSTOMER (
    Username VARCHAR(50) PRIMARY KEY,
    Password VARCHAR(255) NOT NULL,
    First_Name VARCHAR(100),
    Last_Name VARCHAR(100),
    Email VARCHAR(255) UNIQUE,
    Phone_Number VARCHAR(20),
    Shipping_Address TEXT
);

-- 6. Administrator Table
CREATE TABLE ADMINISTRATOR (
    Admin_ID INT PRIMARY KEY,
    Username VARCHAR(50) NOT NULL,
    Password VARCHAR(255) NOT NULL
);

-- 7. Shopping Cart Table
CREATE TABLE SHOPPING_CART (
    Cart_ID INT PRIMARY KEY,
    Username VARCHAR(50) UNIQUE,
    FOREIGN KEY (Username) REFERENCES CUSTOMER(Username)
);

-- 8. Cart Item Table
CREATE TABLE CART_ITEM (
    Cart_ID INT,
    ISBN VARCHAR(20),
    Quantity INT NOT NULL DEFAULT 1,
    PRIMARY KEY (Cart_ID, ISBN),
    FOREIGN KEY (Cart_ID) REFERENCES SHOPPING_CART(Cart_ID),
    FOREIGN KEY (ISBN) REFERENCES BOOK(ISBN)
);

-- 9. Sales Order Table
CREATE TABLE SALES_ORDER (
    Order_ID INT PRIMARY KEY,
    Order_Date DATE NOT NULL,
    Total_Amount DECIMAL(10, 2),
    Credit_Card_No VARCHAR(20),
    Expiry_Date DATE,
    Username VARCHAR(50),
    FOREIGN KEY (Username) REFERENCES CUSTOMER(Username)
);

-- 10. Order Item Table
CREATE TABLE ORDER_ITEM (
    Order_ID INT,
    ISBN VARCHAR(20),
    Quantity INT NOT NULL,
    Unit_Price DECIMAL(10, 2),
    PRIMARY KEY (Order_ID, ISBN),
    FOREIGN KEY (Order_ID) REFERENCES SALES_ORDER(Order_ID),
    FOREIGN KEY (ISBN) REFERENCES BOOK(ISBN)
);

-- 11. Publisher Order (Replenishment)
CREATE TABLE PUBLISHER_ORDER (
    Rep_Order_ID INT PRIMARY KEY,
    Order_Date DATE NOT NULL,
    Status VARCHAR(50),
    Quantity INT,
    ISBN VARCHAR(20),
    Admin_ID INT,
    FOREIGN KEY (ISBN) REFERENCES BOOK(ISBN),
    FOREIGN KEY (Admin_ID) REFERENCES ADMINISTRATOR(Admin_ID)
);
