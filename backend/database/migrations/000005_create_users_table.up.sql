-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('super_admin', 'admin', 'operator')),
    unit_kerja_id UUID REFERENCES unit_kerja(id) ON DELETE SET NULL,
    pptk_id UUID REFERENCES pptk(id) ON DELETE SET NULL,
    avatar_path VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_unit_kerja_id ON users(unit_kerja_id);
CREATE INDEX idx_users_pptk_id ON users(pptk_id);
CREATE INDEX idx_users_is_active ON users(is_active);
