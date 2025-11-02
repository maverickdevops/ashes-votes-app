CREATE TABLE IF NOT EXISTS votes (
  id SERIAL PRIMARY KEY,
  team VARCHAR(32) NOT NULL CHECK (team IN ('australia','england')),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_votes_team ON votes(team);
