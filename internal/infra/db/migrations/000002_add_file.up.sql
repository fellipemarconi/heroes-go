CREATE TYPE file_type AS ENUM (
    'profile',
    'image',
    'video',
    'document'
);

CREATE TABLE files (
   id UUID PRIMARY KEY,
   path TEXT NOT NULL,
   type file_type NOT NULL,
   metadata JSONB NOT NULL DEFAULT '{}',
   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   created_at TIMESTAMP DEFAULT now(),
   updated_at TIMESTAMP DEFAULT now()
);

ALTER TABLE users ADD COLUMN image TEXT;