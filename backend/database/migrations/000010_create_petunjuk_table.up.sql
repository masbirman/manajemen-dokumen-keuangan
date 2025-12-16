-- Create petunjuk table for storing guidance/information
CREATE TABLE petunjuk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    judul VARCHAR(255) NOT NULL,
    konten TEXT NOT NULL,
    halaman VARCHAR(100) NOT NULL, -- 'input_dokumen', 'list_dokumen', etc.
    urutan INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_petunjuk_halaman ON petunjuk(halaman);
CREATE INDEX idx_petunjuk_is_active ON petunjuk(is_active);
CREATE INDEX idx_petunjuk_urutan ON petunjuk(urutan);
