CREATE TABLE IF NOT EXISTS receipts (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  total_items INT NOT NULL,
  total_price NUMERIC(18,2) NOT NULL,
  store_name VARCHAR(255) NOT NULL,
  image_url TEXT NOT NULL,
  public_id TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_receipts_user
    FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_receipts_user_id ON receipts(user_id);
CREATE INDEX IF NOT EXISTS idx_receipts_id ON receipts(id);
