-- Sample data for sale-service database
-- This file contains test data for sale_recipe and sale_unit tables

-- Insert sample sale units (products)
INSERT INTO sale (name, pos_id, price, recipe_id, created_at) VALUES
-- Beverages
('Fresh Milk Latte', 10001, 45000, 1, '2025-01-15 08:00:00+07'),
('Iced Americano', 10002, 35000, 2, '2025-01-15 08:15:00+07'),
('Cappuccino', 10003, 42000, 3, '2025-01-15 08:30:00+07'),
('Espresso', 10004, 30000, 4, '2025-01-15 08:45:00+07'),
('Mocha Frappe', 10005, 55000, 5, '2025-01-15 09:00:00+07'),

-- Food Items
('Croissant', 20001, 25000, 6, '2025-01-15 09:15:00+07'),
('Chocolate Cake', 20002, 48000, 7, '2025-01-15 09:30:00+07'),
('Blueberry Muffin', 20003, 32000, 8, '2025-01-15 09:45:00+07'),
('Sandwich', 20004, 38000, 9, '2025-01-15 10:00:00+07'),
('Bagel with Cream Cheese', 20005, 35000, 10, '2025-01-15 10:15:00+07'),

-- Specialty Drinks
('Matcha Latte', 30001, 50000, 1, '2025-01-16 08:00:00+07'),
('Caramel Macchiato', 30002, 52000, 2, '2025-01-16 08:30:00+07'),
('Vietnamese Coffee', 30003, 40000, 3, '2025-01-16 09:00:00+07'),
('Thai Tea', 30004, 38000, 4, '2025-01-16 09:30:00+07'),
('Fruit Smoothie', 30005, 45000, 5, '2025-01-16 10:00:00+07'),

-- Desserts
('Tiramisu', 40001, 55000, 6, '2025-01-17 08:00:00+07'),
('Cheesecake', 40002, 52000, 7, '2025-01-17 08:30:00+07'),
('Brownie', 40003, 35000, 8, '2025-01-17 09:00:00+07'),
('Ice Cream Sundae', 40004, 42000, 9, '2025-01-17 09:30:00+07'),
('Apple Pie', 40005, 38000, 10, '2025-01-17 10:00:00+07'),

-- Breakfast Items
('Pancakes', 50001, 45000, 1, '2025-01-18 07:00:00+07'),
('French Toast', 50002, 42000, 2, '2025-01-18 07:30:00+07'),
('Breakfast Burrito', 50003, 55000, 3, '2025-01-18 08:00:00+07'),
('Omelette', 50004, 48000, 4, '2025-01-18 08:30:00+07'),
('Granola Bowl', 50005, 40000, 5, '2025-01-18 09:00:00+07');

-- Reset sequence to continue from the last inserted ID
SELECT setval('sale_id_seq', (SELECT MAX(id) FROM sale));

