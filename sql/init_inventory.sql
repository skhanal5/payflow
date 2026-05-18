CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  product_id VARCHAR NOT NULL UNIQUE,
  category VARCHAR NOT NULL DEFAULT '',
  name VARCHAR NOT NULL,
  description VARCHAR NOT NULL DEFAULT '',
  price DOUBLE PRECISION NOT NULL DEFAULT 0,
  available_stock INTEGER NOT NULL DEFAULT 0,
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


INSERT INTO products (product_id, category, name, description, price, available_stock)
VALUES 
  ('PROD001', 'Electronics', 'Wireless Mechanical Keyboard', 'Full-size mechanical keyboard with hot-swappable switches and RGB backlighting.', 79.99, 50),
  ('PROD002', 'Electronics', 'USB-C Fast Charging Cable', '6ft braided USB-C to USB-C cable with 100W Power Delivery support.', 12.99, 200),
  ('PROD003', 'Electronics', 'Bluetooth Noise-Cancelling Headphones', 'Over-ear headphones with active noise cancellation and 30-hour battery life.', 249.99, 30),
  ('PROD004', 'Furniture', 'Ergonomic Office Chair', 'Adjustable lumbar support, mesh back, and breathable cushion for all-day comfort.', 399.99, 15),
  ('PROD005', 'Furniture', 'Standing Desk Converter', 'Height-adjustable desktop riser with dual monitor support and keyboard tray.', 179.99, 25),
  ('PROD006', 'Electronics', '27" 4K UHD Monitor', 'IPS panel with HDR400, USB-C connectivity, and built-in speakers.', 449.99, 20),
  ('PROD007', 'Electronics', 'Qi Wireless Charging Pad', 'Fast wireless charger compatible with all Qi-enabled devices. 15W max output.', 24.99, 100),
  ('PROD008', 'Accessories', 'Adjustable Laptop Stand', 'Aluminum laptop stand with ergonomic elevation and ventilated design.', 34.99, 60)
ON CONFLICT DO NOTHING;
