-- Performance optimization migration script
-- This script adds missing columns and indexes for better performance

-- Add missing columns if they don't exist
DO $$ 
BEGIN
    -- Add group_name column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'matrix_records' AND column_name = 'group_name') THEN
        ALTER TABLE matrix_records ADD COLUMN group_name TEXT;
        UPDATE matrix_records SET group_name = 'default' WHERE group_name IS NULL;
    END IF;
    
    -- Add smallest_xor column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'matrix_records' AND column_name = 'smallest_xor') THEN
        ALTER TABLE matrix_records ADD COLUMN smallest_xor INTEGER;
        -- Calculate smallest_xor for existing records
        UPDATE matrix_records SET smallest_xor = LEAST(
            COALESCE(boyar_xor_count, 999999),
            COALESCE(paar_xor_count, 999999),
            COALESCE(slp_xor_count, 999999),
            ham_xor_count
        ) WHERE smallest_xor IS NULL;
    END IF;
    
    -- Add inverse_matrix_id column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'matrix_records' AND column_name = 'inverse_matrix_id') THEN
        ALTER TABLE matrix_records ADD COLUMN inverse_matrix_id INTEGER;
    END IF;
    
    -- Add inverse_matrix_hash column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'matrix_records' AND column_name = 'inverse_matrix_hash') THEN
        ALTER TABLE matrix_records ADD COLUMN inverse_matrix_hash TEXT;
    END IF;
END $$;

-- Create performance indexes if they don't exist
CREATE INDEX IF NOT EXISTS idx_group_name ON matrix_records(group_name);
CREATE INDEX IF NOT EXISTS idx_smallest_xor ON matrix_records(smallest_xor);
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

-- Update statistics for better query planning
ANALYZE matrix_records;

-- Log completion
DO $$ 
BEGIN
    RAISE NOTICE 'Performance migration completed successfully';
END $$; 