package main

import (
	"bufio"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// MatrixRecord represents a matrix record in the database
type MatrixRecord struct {
	ID                  int       `json:"id"`
	Title               string    `json:"title"`
	Group               string    `json:"group"`
	MatrixBinary        string    `json:"matrix_binary"`
	MatrixHex           string    `json:"matrix_hex"`
	HamXorCount         int       `json:"ham_xor_count"`
	BoyarXorCount       *int      `json:"boyar_xor_count,omitempty"`
	BoyarDepth          *int      `json:"boyar_depth,omitempty"`
	BoyarProgram        *string   `json:"boyar_program,omitempty"`
	PaarXorCount        *int      `json:"paar_xor_count,omitempty"`
	PaarProgram         *string   `json:"paar_program,omitempty"`
	SlpXorCount         *int      `json:"slp_xor_count,omitempty"`
	SlpProgram          *string   `json:"slp_program,omitempty"`
	SmallestXor         *int      `json:"smallest_xor,omitempty"`
	MatrixHash          string    `json:"matrix_hash"`
	InverseMatrixID     *int      `json:"inverse_matrix_id,omitempty"`
	InverseMatrixHash   *string   `json:"inverse_matrix_hash,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Database represents the PostgreSQL database
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{db: db}
	return database, nil
}

// matrixToHex converts a binary matrix to hex representation
func matrixToHex(matrix Matrix) string {
	var hexStrings []string
	for _, row := range matrix {
		binaryStr := strings.Join(row, "")
		// Pad to multiple of 4 bits for hex conversion
		for len(binaryStr)%4 != 0 {
			binaryStr += "0"
		}
		
		var hexStr string
		for i := 0; i < len(binaryStr); i += 4 {
			chunk := binaryStr[i:i+4]
			val, _ := strconv.ParseInt(chunk, 2, 64)
			hexStr += fmt.Sprintf("%X", val)
		}
		hexStrings = append(hexStrings, hexStr)
	}
	return strings.Join(hexStrings, ",")
}

// matrixToBinary converts matrix to string representation
func matrixToBinary(matrix Matrix) string {
	var rows []string
	for _, row := range matrix {
		rows = append(rows, "["+strings.Join(row, " ")+"]")
	}
	return strings.Join(rows, "\n")
}

// calculateMatrixHash creates a unique hash for the matrix
func calculateMatrixHash(matrix Matrix) string {
	matrixStr := matrixToBinary(matrix)
	hash := md5.Sum([]byte(matrixStr))
	return hex.EncodeToString(hash[:])
}

// calculateHammingXOR calculates the Hamming XOR count for a matrix
func calculateHammingXOR(matrix Matrix) int {
	count := 0
	for _, row := range matrix {
		for _, cell := range row {
			if cell == "1" {
				count++
			}
		}
	}
	// Ham XOR = toplam 1'ler - sütun sayısı
	if len(matrix) > 0 {
		count -= len(matrix[0])
	}
	return count
}

// SaveMatrix saves a matrix to the database
func (d *Database) SaveMatrix(title string, matrix Matrix, group string) (*MatrixRecord, error) {
	matrixHash := calculateMatrixHash(matrix)
	
	// Check if matrix already exists
	existing, err := d.GetMatrixByHash(matrixHash)
	if err == nil && existing != nil {
		return existing, nil
	}

	matrixBinary := matrixToBinary(matrix)
	matrixHex := matrixToHex(matrix)
	hamXorCount := calculateHammingXOR(matrix)

	query := `
	INSERT INTO matrix_records (title, group_name, matrix_binary, matrix_hex, ham_xor_count, matrix_hash)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	var id int
	err = d.db.QueryRow(query, title, group, matrixBinary, matrixHex, hamXorCount, matrixHash).Scan(&id)
	if err != nil {
		return nil, err
	}

	return d.GetMatrixByID(id)
}

// UpdateMatrixResults updates the algorithm results for a matrix
func (d *Database) UpdateMatrixResults(id int, boyarResult, paarResult, slpResult *AlgResult) error {
	// Calculate smallest XOR value
	var smallestXor *int
	var xorValues []int
	
	if boyarResult != nil {
		xorValues = append(xorValues, boyarResult.XorCount)
	}
	if paarResult != nil {
		xorValues = append(xorValues, paarResult.XorCount)
	}
	if slpResult != nil {
		xorValues = append(xorValues, slpResult.XorCount)
	}
	
	if len(xorValues) > 0 {
		min := xorValues[0]
		for _, val := range xorValues {
			if val < min {
				min = val
			}
		}
		smallestXor = &min
	}

	query := `
	UPDATE matrix_records 
	SET boyar_xor_count = $1, boyar_depth = $2, boyar_program = $3,
	    paar_xor_count = $4, paar_program = $5,
	    slp_xor_count = $6, slp_program = $7,
	    smallest_xor = $8,
	    updated_at = CURRENT_TIMESTAMP
	WHERE id = $9
	`

	var boyarXor, boyarDepth *int
	var boyarProgram *string
	if boyarResult != nil {
		boyarXor = &boyarResult.XorCount
		boyarDepth = &boyarResult.Depth
		programJson, _ := json.Marshal(boyarResult.Program)
		programStr := string(programJson)
		boyarProgram = &programStr
	}

	var paarXor *int
	var paarProgram *string
	if paarResult != nil {
		paarXor = &paarResult.XorCount
		programJson, _ := json.Marshal(paarResult.Program)
		programStr := string(programJson)
		paarProgram = &programStr
	}

	var slpXor *int
	var slpProgram *string
	if slpResult != nil {
		slpXor = &slpResult.XorCount
		programJson, _ := json.Marshal(slpResult.Program)
		programStr := string(programJson)
		slpProgram = &programStr
	}

	_, err := d.db.Exec(query, boyarXor, boyarDepth, boyarProgram, paarXor, paarProgram, slpXor, slpProgram, smallestXor, id)
	return err
}

// GetMatrixByID retrieves a matrix by its ID
func (d *Database) GetMatrixByID(id int) (*MatrixRecord, error) {
	query := `
	SELECT id, title, group_name, matrix_binary, matrix_hex, ham_xor_count, smallest_xor,
	       boyar_xor_count, boyar_depth, boyar_program,
	       paar_xor_count, paar_program, slp_xor_count, slp_program,
	       matrix_hash, inverse_matrix_id, inverse_matrix_hash, created_at, updated_at
	FROM matrix_records WHERE id = $1
	`
	
	row := d.db.QueryRow(query, id)
	return d.scanMatrixRecord(row)
}

// GetMatrixByHash retrieves a matrix by its hash
func (d *Database) GetMatrixByHash(hash string) (*MatrixRecord, error) {
	query := `
	SELECT id, title, group_name, matrix_binary, matrix_hex, ham_xor_count, smallest_xor,
	       boyar_xor_count, boyar_depth, boyar_program,
	       paar_xor_count, paar_program, slp_xor_count, slp_program,
	       matrix_hash, inverse_matrix_id, inverse_matrix_hash, created_at, updated_at
	FROM matrix_records WHERE matrix_hash = $1
	`
	
	row := d.db.QueryRow(query, hash)
	return d.scanMatrixRecord(row)
}

// GetMatrices retrieves matrices with pagination and filtering
func (d *Database) GetMatrices(page, limit int, titleFilter string, hamXorMin, hamXorMax, boyarXorMin, boyarXorMax, paarXorMin, paarXorMax, slpXorMin, slpXorMax *int) ([]*MatrixRecord, int, error) {
	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1
	
	if titleFilter != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", argIndex))
		args = append(args, "%"+titleFilter+"%")
		argIndex++
	}

	if hamXorMin != nil {
		conditions = append(conditions, fmt.Sprintf("ham_xor_count >= $%d", argIndex))
		args = append(args, *hamXorMin)
		argIndex++
	}

	if hamXorMax != nil {
		conditions = append(conditions, fmt.Sprintf("ham_xor_count <= $%d", argIndex))
		args = append(args, *hamXorMax)
		argIndex++
	}

	if boyarXorMin != nil {
		conditions = append(conditions, fmt.Sprintf("boyar_xor_count IS NOT NULL AND boyar_xor_count >= $%d", argIndex))
		args = append(args, *boyarXorMin)
		argIndex++
	}

	if boyarXorMax != nil {
		conditions = append(conditions, fmt.Sprintf("boyar_xor_count IS NOT NULL AND boyar_xor_count <= $%d", argIndex))
		args = append(args, *boyarXorMax)
		argIndex++
	}

	if paarXorMin != nil {
		conditions = append(conditions, fmt.Sprintf("paar_xor_count IS NOT NULL AND paar_xor_count >= $%d", argIndex))
		args = append(args, *paarXorMin)
		argIndex++
	}

	if paarXorMax != nil {
		conditions = append(conditions, fmt.Sprintf("paar_xor_count IS NOT NULL AND paar_xor_count <= $%d", argIndex))
		args = append(args, *paarXorMax)
		argIndex++
	}

	if slpXorMin != nil {
		conditions = append(conditions, fmt.Sprintf("slp_xor_count IS NOT NULL AND slp_xor_count >= $%d", argIndex))
		args = append(args, *slpXorMin)
		argIndex++
	}

	if slpXorMax != nil {
		conditions = append(conditions, fmt.Sprintf("slp_xor_count IS NOT NULL AND slp_xor_count <= $%d", argIndex))
		args = append(args, *slpXorMax)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM matrix_records %s", whereClause)
	var total int
	err := d.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * limit
	query := fmt.Sprintf(`
	SELECT id, title, group_name, matrix_binary, matrix_hex, ham_xor_count, smallest_xor,
	       boyar_xor_count, boyar_depth, boyar_program,
	       paar_xor_count, paar_program, slp_xor_count, slp_program,
	       matrix_hash, inverse_matrix_id, inverse_matrix_hash, created_at, updated_at
	FROM matrix_records %s
	ORDER BY COALESCE(smallest_xor, ham_xor_count) ASC, created_at DESC
	LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)
	
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matrices []*MatrixRecord
	for rows.Next() {
		matrix, err := d.scanMatrixRecord(rows)
		if err != nil {
			return nil, 0, err
		}
		matrices = append(matrices, matrix)
	}

	return matrices, total, nil
}

// scanMatrixRecord scans a row into a MatrixRecord
func (d *Database) scanMatrixRecord(scanner interface{}) (*MatrixRecord, error) {
	var record MatrixRecord
	var groupName sql.NullString
	var smallestXor, boyarXor, boyarDepth, paarXor, slpXor, inverseMatrixID sql.NullInt64
	var boyarProgram, paarProgram, slpProgram, inverseMatrixHash sql.NullString

	var err error
	switch s := scanner.(type) {
	case *sql.Row:
		err = s.Scan(&record.ID, &record.Title, &groupName, &record.MatrixBinary, &record.MatrixHex,
			&record.HamXorCount, &smallestXor, &boyarXor, &boyarDepth, &boyarProgram,
			&paarXor, &paarProgram, &slpXor, &slpProgram,
			&record.MatrixHash, &inverseMatrixID, &inverseMatrixHash, &record.CreatedAt, &record.UpdatedAt)
	case *sql.Rows:
		err = s.Scan(&record.ID, &record.Title, &groupName, &record.MatrixBinary, &record.MatrixHex,
			&record.HamXorCount, &smallestXor, &boyarXor, &boyarDepth, &boyarProgram,
			&paarXor, &paarProgram, &slpXor, &slpProgram,
			&record.MatrixHash, &inverseMatrixID, &inverseMatrixHash, &record.CreatedAt, &record.UpdatedAt)
	default:
		return nil, fmt.Errorf("unsupported scanner type")
	}

	if err != nil {
		return nil, err
	}

	// Handle nullable fields
	if groupName.Valid {
		record.Group = groupName.String
	}
	if smallestXor.Valid {
		val := int(smallestXor.Int64)
		record.SmallestXor = &val
	}
	if boyarXor.Valid {
		val := int(boyarXor.Int64)
		record.BoyarXorCount = &val
	}
	if boyarDepth.Valid {
		val := int(boyarDepth.Int64)
		record.BoyarDepth = &val
	}
	if boyarProgram.Valid {
		record.BoyarProgram = &boyarProgram.String
	}
	if paarXor.Valid {
		val := int(paarXor.Int64)
		record.PaarXorCount = &val
	}
	if paarProgram.Valid {
		record.PaarProgram = &paarProgram.String
	}
	if slpXor.Valid {
		val := int(slpXor.Int64)
		record.SlpXorCount = &val
	}
	if slpProgram.Valid {
		record.SlpProgram = &slpProgram.String
	}
	if inverseMatrixID.Valid {
		val := int(inverseMatrixID.Int64)
		record.InverseMatrixID = &val
	}
	if inverseMatrixHash.Valid {
		record.InverseMatrixHash = &inverseMatrixHash.String
	}

	return &record, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// GetMatrixCount returns the total number of matrices in the database
func (d *Database) GetMatrixCount() (int, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM matrix_records").Scan(&count)
	return count, err
}

// ImportMatricesFromFiles imports matrices from the matrices-data directory
func (d *Database) ImportMatricesFromFiles(dataPath string) error {
	importStartTime := time.Now()
	log.Printf("🚀 [IMPORT] Matrices-data klasöründen matrisler import ediliyor: %s", dataPath)
	
	// Get list of .txt files in the directory
	fileListStartTime := time.Now()
	files, err := filepath.Glob(filepath.Join(dataPath, "*.txt"))
	fileListDuration := time.Since(fileListStartTime)
	log.Printf("📁 [IMPORT] Dosya listesi alındı (%v): %d dosya bulundu", fileListDuration, len(files))
	
	if err != nil {
		return fmt.Errorf("dosya listesi alınamadı: %v", err)
	}

	if len(files) == 0 {
		log.Printf("⚠️  [IMPORT] Matrices-data klasöründe .txt dosyası bulunamadı")
		return nil
	}

	log.Printf("📋 [IMPORT] %d dosya bulundu, import işlemi başlıyor...", len(files))

	totalImported := 0
	for i, filePath := range files {
		fileName := filepath.Base(filePath)
		fileStartTime := time.Now()
		log.Printf("📄 [IMPORT] Dosya işleniyor (%d/%d): %s", i+1, len(files), fileName)
		
		count, err := d.importMatricesFromFile(filePath)
		fileDuration := time.Since(fileStartTime)
		
		if err != nil {
			log.Printf("❌ [IMPORT] %s dosyası işlenirken hata oluştu (%v): %v", fileName, fileDuration, err)
			continue
		}
		
		totalImported += count
		log.Printf("✅ [IMPORT] %s dosyasından %d matris import edildi (%v)", fileName, count, fileDuration)
	}

	totalImportDuration := time.Since(importStartTime)
	log.Printf("🎉 [IMPORT] Import işlemi tamamlandı (Toplam süre: %v). Toplam %d matris import edildi.", totalImportDuration, totalImported)
	return nil
}

// importMatricesFromFile imports matrices from a single file
func (d *Database) importMatricesFromFile(filePath string) (int, error) {
	fileStartTime := time.Now()
	fileName := filepath.Base(filePath)
	log.Printf("📖 [FILE] Dosya okunuyor: %s", fileName)
	
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Extract filename without extension as group name
	groupName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	scanner := bufio.NewScanner(file)
	var currentMatrix [][]string
	var currentTitle string
	importedCount := 0
	lineNumber := 0
	matrixCount := 0

	parseStartTime := time.Now()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNumber++

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for separator
		if strings.HasPrefix(line, "------------------------------") {
			// Process current matrix if we have one
			if len(currentMatrix) > 0 && currentTitle != "" {
				matrixStartTime := time.Now()
				err := d.saveMatrixFromImport(currentTitle, currentMatrix, groupName)
				matrixDuration := time.Since(matrixStartTime)
				
				if err != nil {
					log.Printf("❌ [FILE] Matris kaydedilemedi (%v) (satır %d): %v", matrixDuration, lineNumber, err)
				} else {
					importedCount++
					log.Printf("💾 [FILE] Matris kaydedildi (%v): %s", matrixDuration, currentTitle)
				}
			}
			
			// Reset for next matrix
			currentMatrix = [][]string{}
			currentTitle = ""
			matrixCount++
			continue
		}

		// Check if this is a title line (contains "matrisi" and ends with ":")
		if strings.Contains(line, "matrisi") && strings.HasSuffix(line, ":") {
			currentTitle = strings.TrimSuffix(line, ":")
			continue
		}

		// Check if this is a matrix row (starts with [ and ends with ])
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Parse matrix row
			rowStr := strings.Trim(line, "[]")
			elements := strings.Fields(rowStr)
			if len(elements) > 0 {
				currentMatrix = append(currentMatrix, elements)
			}
			continue
		}

		// Skip other lines (like "HamXOR Sayisi:" etc.)
	}

	// Process the last matrix if exists
	if len(currentMatrix) > 0 && currentTitle != "" {
		matrixStartTime := time.Now()
		err := d.saveMatrixFromImport(currentTitle, currentMatrix, groupName)
		matrixDuration := time.Since(matrixStartTime)
		
		if err != nil {
			log.Printf("❌ [FILE] Son matris kaydedilemedi (%v): %v", matrixDuration, err)
		} else {
			importedCount++
			log.Printf("💾 [FILE] Son matris kaydedildi (%v): %s", matrixDuration, currentTitle)
		}
		matrixCount++
	}

	parseDuration := time.Since(parseStartTime)
	totalFileDuration := time.Since(fileStartTime)

	if err := scanner.Err(); err != nil {
		return importedCount, err
	}

	log.Printf("📊 [FILE] %s tamamlandı - %d satır okundu, %d matris bulundu, %d matris import edildi (Parse: %v, Toplam: %v)", 
		fileName, lineNumber, matrixCount, importedCount, parseDuration, totalFileDuration)

	return importedCount, nil
}

// saveMatrixFromImport saves a matrix during import process
func (d *Database) saveMatrixFromImport(title string, matrix [][]string, group string) error {
	startTime := time.Now()
	log.Printf("📊 [IMPORT] Matris işleme başlıyor: %s", title)
	
	// Check if matrix already exists by hash
	hashStartTime := time.Now()
	matrixHash := calculateMatrixHash(matrix)
	existing, err := d.GetMatrixByHash(matrixHash)
	hashDuration := time.Since(hashStartTime)
	log.Printf("⏱️  [IMPORT] Hash kontrolü tamamlandı (%v): %s", hashDuration, title)
	
	if err == nil && existing != nil {
		// Matrix already exists, skip
		log.Printf("⏭️  [IMPORT] Matris zaten mevcut, atlanıyor: %s", title)
		return nil
	}

	// Save the matrix
	saveStartTime := time.Now()
	savedMatrix, err := d.SaveMatrix(title, matrix, group)
	saveDuration := time.Since(saveStartTime)
	log.Printf("💾 [IMPORT] Matris veritabanına kaydedildi (%v): %s", saveDuration, title)
	
	if err != nil {
		return err
	}

	// Calculate algorithms for the newly saved matrix
	log.Printf("🧮 [IMPORT] Algoritma hesaplamaları başlıyor: %s", title)
	
	// Run algorithms in background
	go func() {
		algorithmStartTime := time.Now()
		
		// Calculate Boyar SLP
		boyarStartTime := time.Now()
		boyarResult, err := runBoyarSLP(matrix)
		boyarDuration := time.Since(boyarStartTime)
		if err != nil {
			log.Printf("❌ [BOYAR] %s için Boyar SLP hesaplanamadı (%v): %v", title, boyarDuration, err)
		} else {
			log.Printf("✅ [BOYAR] %s için Boyar SLP tamamlandı (%v) - XOR: %d", title, boyarDuration, boyarResult.XorCount)
		}

		// Calculate Paar Algorithm
		paarStartTime := time.Now()
		paarResult, err := runPaarAlgorithm(matrix)
		paarDuration := time.Since(paarStartTime)
		if err != nil {
			log.Printf("❌ [PAAR] %s için Paar algoritması hesaplanamadı (%v): %v", title, paarDuration, err)
		} else {
			log.Printf("✅ [PAAR] %s için Paar algoritması tamamlandı (%v) - XOR: %d", title, paarDuration, paarResult.XorCount)
		}

		// Calculate SLP Heuristic
		slpStartTime := time.Now()
		slpResult, err := runSLPHeuristic(matrix)
		slpDuration := time.Since(slpStartTime)
		if err != nil {
			log.Printf("❌ [SLP] %s için SLP Heuristic hesaplanamadı (%v): %v", title, slpDuration, err)
		} else {
			log.Printf("✅ [SLP] %s için SLP Heuristic tamamlandı (%v) - XOR: %d", title, slpDuration, slpResult.XorCount)
		}

		// Update matrix with results
		updateStartTime := time.Now()
		err = d.UpdateMatrixResults(savedMatrix.ID, boyarResult, paarResult, slpResult)
		updateDuration := time.Since(updateStartTime)
		
		totalAlgorithmDuration := time.Since(algorithmStartTime)
		
		if err != nil {
			log.Printf("❌ [UPDATE] %s için sonuçlar kaydedilemedi (%v): %v", title, updateDuration, err)
		} else {
			log.Printf("✅ [UPDATE] %s için sonuçlar kaydedildi (%v)", title, updateDuration)
			log.Printf("🎯 [TOPLAM] %s için tüm algoritmalar tamamlandı (Toplam: %v, Boyar: %v, Paar: %v, SLP: %v)", 
				title, totalAlgorithmDuration, boyarDuration, paarDuration, slpDuration)
		}
	}()

	totalDuration := time.Since(startTime)
	log.Printf("📈 [IMPORT] Matris işleme tamamlandı (%v): %s", totalDuration, title)
	return nil
}

// GetAllMatrixHashes returns all matrix hashes from database
func (d *Database) GetAllMatrixHashes() (map[string]bool, error) {
	query := "SELECT matrix_hash FROM matrix_records"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hashes := make(map[string]bool)
	for rows.Next() {
		var hash string
		if err := rows.Scan(&hash); err != nil {
			return nil, err
		}
		hashes[hash] = true
	}

	return hashes, nil
}

// GetFileMatrixHashes returns all matrix hashes from files
func (d *Database) GetFileMatrixHashes(dataPath string) (map[string]bool, error) {
	files, err := filepath.Glob(filepath.Join(dataPath, "*.txt"))
	if err != nil {
		return nil, err
	}

	hashes := make(map[string]bool)
	for _, filePath := range files {
		fileHashes, err := d.getMatrixHashesFromFile(filePath)
		if err != nil {
			log.Printf("HATA: %s dosyasındaki hash'ler alınamadı: %v", filepath.Base(filePath), err)
			continue
		}
		for hash := range fileHashes {
			hashes[hash] = true
		}
	}

	return hashes, nil
}

// getMatrixHashesFromFile extracts matrix hashes from a single file
func (d *Database) getMatrixHashesFromFile(filePath string) (map[string]bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentMatrix [][]string
	var currentTitle string
	hashes := make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for separator
		if strings.HasPrefix(line, "------------------------------") {
			// Process current matrix if we have one
			if len(currentMatrix) > 0 && currentTitle != "" {
				hash := calculateMatrixHash(currentMatrix)
				hashes[hash] = true
			}
			
			// Reset for next matrix
			currentMatrix = [][]string{}
			currentTitle = ""
			continue
		}

		// Check if this is a title line (contains "matrisi" and ends with ":")
		if strings.Contains(line, "matrisi") && strings.HasSuffix(line, ":") {
			currentTitle = strings.TrimSuffix(line, ":")
			continue
		}

		// Check if this is a matrix row (starts with [ and ends with ])
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// Parse matrix row
			rowStr := strings.Trim(line, "[]")
			elements := strings.Fields(rowStr)
			if len(elements) > 0 {
				currentMatrix = append(currentMatrix, elements)
			}
			continue
		}
	}

	// Process the last matrix if exists
	if len(currentMatrix) > 0 && currentTitle != "" {
		hash := calculateMatrixHash(currentMatrix)
		hashes[hash] = true
	}

	return hashes, scanner.Err()
}

// GetMatricesWithoutAlgorithms returns matrices that don't have algorithm results
func (d *Database) GetMatricesWithoutAlgorithms(limit int) ([]*MatrixRecord, error) {
	query := `
	SELECT id, title, matrix_binary, matrix_hex, ham_xor_count, 
	       boyar_xor_count, boyar_depth, boyar_program,
	       paar_xor_count, paar_program, slp_xor_count, slp_program,
	       matrix_hash, created_at, updated_at
	FROM matrix_records 
	WHERE (boyar_xor_count IS NULL OR paar_xor_count IS NULL OR slp_xor_count IS NULL)
	ORDER BY created_at ASC
	LIMIT $1
	`
	
	rows, err := d.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matrices []*MatrixRecord
	for rows.Next() {
		matrix, err := d.scanMatrixRecord(rows)
		if err != nil {
			return nil, err
		}
		matrices = append(matrices, matrix)
	}

	return matrices, nil
}

// Algorithm runner functions for import process

// runBoyarSLP runs the Boyar SLP algorithm on a matrix
func runBoyarSLP(matrix [][]string) (*AlgResult, error) {
	boyar := NewBoyarSLP(10) // depth limit
	
	err := boyar.ReadTargetMatrix(matrix)
	if err != nil {
		return nil, err
	}
	
	err = boyar.InitBase()
	if err != nil {
		return nil, err
	}
	
	result, err := boyar.Solve(matrix)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// runPaarAlgorithm runs the Paar algorithm on a matrix
func runPaarAlgorithm(matrix [][]string) (*AlgResult, error) {
	paar := NewPaarAlgorithm()
	
	err := paar.ReadTargetMatrix(matrix)
	if err != nil {
		return nil, err
	}
	
	err = paar.InitBase()
	if err != nil {
		return nil, err
	}
	
	result, err := paar.Solve(matrix)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// runSLPHeuristic runs the SLP Heuristic algorithm on a matrix
func runSLPHeuristic(matrix [][]string) (*AlgResult, error) {
	slp := NewSLPHeuristic()
	
	err := slp.ReadTargetMatrix(matrix)
	if err != nil {
		return nil, err
	}
	
	err = slp.InitBase()
	if err != nil {
		return nil, err
	}
	
	result, err := slp.Solve(matrix)
	if err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Global database instance
var db *Database

// createTables creates the necessary database tables
func createTables(database *sql.DB) error {
	// First, try to add new columns if they don't exist
	migrationSQL := `
	-- Add group_name column if it doesn't exist
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='matrix_records' AND column_name='group_name') THEN
			ALTER TABLE matrix_records ADD COLUMN group_name VARCHAR(255);
		END IF;
	END $$;

	-- Add smallest_xor column if it doesn't exist
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='matrix_records' AND column_name='smallest_xor') THEN
			ALTER TABLE matrix_records ADD COLUMN smallest_xor INTEGER;
		END IF;
	END $$;

	-- Add inverse_matrix_id column if it doesn't exist
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='matrix_records' AND column_name='inverse_matrix_id') THEN
			ALTER TABLE matrix_records ADD COLUMN inverse_matrix_id INTEGER;
		END IF;
	END $$;

	-- Add inverse_matrix_hash column if it doesn't exist
	DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='matrix_records' AND column_name='inverse_matrix_hash') THEN
			ALTER TABLE matrix_records ADD COLUMN inverse_matrix_hash VARCHAR(32);
		END IF;
	END $$;
	`

	_, err := database.Exec(migrationSQL)
	if err != nil {
		log.Printf("Migration hatası (devam ediliyor): %v", err)
	}

	createTableSQL := `
	-- Create matrix_records table
	CREATE TABLE IF NOT EXISTS matrix_records (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		group_name VARCHAR(255),
		matrix_binary TEXT NOT NULL,
		matrix_hex TEXT NOT NULL,
		ham_xor_count INTEGER NOT NULL,
		smallest_xor INTEGER,
		boyar_xor_count INTEGER,
		boyar_depth INTEGER,
		boyar_program TEXT,
		paar_xor_count INTEGER,
		paar_program TEXT,
		slp_xor_count INTEGER,
		slp_program TEXT,
		matrix_hash VARCHAR(32) NOT NULL UNIQUE,
		inverse_matrix_id INTEGER,
		inverse_matrix_hash VARCHAR(32),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Create indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_matrix_records_hash ON matrix_records(matrix_hash);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_title ON matrix_records(title);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_group ON matrix_records(group_name);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_ham_xor ON matrix_records(ham_xor_count);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_smallest_xor ON matrix_records(smallest_xor);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_boyar_xor ON matrix_records(boyar_xor_count);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_paar_xor ON matrix_records(paar_xor_count);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_slp_xor ON matrix_records(slp_xor_count);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_inverse_id ON matrix_records(inverse_matrix_id);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_inverse_hash ON matrix_records(inverse_matrix_hash);
	CREATE INDEX IF NOT EXISTS idx_matrix_records_created_at ON matrix_records(created_at);
	`

	_, err = database.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("tablo oluşturma hatası: %v", err)
	}

	log.Printf("Veritabanı tabloları başarıyla oluşturuldu/kontrol edildi")
	
	// Update smallest_xor for existing records
	go func() {
		time.Sleep(2 * time.Second) // Wait for database to be ready
		updateSmallestXorForExistingRecords(database)
	}()
	
	return nil
}

// updateSmallestXorForExistingRecords updates smallest_xor for existing records
func updateSmallestXorForExistingRecords(database *sql.DB) {
	log.Printf("Mevcut kayıtlar için smallest_xor değerleri güncelleniyor...")
	
	query := `
	UPDATE matrix_records 
	SET smallest_xor = (
		SELECT MIN(xor_value) FROM (
			SELECT COALESCE(boyar_xor_count, 999999) as xor_value
			UNION ALL
			SELECT COALESCE(paar_xor_count, 999999) as xor_value
			UNION ALL
			SELECT COALESCE(slp_xor_count, 999999) as xor_value
		) AS xor_values
		WHERE xor_value < 999999
	)
	WHERE smallest_xor IS NULL 
	AND (boyar_xor_count IS NOT NULL OR paar_xor_count IS NOT NULL OR slp_xor_count IS NOT NULL)
	`
	
	result, err := database.Exec(query)
	if err != nil {
		log.Printf("Smallest XOR güncelleme hatası: %v", err)
		return
	}
	
	rowsAffected, _ := result.RowsAffected()
	log.Printf("✓ %d kayıt için smallest_xor değeri güncellendi", rowsAffected)
}

// InitDatabase initializes the database connection
func InitDatabase(config *Config) error {
	// Get database connection parameters from environment
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "xor_opt"
	}
	
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "xor_user"
	}
	
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "xor_password"
	}
	
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var err error
	db, err = NewDatabase(connStr)
	if err != nil {
		return fmt.Errorf("veritabanı bağlantısı kurulamadı: %v", err)
	}

	log.Printf("PostgreSQL veritabanına başarıyla bağlanıldı: %s:%s/%s", host, port, dbname)

	// Create tables if they don't exist
	err = createTables(db.db)
	if err != nil {
		return fmt.Errorf("veritabanı tabloları oluşturulamadı: %v", err)
	}

	// Check if we need to import matrices using hash comparison (only if enabled in config)
	if config != nil && config.Import.Enabled && config.Import.ProcessOnStart {
		go func() {
			time.Sleep(5 * time.Second) // Wait a bit for the application to fully start
			
			autoImportStartTime := time.Now()
			log.Printf("🔄 [AUTO-IMPORT] Otomatik import işlemi başlıyor...")
			
			dataPath := config.Import.DataDirectory
			if dataPath == "" {
				dataPath = "./matrices-data"
			}
			
			// Check for environment variable override
			if envPath := os.Getenv("MATRICES_DATA_PATH"); envPath != "" {
				dataPath = envPath
			}

			log.Printf("📂 [AUTO-IMPORT] Config'e göre otomatik import başlatılıyor: %s", dataPath)

			// Get hashes from database
			dbHashStartTime := time.Now()
			dbHashes, err := db.GetAllMatrixHashes()
			dbHashDuration := time.Since(dbHashStartTime)
			if err != nil {
				log.Printf("❌ [AUTO-IMPORT] Veritabanı hash'leri alınamadı (%v): %v", dbHashDuration, err)
				return
			}
			log.Printf("🗄️  [AUTO-IMPORT] Veritabanı hash'leri alındı (%v): %d hash", dbHashDuration, len(dbHashes))

			// Get hashes from files
			fileHashStartTime := time.Now()
			fileHashes, err := db.GetFileMatrixHashes(dataPath)
			fileHashDuration := time.Since(fileHashStartTime)
			if err != nil {
				log.Printf("❌ [AUTO-IMPORT] Dosya hash'leri alınamadı (%v): %v", fileHashDuration, err)
				return
			}
			log.Printf("📁 [AUTO-IMPORT] Dosya hash'leri alındı (%v): %d hash", fileHashDuration, len(fileHashes))

			log.Printf("📊 [AUTO-IMPORT] Veritabanında %d matris hash'i bulundu", len(dbHashes))
			log.Printf("📊 [AUTO-IMPORT] Dosyalarda %d matris hash'i bulundu", len(fileHashes))

			// Check for missing matrices
			compareStartTime := time.Now()
			missingCount := 0
			for hash := range fileHashes {
				if !dbHashes[hash] {
					missingCount++
				}
			}
			compareDuration := time.Since(compareStartTime)
			log.Printf("🔍 [AUTO-IMPORT] Hash karşılaştırması tamamlandı (%v): %d eksik matris", compareDuration, missingCount)

			if missingCount > 0 {
				log.Printf("🚀 [AUTO-IMPORT] Veritabanında %d eksik matris var, import işlemi başlatılıyor...", missingCount)
				importStartTime := time.Now()
				err := db.ImportMatricesFromFiles(dataPath)
				importDuration := time.Since(importStartTime)
				if err != nil {
					log.Printf("❌ [AUTO-IMPORT] Matris import işlemi başarısız (%v): %v", importDuration, err)
				} else {
					log.Printf("✅ [AUTO-IMPORT] Matris import işlemi tamamlandı (%v)", importDuration)
				}
			} else {
				log.Printf("✅ [AUTO-IMPORT] Tüm matrisler zaten veritabanında mevcut")
			}
			
			totalAutoImportDuration := time.Since(autoImportStartTime)
			log.Printf("🎯 [AUTO-IMPORT] Otomatik import işlemi tamamlandı (Toplam süre: %v)", totalAutoImportDuration)
		}()
	} else {
		log.Printf("Otomatik import devre dışı (enabled: %v, process_on_start: %v)", 
			config != nil && config.Import.Enabled, 
			config != nil && config.Import.ProcessOnStart)
	}

	return nil
}

// MatrixExists checks if a matrix with given name exists in database
func MatrixExists(name string) (bool, error) {
	if db == nil {
		return false, fmt.Errorf("veritabanı bağlantısı yok")
	}

	var count int
	err := db.db.QueryRow("SELECT COUNT(*) FROM matrix_records WHERE title = $1", name).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// SaveMatrixToDB saves a matrix to database and returns the ID
func SaveMatrixToDB(name string, matrix Matrix, source string) (int, error) {
	if db == nil {
		return 0, fmt.Errorf("veritabanı bağlantısı yok")
	}

	record, err := db.SaveMatrix(name, matrix, "")
	if err != nil {
		return 0, err
	}

	return record.ID, nil
}

// CalculateAndSaveAlgorithm calculates and saves algorithm result for a matrix
func CalculateAndSaveAlgorithm(matrixID int, matrix Matrix, algorithm string) error {
	if db == nil {
		return fmt.Errorf("veritabanı bağlantısı yok")
	}

	var boyarResult, paarResult, slpResult *AlgResult

	switch strings.ToLower(algorithm) {
	case "boyar":
		result, err := runBoyarSLP(matrix)
		if err != nil {
			return fmt.Errorf("boyar algoritması hatası: %v", err)
		}
		boyarResult = result

	case "paar":
		result, err := runPaarAlgorithm(matrix)
		if err != nil {
			return fmt.Errorf("paar algoritması hatası: %v", err)
		}
		paarResult = result

	case "slp":
		result, err := runSLPHeuristic(matrix)
		if err != nil {
			return fmt.Errorf("slp algoritması hatası: %v", err)
		}
		slpResult = result

	default:
		return fmt.Errorf("desteklenmeyen algoritma: %s", algorithm)
	}

	return db.UpdateMatrixResults(matrixID, boyarResult, paarResult, slpResult)
}

// calculateMatrixInverse calculates the inverse of a binary matrix using Gaussian elimination
func calculateMatrixInverse(matrix Matrix) (Matrix, error) {
	n := len(matrix)
	if n == 0 {
		return nil, fmt.Errorf("matris boş")
	}
	
	// Check if matrix is square
	if len(matrix[0]) != n {
		return nil, fmt.Errorf("matris kare değil: %dx%d", n, len(matrix[0]))
	}
	
	// Create augmented matrix [A|I]
	augmented := make([][]int, n)
	for i := 0; i < n; i++ {
		augmented[i] = make([]int, 2*n)
		// Copy original matrix
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(strings.TrimSpace(matrix[i][j]))
			if err != nil {
				return nil, fmt.Errorf("geçersiz matris değeri: %s", matrix[i][j])
			}
			augmented[i][j] = val
		}
		// Add identity matrix
		for j := n; j < 2*n; j++ {
			if j-n == i {
				augmented[i][j] = 1
			} else {
				augmented[i][j] = 0
			}
		}
	}
	
	// Gaussian elimination in GF(2)
	for i := 0; i < n; i++ {
		// Find pivot
		pivotRow := -1
		for k := i; k < n; k++ {
			if augmented[k][i] == 1 {
				pivotRow = k
				break
			}
		}
		
		if pivotRow == -1 {
			return nil, fmt.Errorf("matris tersi alınamaz (determinant = 0)")
		}
		
		// Swap rows if needed
		if pivotRow != i {
			augmented[i], augmented[pivotRow] = augmented[pivotRow], augmented[i]
		}
		
		// Eliminate column
		for k := 0; k < n; k++ {
			if k != i && augmented[k][i] == 1 {
				for j := 0; j < 2*n; j++ {
					augmented[k][j] ^= augmented[i][j] // XOR operation in GF(2)
				}
			}
		}
	}
	
	// Extract inverse matrix from right side of augmented matrix
	inverse := make(Matrix, n)
	for i := 0; i < n; i++ {
		inverse[i] = make([]string, n)
		for j := 0; j < n; j++ {
			inverse[i][j] = strconv.Itoa(augmented[i][j+n])
		}
	}
	
	return inverse, nil
}

// SaveMatrixInverse calculates and saves the inverse of a matrix
func (d *Database) SaveMatrixInverse(originalID int) (*MatrixRecord, error) {
	// Get original matrix
	original, err := d.GetMatrixByID(originalID)
	if err != nil {
		return nil, fmt.Errorf("orijinal matris alınamadı: %v", err)
	}
	
	if original == nil {
		return nil, fmt.Errorf("orijinal matris bulunamadı")
	}
	
	// Check if inverse already exists
	if original.InverseMatrixID != nil {
		// Return existing inverse
		return d.GetMatrixByID(*original.InverseMatrixID)
	}
	
	// Parse matrix from binary string
	matrix, err := parseMatrixFromBinary(original.MatrixBinary)
	if err != nil {
		return nil, fmt.Errorf("matris parse edilemedi: %v", err)
	}
	
	// Calculate inverse
	inverse, err := calculateMatrixInverse(matrix)
	if err != nil {
		return nil, fmt.Errorf("ters matris hesaplanamadı: %v", err)
	}
	
	// Create title for inverse matrix
	inverseTitle := original.Title + " (Ters)"
	
	// Check if inverse already exists by hash
	inverseHash := calculateMatrixHash(inverse)
	existing, err := d.GetMatrixByHash(inverseHash)
	if err == nil && existing != nil {
		// Update original matrix with inverse reference
		err = d.updateMatrixInverseReference(originalID, existing.ID, inverseHash)
		if err != nil {
			log.Printf("❌ Orijinal matrise ters matris referansı eklenemedi: %v", err)
		}
		return existing, nil
	}
	
	// Save inverse matrix
	inverseRecord, err := d.SaveMatrix(inverseTitle, inverse, original.Group)
	if err != nil {
		return nil, fmt.Errorf("ters matris kaydedilemedi: %v", err)
	}
	
	// Update original matrix with inverse reference
	err = d.updateMatrixInverseReference(originalID, inverseRecord.ID, inverseHash)
	if err != nil {
		log.Printf("❌ Orijinal matrise ters matris referansı eklenemedi: %v", err)
	}
	
	// Calculate algorithms for inverse matrix in background
	go func() {
		log.Printf("🔄 [INVERSE] %s için algoritma hesaplamaları başlıyor", inverseTitle)
		
		// Calculate Boyar SLP
		boyarResult, err := runBoyarSLP(inverse)
		if err != nil {
			log.Printf("❌ [INVERSE-BOYAR] %s için Boyar SLP hesaplanamadı: %v", inverseTitle, err)
		} else {
			log.Printf("✅ [INVERSE-BOYAR] %s için Boyar SLP tamamlandı - XOR: %d", inverseTitle, boyarResult.XorCount)
		}

		// Calculate Paar Algorithm
		paarResult, err := runPaarAlgorithm(inverse)
		if err != nil {
			log.Printf("❌ [INVERSE-PAAR] %s için Paar algoritması hesaplanamadı: %v", inverseTitle, err)
		} else {
			log.Printf("✅ [INVERSE-PAAR] %s için Paar algoritması tamamlandı - XOR: %d", inverseTitle, paarResult.XorCount)
		}

		// Calculate SLP Heuristic
		slpResult, err := runSLPHeuristic(inverse)
		if err != nil {
			log.Printf("❌ [INVERSE-SLP] %s için SLP Heuristic hesaplanamadı: %v", inverseTitle, err)
		} else {
			log.Printf("✅ [INVERSE-SLP] %s için SLP Heuristic tamamlandı - XOR: %d", inverseTitle, slpResult.XorCount)
		}

		// Update matrix with results
		err = d.UpdateMatrixResults(inverseRecord.ID, boyarResult, paarResult, slpResult)
		if err != nil {
			log.Printf("❌ [INVERSE-UPDATE] %s için sonuçlar kaydedilemedi: %v", inverseTitle, err)
		} else {
			log.Printf("✅ [INVERSE-UPDATE] %s için sonuçlar kaydedildi", inverseTitle)
		}
	}()
	
	return inverseRecord, nil
}

// updateMatrixInverseReference updates the original matrix with inverse reference
func (d *Database) updateMatrixInverseReference(originalID, inverseID int, inverseHash string) error {
	query := `
	UPDATE matrix_records 
	SET inverse_matrix_id = $1, inverse_matrix_hash = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	`
	
	_, err := d.db.Exec(query, inverseID, inverseHash, originalID)
	if err != nil {
		return fmt.Errorf("ters matris referansı güncellenemedi: %v", err)
	}
	
	return nil
}

// parseMatrixFromBinary parses a matrix from its binary string representation
func parseMatrixFromBinary(matrixBinary string) (Matrix, error) {
	lines := strings.Split(strings.TrimSpace(matrixBinary), "\n")
	matrix := make(Matrix, len(lines))
	
	for i, line := range lines {
		// Remove brackets and split by spaces
		line = strings.Trim(line, "[]")
		elements := strings.Fields(line)
		matrix[i] = elements
	}
	
	return matrix, nil
} 