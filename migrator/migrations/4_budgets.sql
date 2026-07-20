CREATE TABLE budgets(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    category_id INT REFERENCES categories(id),
    limit_value NUMERIC(11,2) NOT NULL
);

CREATE UNIQUE INDEX idx_unique_budgets_user ON budgets(user_id, category_id);