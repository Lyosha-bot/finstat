CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    amount NUMERIC(11, 2),
    description VARCHAR(100),
    date DATE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_date ON transactions (user_id, date DESC);