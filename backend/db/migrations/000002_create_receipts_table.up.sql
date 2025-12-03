CREATE TABLE IF NOT EXISTS receipts (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  total_items INT NOT NULL,
  total_price NUMERIC(18,0) NOT NULL,
  store_name VARCHAR(255) NOT NULL,
  image_url TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_receipts_user
    FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE
);
