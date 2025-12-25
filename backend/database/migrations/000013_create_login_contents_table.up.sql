-- Login Contents table for scheduled login page content
CREATE TABLE IF NOT EXISTS login_contents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    image_width INTEGER DEFAULT 400,
    title_size INTEGER DEFAULT 28,
    desc_size INTEGER DEFAULT 16,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for active content lookup
CREATE INDEX IF NOT EXISTS idx_login_contents_active ON login_contents(is_active, start_date, end_date);
