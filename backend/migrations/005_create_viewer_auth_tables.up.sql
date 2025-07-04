-- Create viewer password settings table
CREATE TABLE viewer_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    password_hash VARCHAR(255) NOT NULL,
    session_duration_days INTEGER NOT NULL DEFAULT 7,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create viewer login history table
CREATE TABLE viewer_login_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address INET,
    user_agent TEXT,
    login_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    session_token VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX idx_viewer_login_history_login_time ON viewer_login_history(login_time);
CREATE INDEX idx_viewer_login_history_session_token ON viewer_login_history(session_token);
CREATE INDEX idx_viewer_login_history_expires_at ON viewer_login_history(expires_at);

-- Insert default viewer settings (password: viewer123)
INSERT INTO viewer_settings (password_hash, session_duration_days) VALUES (
    '$2a$10$vPZxOoHW8tRYvBhDHN4yBOmJQfgVzv7rVHvLFxEGIGsNTVcBjJqhS', -- bcrypt hash of 'viewer123'
    7
);