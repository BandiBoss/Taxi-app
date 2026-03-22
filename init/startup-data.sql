/* This is an initial admin user. Admin user password is "0000" */
INSERT INTO users (username, password_hash, role) VALUES ('admin', '$2a$10$N.HONssc0FctGNeWhoXpyOATK2lH1BNDYk4eylczvCSeCMP8x//wG', 'admin');

/* Sample drivers */
INSERT INTO drivers (name, phone, car_model, license_plate) VALUES 
('John Doe', '123-456-7890', 'Toyota Camry', 'ABC-1234'),
('Jane Smith', '987-654-3210', 'Honda Accord', 'XYZ-5678');