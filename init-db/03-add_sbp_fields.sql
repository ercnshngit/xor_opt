-- Add SBP algorithm fields to matrix_records table
ALTER TABLE matrix_records 
ADD COLUMN sbp_xor_count INTEGER,
ADD COLUMN sbp_depth INTEGER,
ADD COLUMN sbp_program TEXT;

-- Add SBP algorithm fields to inverse_pairs table
ALTER TABLE inverse_pairs 
ADD COLUMN sbp_xor_count INTEGER,
ADD COLUMN sbp_depth INTEGER,
ADD COLUMN sbp_program TEXT; 