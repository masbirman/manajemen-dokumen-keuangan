-- Create unit_kerja table
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE unit_kerja (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode VARCHAR(50) NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_unit_kerja_kode ON unit_kerja(kode);
CREATE INDEX idx_unit_kerja_is_active ON unit_kerja(is_active);
