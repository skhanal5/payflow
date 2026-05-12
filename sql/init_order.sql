CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  order_id VARCHAR NOT NULL UNIQUE,
  status VARCHAR NOT NULL,
  created_at TIMESTAMP  DEFAULT now(),
  updated_at TIMESTAMP  DEFAULT now(),
  deleted_at TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS order_items (
  id SERIAL PRIMARY KEY,
  order_id VARCHAR NOT NULL,
  product_id VARCHAR NOT NULL,
  quantity INTEGER NOT NULL,
  created_at TIMESTAMP  DEFAULT now(),
  updated_at TIMESTAMP  DEFAULT now(),
  deleted_at TIMESTAMP ,
  FOREIGN KEY (order_id) REFERENCES orders(order_id)
);


CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER trigger_set_updated_at_orders
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


CREATE TRIGGER trigger_set_updated_at_order_items
BEFORE UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

INSERT INTO orders (order_id, status)
VALUES ('ORD123456', 'PROCESSING')
ON CONFLICT (order_id) DO NOTHING;


INSERT INTO order_items (order_id, product_id, quantity)
VALUES 
  ('ORD123456', 'PROD001', 2),
  ('ORD123456', 'PROD002', 1)
ON CONFLICT DO NOTHING;