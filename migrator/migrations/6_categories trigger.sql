DELETE FROM categories
WHERE user_id IS NOT NULL;


CREATE OR REPLACE FUNCTION check_system_category_collision()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.user_id IS NOT NULL AND EXISTS (
        SELECT 1
        FROM categories
        WHERE user_id IS NULL AND NEW.name = name
    ) THEN
        RAISE EXCEPTION 'Category with name %s already exists', NEW.name
        USING ERRCODE = 'unique_violation';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_system_category_collision
BEFORE INSERT OR UPDATE ON categories
FOR EACH ROW EXECUTE FUNCTION check_system_category_collision();