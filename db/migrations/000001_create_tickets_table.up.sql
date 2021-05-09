CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    price DECIMAL NOT NULL CHECK (price > 0),
    user_id INTEGER NOT NULL,
    order_id INTEGER DEFAULT 0,
    version INTEGER DEFAULT 1
);