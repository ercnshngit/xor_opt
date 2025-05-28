-- Create the matrix_records table
CREATE TABLE IF NOT EXISTS matrix_records (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    group_name TEXT,
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
    smallest_xor INTEGER,
    matrix_hash TEXT UNIQUE NOT NULL,
    inverse_matrix_id INTEGER,
    inverse_matrix_hash TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_matrix_hash ON matrix_records(matrix_hash);
CREATE INDEX IF NOT EXISTS idx_title ON matrix_records(title);
CREATE INDEX IF NOT EXISTS idx_group_name ON matrix_records(group_name);
CREATE INDEX IF NOT EXISTS idx_ham_xor_count ON matrix_records(ham_xor_count);
CREATE INDEX IF NOT EXISTS idx_boyar_xor_count ON matrix_records(boyar_xor_count);
CREATE INDEX IF NOT EXISTS idx_paar_xor_count ON matrix_records(paar_xor_count);
CREATE INDEX IF NOT EXISTS idx_slp_xor_count ON matrix_records(slp_xor_count);
CREATE INDEX IF NOT EXISTS idx_smallest_xor ON matrix_records(smallest_xor);
CREATE INDEX IF NOT EXISTS idx_created_at ON matrix_records(created_at);
CREATE INDEX IF NOT EXISTS idx_inverse_matrix_id ON matrix_records(inverse_matrix_id);

-- Composite indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_smallest_xor_created_at ON matrix_records(smallest_xor ASC, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_ham_xor_created_at ON matrix_records(ham_xor_count ASC, created_at DESC);

-- Index for title search (case insensitive)
CREATE INDEX IF NOT EXISTS idx_title_lower ON matrix_records(LOWER(title));

-- Partial indexes for algorithm results
CREATE INDEX IF NOT EXISTS idx_has_boyar ON matrix_records(id) WHERE boyar_xor_count IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_has_paar ON matrix_records(id) WHERE paar_xor_count IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_has_slp ON matrix_records(id) WHERE slp_xor_count IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_missing_algorithms ON matrix_records(id) WHERE boyar_xor_count IS NULL OR paar_xor_count IS NULL OR slp_xor_count IS NULL;

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