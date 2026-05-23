-- Add tsvector column for full-text search
ALTER TABLE heroes ADD COLUMN search_vector tsvector;

-- Function to update search_vector
CREATE OR REPLACE FUNCTION update_heroes_search_vector()
RETURNS TRIGGER AS $$
BEGIN
  NEW.search_vector := 
    setweight(to_tsvector('english', COALESCE(NEW.name, '')), 'A') ||
    setweight(to_tsvector('english', COALESCE(NEW.description, '')), 'B') ||
    setweight(to_tsvector('english', COALESCE(NEW.universe, '')), 'C');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update search_vector on insert/update
CREATE TRIGGER heroes_search_vector_trigger
BEFORE INSERT OR UPDATE ON heroes
FOR EACH ROW
EXECUTE FUNCTION update_heroes_search_vector();

-- Create GIN index for full-text search performance
CREATE INDEX idx_heroes_search_vector ON heroes USING GIN(search_vector);

-- Populate search_vector for existing records
UPDATE heroes SET search_vector = 
  setweight(to_tsvector('english', COALESCE(name, '')), 'A') ||
  setweight(to_tsvector('english', COALESCE(description, '')), 'B') ||
  setweight(to_tsvector('english', COALESCE(universe, '')), 'C');
