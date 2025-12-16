-- Add jabatan column to pptk table
ALTER TABLE pptk ADD COLUMN IF NOT EXISTS jabatan VARCHAR(255);
