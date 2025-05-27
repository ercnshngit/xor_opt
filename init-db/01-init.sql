-- Create the matrix_records table
CREATE TABLE IF NOT EXISTS matrix_records (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    matrix_binary TEXT NOT NULL,
    matrix_hex TEXT NOT NULL,
    ham_xor_count INTEGER NOT NULL,
    boyar_xor_count INTEGER,
    boyar_depth INTEGER,
    boyar_program TEXT,
    paar_xor_count INTEGER,
    paar_program TEXT,
    slp_xor_count INTEGER,
    slp_program TEXT,
    matrix_hash TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_matrix_hash ON matrix_records(matrix_hash);
CREATE INDEX IF NOT EXISTS idx_title ON matrix_records(title);
CREATE INDEX IF NOT EXISTS idx_ham_xor_count ON matrix_records(ham_xor_count);
CREATE INDEX IF NOT EXISTS idx_boyar_xor_count ON matrix_records(boyar_xor_count);
CREATE INDEX IF NOT EXISTS idx_paar_xor_count ON matrix_records(paar_xor_count);
CREATE INDEX IF NOT EXISTS idx_slp_xor_count ON matrix_records(slp_xor_count);
CREATE INDEX IF NOT EXISTS idx_created_at ON matrix_records(created_at);

-- Create a function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create a trigger to automatically update the updated_at column
CREATE TRIGGER update_matrix_records_updated_at 
    BEFORE UPDATE ON matrix_records 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column(); 