DROP INDEX IF EXISTS idx_orders_user_status;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_product_id;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP TABLE IF EXISTS orders CASCADE;
DROP TYPE IF EXISTS order_status;
