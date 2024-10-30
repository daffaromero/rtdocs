CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) CHECK (role IN ('authenticated', 'guest')) NOT NULL DEFAULT 'authenticated',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE docs
  ADD COLUMN "owner_id" uuid REFERENCES users(id) ON DELETE CASCADE;
  ADD COLUMN "is_public" BOOLEAN DEFAULT FALSE;
  ADD COLUMN "can_edit" BOOLEAN DEFAULT TRUE;
  ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
  ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

CREATE INDEX idx_docs_owner_id ON docs(owner_id);
CREATE INDEX idx_docs_is_public ON docs(is_public);