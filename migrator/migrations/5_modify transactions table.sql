ALTER TABLE transactions RENAME COLUMN amount TO value;
ALTER TABLE transactions ALTER COLUMN value SET NOT NULL;

ALTER TABLE transactions ALTER COLUMN user_id SET NOT NULL;