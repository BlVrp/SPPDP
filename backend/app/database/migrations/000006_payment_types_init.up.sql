INSERT INTO payment_types(type) VALUES
('STRIPE')
ON CONFLICT DO NOTHING;
