-- Create matrix_records table
CREATE TABLE IF NOT EXISTS matrix_records (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
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
    matrix_hash VARCHAR(32) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_matrix_records_hash ON matrix_records(matrix_hash);
CREATE INDEX IF NOT EXISTS idx_matrix_records_title ON matrix_records(title);
CREATE INDEX IF NOT EXISTS idx_matrix_records_ham_xor ON matrix_records(ham_xor_count);
CREATE INDEX IF NOT EXISTS idx_matrix_records_boyar_xor ON matrix_records(boyar_xor_count);
CREATE INDEX IF NOT EXISTS idx_matrix_records_paar_xor ON matrix_records(paar_xor_count);
CREATE INDEX IF NOT EXISTS idx_matrix_records_slp_xor ON matrix_records(slp_xor_count);
CREATE INDEX IF NOT EXISTS idx_matrix_records_created_at ON matrix_records(created_at); 