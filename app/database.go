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
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	MatrixBinary    string    `json:"matrix_binary"`
	MatrixHex       string    `json:"matrix_hex"`
	HamXorCount     int       `json:"ham_xor_count"`
	BoyarXorCount   *int      `json:"boyar_xor_count,omitempty"`
	BoyarDepth      *int      `json:"boyar_depth,omitempty"`
	BoyarProgram    *string   `json:"boyar_program,omitempty"`
	PaarXorCount    *int      `json:"paar_xor_count,omitempty"`
	PaarProgram     *string   `json:"paar_program,omitempty"`
	SlpXorCount     *int      `json:"slp_xor_count,omitempty"`
	SlpProgram      *string   `json:"slp_program,omitempty"`
	MatrixHash      string    `json:"matrix_hash"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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
	return count
}

// SaveMatrix saves a matrix to the database
func (d *Database) SaveMatrix(title string, matrix Matrix) (*MatrixRecord, error) {
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
	INSERT INTO matrix_records (title, matrix_binary, matrix_hex, ham_xor_count, matrix_hash)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	var id int
	err = d.db.QueryRow(query, title, matrixBinary, matrixHex, hamXorCount, matrixHash).Scan(&id)
	if err != nil {
		return nil, err
	}

	return d.GetMatrixByID(id)
}

// UpdateMatrixResults updates the algorithm results for a matrix
func (d *Database) UpdateMatrixResults(id int, boyarResult, paarResult, slpResult *AlgResult) error {
	query := `
	UPDATE matrix_records 
	SET boyar_xor_count = $1, boyar_depth = $2, boyar_program = $3,
	    paar_xor_count = $4, paar_program = $5,
	    slp_xor_count = $6, slp_program = $7,
	    updated_at = CURRENT_TIMESTAMP
	WHERE id = $8
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

	_, err := d.db.Exec(query, boyarXor, boyarDepth, boyarProgram, paarXor, paarProgram, slpXor, slpProgram, id)
	return err
}

// GetMatrixByID retrieves a matrix by its ID
func (d *Database) GetMatrixByID(id int) (*MatrixRecord, error) {
	query := `
	SELECT id, title, matrix_binary, matrix_hex, ham_xor_count, 
	       boyar_xor_count, boyar_depth, boyar_program,
	       paar_xor_count, paar_program, slp_xor_count, slp_program,
	       matrix_hash, created_at, updated_at
	FROM matrix_records WHERE id = $1
	`
	
	row := d.db.QueryRow(query, id)
	return d.scanMatrixRecord(row)
}

