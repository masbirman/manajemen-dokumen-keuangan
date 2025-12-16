-- Add nomor_dokumen and tanggal_dokumen columns to dokumen table
ALTER TABLE dokumen ADD COLUMN IF NOT EXISTS nomor_dokumen VARCHAR(255);
ALTER TABLE dokumen ADD COLUMN IF NOT EXISTS tanggal_dokumen DATE;

-- Create index for nomor_dokumen
CREATE INDEX IF NOT EXISTS idx_dokumen_nomor_dokumen ON dokumen(nomor_dokumen);
CREATE INDEX IF NOT EXISTS idx_dokumen_tanggal_dokumen ON dokumen(tanggal_dokumen);
