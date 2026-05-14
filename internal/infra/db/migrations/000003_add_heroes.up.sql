CREATE TABLE heroes (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    universe TEXT NOT NULL,
    powers TEXT[] NOT NULL DEFAULT '{}',
    alignment TEXT NOT NULL,
    description TEXT,
    image TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_heroes_user_id ON heroes(user_id);
CREATE INDEX idx_heroes_alignment ON heroes(alignment);
CREATE INDEX idx_heroes_universe ON heroes(universe);
CREATE INDEX idx_heroes_created_at ON heroes(created_at DESC);