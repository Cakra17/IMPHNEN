CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  price NUMERIC(18,2) NOT NULL,
  stock INT NOT NULL,
  image_url TEXT NOT NULL,
  public_id TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_user_product
    FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_products_user ON products(user_id)
