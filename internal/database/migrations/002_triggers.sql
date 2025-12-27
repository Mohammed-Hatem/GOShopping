
CREATE OR REPLACE FUNCTION prevent_negative_stock()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.stock_quantity < 0 THEN
        RAISE EXCEPTION 'Stock cannot be negative';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_negative_stock_trigger
BEFORE UPDATE OF stock_quantity ON book
FOR EACH ROW
EXECUTE FUNCTION prevent_negative_stock();


CREATE OR REPLACE FUNCTION auto_order_low_stock()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.stock_quantity <= NEW.threshold THEN
        INSERT INTO publisher_order (order_date, status, quantity, isbn, admin_id)
        VALUES (CURRENT_DATE, 'Pending', 20, NEW.isbn, 1);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER auto_order_low_stock_trigger
AFTER UPDATE OF stock_quantity ON book
FOR EACH ROW
EXECUTE FUNCTION auto_order_low_stock();


CREATE OR REPLACE FUNCTION add_stock_on_confirm()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'Confirmed' THEN
        UPDATE book
        SET stock_quantity = stock_quantity + NEW.quantity
        WHERE isbn = NEW.isbn;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER add_stock_on_confirm_trigger
AFTER UPDATE OF status ON publisher_order
FOR EACH ROW
EXECUTE FUNCTION add_stock_on_confirm();


ALTER TABLE book
ADD CONSTRAINT stock_not_negative CHECK (stock_quantity >= 0);

ALTER TABLE book
ADD CONSTRAINT threshold_positive CHECK (threshold > 0);

ALTER TABLE publisher_order
ADD CONSTRAINT valid_order_status
CHECK (status IN ('Pending', 'Confirmed', 'Cancelled'));
