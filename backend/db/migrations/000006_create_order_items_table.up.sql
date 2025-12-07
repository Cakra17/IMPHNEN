CREATE TABLE IF NOT EXISTS order_items (
  id UUID PRIMARY KEY,
  order_id UUID NOT NULL,
  product_id UUID NOT NULL,
  quantity INT NOT NULL,
  total_price NUMERIC(18,2) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_order_item_product
    FOREIGN KEY (product_id)
    REFERENCES products(id) ON DELETE CASCADE
);
