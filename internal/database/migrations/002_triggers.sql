-- Database triggers for bookstore order processing system

-- Trigger 1: Prevent negative stock quantities (before update)
CREATE OR REPLACE FUNCTION prevent_negative_stock()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.stock_quantity < 0 THEN
        RAISE EXCEPTION 'Stock quantity cannot be negative. Current stock: %, Attempted update: %', 
                        OLD.stock_quantity, NEW.stock_quantity;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_prevent_negative_stock ON book;
CREATE TRIGGER trigger_prevent_negative_stock
    BEFORE UPDATE ON book
    FOR EACH ROW
    WHEN (OLD.stock_quantity IS DISTINCT FROM NEW.stock_quantity)
    EXECUTE FUNCTION prevent_negative_stock();

-- Trigger 2: Automatic publisher order when stock drops below threshold (after update)
CREATE OR REPLACE FUNCTION auto_order_when_below_threshold()
RETURNS TRIGGER AS $$
DECLARE
    order_quantity INT := 20; -- Constant order quantity as specified
    admin_id INT := 1; -- Default admin ID (should be configurable)
BEGIN
    -- Check if stock dropped from above threshold to below or equal threshold
    IF OLD.stock_quantity > OLD.threshold AND NEW.stock_quantity <= NEW.threshold THEN
        -- Insert automatic publisher order
        INSERT INTO publisher_order (order_date, status, quantity, isbn, admin_id)
        VALUES (CURRENT_DATE, 'Pending', order_quantity, NEW.isbn, admin_id);
        
        RAISE NOTICE 'Automatic order placed for ISBN %: % units', NEW.isbn, order_quantity;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_auto_order ON book;
CREATE TRIGGER trigger_auto_order
    AFTER UPDATE ON book
    FOR EACH ROW
    WHEN (OLD.stock_quantity IS DISTINCT FROM NEW.stock_quantity)
    EXECUTE FUNCTION auto_order_when_below_threshold();

-- Trigger 3: Update stock when publisher order is confirmed
CREATE OR REPLACE FUNCTION update_stock_on_order_confirmation()
RETURNS TRIGGER AS $$
BEGIN
    -- Only update stock if status is being changed to 'Confirmed'
    IF NEW.status = 'Confirmed' AND OLD.status != 'Confirmed' THEN
        UPDATE book 
        SET stock_quantity = stock_quantity + NEW.quantity 
        WHERE isbn = NEW.isbn;
        
        RAISE NOTICE 'Stock updated for ISBN %: added % units', NEW.isbn, NEW.quantity;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_stock_on_confirmation ON publisher_order;
CREATE TRIGGER trigger_update_stock_on_confirmation
    AFTER UPDATE ON publisher_order
    FOR EACH ROW
    WHEN (OLD.status IS DISTINCT FROM NEW.status)
    EXECUTE FUNCTION update_stock_on_order_confirmation();

-- Add check constraint for book categories
ALTER TABLE book 
ADD CONSTRAINT chk_book_category 
CHECK (category IN ('Science', 'Art', 'Religion', 'History', 'Geography'));

-- Add check constraint for publisher order status
ALTER TABLE publisher_order 
ADD CONSTRAINT chk_order_status 
CHECK (status IN ('Pending', 'Confirmed', 'Cancelled'));

-- Add check constraint for positive stock quantity
ALTER TABLE book 
ADD CONSTRAINT chk_positive_stock 
CHECK (stock_quantity >= 0);

-- Add check constraint for positive threshold
ALTER TABLE book 
ADD CONSTRAINT chk_positive_threshold 
CHECK (threshold > 0);