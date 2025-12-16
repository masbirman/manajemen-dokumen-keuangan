-- Remove nomor_dokumen and tanggal_dokumen columns from dokumen table
DROP INDEX IF EXISTS idx_dokumen_nomor_dokumen;
DROP INDEX IF EXISTS idx_dokumen_tanggal_dokumen;
ALTER TABLE dokumen DROP COLUMN IF EXISTS nomor_dokumen;
ALTER TABLE dokumen DROP COLUMN IF EXISTS tanggal_dokumen;
