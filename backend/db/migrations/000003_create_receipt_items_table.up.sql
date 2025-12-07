CREATE TABLE IF NOT EXISTS receipt_items (
  id UUID PRIMARY KEY,
  receipt_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  price NUMERIC(18,2) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_receipt_items_receipts
    FOREIGN KEY (receipt_id)
    REFERENCES receipts(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_receipt_id ON receipt_items(receipt_id);
