-- Create jenis_dokumen table
CREATE TABLE jenis_dokumen (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kode VARCHAR(50) NOT NULL UNIQUE,
    nama VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_jenis_dokumen_kode ON jenis_dokumen(kode);
CREATE INDEX idx_jenis_dokumen_is_active ON jenis_dokumen(is_active);
