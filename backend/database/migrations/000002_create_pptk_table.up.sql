-- Create pptk table
CREATE TABLE pptk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nip VARCHAR(50) NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    unit_kerja_id UUID NOT NULL REFERENCES unit_kerja(id) ON DELETE RESTRICT,
    avatar_path VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pptk_nip ON pptk(nip);
CREATE INDEX idx_pptk_unit_kerja_id ON pptk(unit_kerja_id);
CREATE INDEX idx_pptk_is_active ON pptk(is_active);
