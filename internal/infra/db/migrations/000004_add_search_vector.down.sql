-- Rollback: Remove search vector functionality
DROP INDEX IF EXISTS idx_heroes_search_vector;
DROP TRIGGER IF EXISTS heroes_search_vector_trigger ON heroes;
DROP FUNCTION IF EXISTS update_heroes_search_vector();
ALTER TABLE heroes DROP COLUMN IF EXISTS search_vector;