// GetMatrixByHash retrieves a matrix by its hash
func (d *Database) GetMatrixByHash(hash string) (*MatrixRecord, error) {
	query := `
	SELECT id, title, matrix_binary, matrix_hex, ham_xor_count, 
	       boyar_xor_count, boyar_depth, boyar_program,
	       paar_xor_count, paar_program, slp_xor_count, slp_program,
	       matrix_hash, created_at, updated_at
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
		conditions = append(conditions, fmt.Sprintf("boyar_xor_count >= $%d", argIndex))
		args = append(args, *boyarXorMin)
		argIndex++
	}

	if boyarXorMax != nil {
		conditions = append(conditions, fmt.Sprintf("boyar_xor_count <= $%d", argIndex))
		args = append(args, *boyarXorMax)
		argIndex++
	}

	if paarXorMin != nil {
		conditions = append(conditions, fmt.Sprintf("paar_xor_count >= $%d", argIndex))
		args = append(args, *paarXorMin)
		argIndex++
	}

	if paarXorMax != nil {
		conditions = append(conditions, fmt.Sprintf("paar_xor_count <= $%d", argIndex))
		args = append(args, *paarXorMax)
		argIndex++
	}

	if slpXorMin != nil {
		conditions = append(conditions, fmt.Sprintf("slp_xor_count >= $%d", argIndex))
		args = append(args, *slpXorMin)
		argIndex++
	}

	if slpXorMax != nil {
		conditions = append(conditions, fmt.Sprintf("slp_xor_count <= $%d", argIndex))
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
	SELECT id, title, matrix_binary, matrix_hex, ham_xor_count, 
	       boyar_xor_count, boyar_depth, boyar_program,
	       paar_xor_count, paar_program, slp_xor_count, slp_program,
	       matrix_hash, created_at, updated_at
	FROM matrix_records %s
	ORDER BY created_at DESC
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
	var boyarXor, boyarDepth, paarXor, slpXor sql.NullInt64
	var boyarProgram, paarProgram, slpProgram sql.NullString

	var err error
	switch s := scanner.(type) {
	case *sql.Row:
		err = s.Scan(&record.ID, &record.Title, &record.MatrixBinary, &record.MatrixHex,
			&record.HamXorCount, &boyarXor, &boyarDepth, &boyarProgram,
			&paarXor, &paarProgram, &slpXor, &slpProgram,
			&record.MatrixHash, &record.CreatedAt, &record.UpdatedAt)
	case *sql.Rows:
		err = s.Scan(&record.ID, &record.Title, &record.MatrixBinary, &record.MatrixHex,
			&record.HamXorCount, &boyarXor, &boyarDepth, &boyarProgram,
			&paarXor, &paarProgram, &slpXor, &slpProgram,
			&record.MatrixHash, &record.CreatedAt, &record.UpdatedAt)
	default:
		return nil, fmt.Errorf("unsupported scanner type")
	}

	if err != nil {
		return nil, err
	}

	// Handle nullable fields
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
	log.Printf("Matrices-data klasöründen matrisler import ediliyor: %s", dataPath)
	
	// Get list of .txt files in the directory
	files, err := filepath.Glob(filepath.Join(dataPath, "*.txt"))
	if err != nil {
		return fmt.Errorf("dosya listesi alınamadı: %v", err)
	}

	if len(files) == 0 {
		log.Printf("Matrices-data klasöründe .txt dosyası bulunamadı")
		return nil
	}

	log.Printf("%d dosya bulundu, import işlemi başlıyor...", len(files))

	totalImported := 0
	for _, filePath := range files {
		fileName := filepath.Base(filePath)
		log.Printf("Dosya işleniyor: %s", fileName)
		
		count, err := d.importMatricesFromFile(filePath)
		if err != nil {
			log.Printf("HATA: %s dosyası işlenirken hata oluştu: %v", fileName, err)
			continue
		}
		
		totalImported += count
		log.Printf("%s dosyasından %d matris import edildi", fileName, count)
	}

	log.Printf("Import işlemi tamamlandı. Toplam %d matris import edildi.", totalImported)
	return nil
}

// importMatricesFromFile imports matrices from a single file
func (d *Database) importMatricesFromFile(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentMatrix [][]string
	var currentTitle string
	importedCount := 0
	lineNumber := 0

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
				err := d.saveMatrixFromImport(currentTitle, currentMatrix)
				if err != nil {
					log.Printf("HATA: Matris kaydedilemedi (satır %d): %v", lineNumber, err)
				} else {
					importedCount++
				}
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

		// Skip other lines (like "HamXOR Sayisi:" etc.)
	}

	// Process the last matrix if exists
	if len(currentMatrix) > 0 && currentTitle != "" {
		err := d.saveMatrixFromImport(currentTitle, currentMatrix)
		if err != nil {
			log.Printf("HATA: Son matris kaydedilemedi: %v", err)
		} else {
			importedCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return importedCount, err
	}

	return importedCount, nil
}

// saveMatrixFromImport saves a matrix during import process with algorithm calculation
func (d *Database) saveMatrixFromImport(title string, matrix [][]string) error {
	// Check if matrix already exists by hash
	matrixHash := calculateMatrixHash(matrix)
	existing, err := d.GetMatrixByHash(matrixHash)
	if err == nil && existing != nil {
		// Matrix already exists, skip
		return nil
	}

	// Save the matrix first
	record, err := d.SaveMatrix(title, matrix)
	if err != nil {
		return err
	}

	// Calculate algorithms immediately during import
	go func() {
		log.Printf("Matris %d için algoritmalar hesaplanıyor...", record.ID)
		
		var boyarResult, paarResult, slpResult *AlgResult

		// Boyar algorithm
		boyar := NewBoyarSLP(10)
		if result, err := boyar.Solve(matrix); err == nil {
			boyarResult = &result
		} else {
			log.Printf("Boyar algoritması hatası (ID %d): %v", record.ID, err)
		}

		// Paar algorithm
		paar := NewPaarAlgorithm()
		if result, err := paar.Solve(matrix); err == nil {
			paarResult = &result
		} else {
			log.Printf("Paar algoritması hatası (ID %d): %v", record.ID, err)
		}

		// SLP algorithm
		slp := NewSLPHeuristic()
		if result, err := slp.Solve(matrix); err == nil {
			slpResult = &result
		} else {
			log.Printf("SLP algoritması hatası (ID %d): %v", record.ID, err)
		}

		// Update database with results
		err = d.UpdateMatrixResults(record.ID, boyarResult, paarResult, slpResult)
		if err != nil {
			log.Printf("Algoritma sonuçları güncellenemedi (ID %d): %v", record.ID, err)
		} else {
			log.Printf("Matris %d algoritmaları tamamlandı", record.ID)
		}
	}()

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

// Global database instance
var db *Database

// InitDatabase initializes the database connection
func InitDatabase() error {
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

	// Check if we need to import matrices using hash comparison
	go func() {
		time.Sleep(5 * time.Second) // Wait a bit for the application to fully start
		
		dataPath := os.Getenv("MATRICES_DATA_PATH")
		if dataPath == "" {
			dataPath = "./matrices-data"
		}

		// Get hashes from database
		dbHashes, err := db.GetAllMatrixHashes()
		if err != nil {
			log.Printf("HATA: Veritabanı hash'leri alınamadı: %v", err)
			return
		}

		// Get hashes from files
		fileHashes, err := db.GetFileMatrixHashes(dataPath)
		if err != nil {
			log.Printf("HATA: Dosya hash'leri alınamadı: %v", err)
			return
		}

		log.Printf("Veritabanında %d matris hash'i bulundu", len(dbHashes))
		log.Printf("Dosyalarda %d matris hash'i bulundu", len(fileHashes))

		// Check for missing matrices
		missingCount := 0
		for hash := range fileHashes {
			if !dbHashes[hash] {
				missingCount++
			}
		}

		if missingCount > 0 {
			log.Printf("Veritabanında %d eksik matris var, import işlemi başlatılıyor...", missingCount)
			err := db.ImportMatricesFromFiles(dataPath)
			if err != nil {
				log.Printf("HATA: Matris import işlemi başarısız: %v", err)
			}
		} else {
			log.Printf("Tüm matrisler zaten veritabanında mevcut")
		}
	}()

	return nil
} 