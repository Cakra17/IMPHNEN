CREATE TYPE type AS ENUM ('expense', 'income');
CREATE TYPE source AS ENUM ('receipt', 'bot', 'manual');

CREATE TABLE transactions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  type type NOT NULL,
  source source NOT NULL,
  amount DECIMAL(18,2) NOT NULL,
  transaction_date TIMESTAMP NOT NULL,
  receipt_id UUID DEFAULT NULL,
  order_id UUID DEFAULT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CONSTRAINT fk_user_transactions
    FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_transactions
    FOREIGN KEY (receipt_id)
    REFERENCES receipts(id) ON DELETE CASCADE,
  CONSTRAINT fk_order_transactions
    FOREIGN KEY (order_id)
    REFERENCES orders(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_date ON transactions(user_id, transaction_date);
CREATE INDEX IF NOT EXISTS idx_user_type ON transactions(user_id, type);
CREATE INDEX IF NOT EXISTS idx_user_type_date ON transactions(user_id, type, transaction_date);
