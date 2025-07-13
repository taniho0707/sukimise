CREATE TABLE discord_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discord_id VARCHAR(255) NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL,
    access_token TEXT,
    refresh_token TEXT,
    token_expiry TIMESTAMP,
    linked_at TIMESTAMP DEFAULT NOW(),
    last_used_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_discord_links_discord_id ON discord_links(discord_id);
CREATE INDEX idx_discord_links_user_id ON discord_links(user_id);
CREATE INDEX idx_discord_links_username ON discord_links(username);