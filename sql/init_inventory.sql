CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  product_id VARCHAR NOT NULL UNIQUE,
  name VARCHAR NOT NULL,
  quantity INTEGER NOT NULL,
  created_at TIMESTAMP  DEFAULT now(),
  updated_at TIMESTAMP  DEFAULT now(),
  deleted_at TIMESTAMP 
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_updated_at_products
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


INSERT INTO products (product_id, name, quantity)
VALUES 
  ('PROD001', 'USB Keyboard', 50),
  ('PROD002', 'Wireless Mouse', 80)
ON CONFLICT DO NOTHING;
