-- Create dokumen table
CREATE TABLE dokumen (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_kerja_id UUID NOT NULL REFERENCES unit_kerja(id) ON DELETE RESTRICT,
    pptk_id UUID NOT NULL REFERENCES pptk(id) ON DELETE RESTRICT,
    jenis_dokumen_id UUID NOT NULL REFERENCES jenis_dokumen(id) ON DELETE RESTRICT,
    sumber_dana_id UUID NOT NULL REFERENCES sumber_dana(id) ON DELETE RESTRICT,
    nilai DECIMAL(15, 2) NOT NULL,
    uraian TEXT NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_dokumen_unit_kerja_id ON dokumen(unit_kerja_id);
CREATE INDEX idx_dokumen_pptk_id ON dokumen(pptk_id);
CREATE INDEX idx_dokumen_jenis_dokumen_id ON dokumen(jenis_dokumen_id);
CREATE INDEX idx_dokumen_sumber_dana_id ON dokumen(sumber_dana_id);
CREATE INDEX idx_dokumen_created_by ON dokumen(created_by);
CREATE INDEX idx_dokumen_created_at ON dokumen(created_at);
