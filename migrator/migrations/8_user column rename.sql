ALTER TABLE users RENAME COLUMN applied_at TO created_at;
ALTER TABLE users ALTER COLUMN created_at SET DEFAULT NOW();