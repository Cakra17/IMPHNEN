CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'cancelled');

CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  customer_id INT NOT NULL,
  total_price NUMERIC(18,2) NOT NULL,
  status order_status NOT NULL DEFAULT 'pending',
  order_date TIMESTAMPTZ DEFAULT NOW(),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_orders_user
    FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_customer_orders
    FOREIGN KEY (customer_id)
    REFERENCES customers(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_customer_status ON orders(customer_id, status);
CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status);
