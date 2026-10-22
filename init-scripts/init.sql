-- 1. Create the table if it doesn't exist
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    worker_id INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Insert mock data for pool testing
INSERT INTO tasks (title, status) VALUES
('Process user profile uploads', 'pending'),
('Generate monthly financial invoice', 'pending'),
('Sync inventory with supplier API', 'pending'),
('Send welcome email sequences', 'pending'),
('Optimize database indexes', 'pending');
