CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    name VARCHAR(25) NOT NULL
);

CREATE UNIQUE INDEX idx_unique_system_categories ON categories(name) WHERE user_id IS NULL;
CREATE UNIQUE INDEX idx_unique_user_categories ON categories(name, user_id) WHERE user_id IS NOT NULL;

TRUNCATE transactions;

ALTER TABLE transactions ADD category_id INT REFERENCES categories(id) NOT NULL;