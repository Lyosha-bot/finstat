INSERT INTO categories(id, name) VALUES (0, 'Без категории');

CREATE OR REPLACE FUNCTION reset_category_to_zero()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE transactions 
    SET category_id = 0 
    WHERE category_id = OLD.id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_delete_category
BEFORE DELETE ON categories
FOR EACH ROW
EXECUTE FUNCTION reset_category_to_zero();

ALTER TABLE budgets DROP CONSTRAINT budgets_category_id_fkey;

ALTER TABLE budgets 
ADD CONSTRAINT budgets_category_id_fkey 
FOREIGN KEY (category_id)
REFERENCES categories(id)
ON DELETE CASCADE;