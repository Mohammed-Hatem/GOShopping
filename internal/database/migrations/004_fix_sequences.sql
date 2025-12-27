-- Fix sequences after seed data insertion
-- This ensures sequences are synchronized with existing data

-- Fix sales_order sequence
SELECT setval('sales_order_order_id_seq', COALESCE((SELECT MAX(order_id) FROM sales_order), 1), true);

-- Fix publisher_order sequence
SELECT setval('publisher_order_rep_order_id_seq', COALESCE((SELECT MAX(rep_order_id) FROM publisher_order), 1), true);

-- Fix shopping_cart sequence
SELECT setval('shopping_cart_cart_id_seq', COALESCE((SELECT MAX(cart_id) FROM shopping_cart), 1), true);

-- Fix administrator sequence
SELECT setval('administrator_admin_id_seq', COALESCE((SELECT MAX(admin_id) FROM administrator), 1), true);

-- Fix publisher sequence
SELECT setval('publisher_publisher_id_seq', COALESCE((SELECT MAX(publisher_id) FROM publisher), 1), true);

