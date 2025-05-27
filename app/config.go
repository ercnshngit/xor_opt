package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig `json:"database"`
	Import   ImportConfig   `json:"import"`
	Server   ServerConfig   `json:"server"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}

// ImportConfig holds auto import configuration
type ImportConfig struct {
	Enabled         bool     `json:"enabled"`
	DataDirectory   string   `json:"data_directory"`
	FileExtensions  []string `json:"file_extensions"`
	MaxFileSize     int64    `json:"max_file_size_mb"`
	ProcessOnStart  bool     `json:"process_on_start"`
	WatchDirectory  bool     `json:"watch_directory"`
	SkipExisting    bool     `json:"skip_existing"`
	BatchSize       int      `json:"batch_size"`
	AutoCalculate   bool     `json:"auto_calculate"`
	Algorithms      []string `json:"algorithms"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string `json:"port"`
	Host         string `json:"host"`
	EnableCORS   bool   `json:"enable_cors"`
	LogLevel     string `json:"log_level"`
	StaticDir    string `json:"static_dir"`
}

// Default configuration
var defaultConfig = Config{
	Database: DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "xor_optimization",
		SSLMode:  "disable",
	},
	Import: ImportConfig{
		Enabled:         true,
		DataDirectory:   "./matrices-data",
		FileExtensions:  []string{".txt", ".csv", ".json"},
		MaxFileSize:     500, // MB
		ProcessOnStart:  true,
		WatchDirectory:  false,
		SkipExisting:    true,
		BatchSize:       10,
		AutoCalculate:   true,
		Algorithms:      []string{"boyar", "paar", "slp"},
	},
	Server: ServerConfig{
		Port:         ":3000",
		Host:         "localhost",
		EnableCORS:   true,
		LogLevel:     "info",
		StaticDir:    "./web",
	},
}

// LoadConfig loads configuration from file or creates default
func LoadConfig(configPath string) (*Config, error) {
	config := defaultConfig

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config dosyası bulunamadı (%s), varsayılan config oluşturuluyor...", configPath)
		
		// Create config file with default values
		if err := SaveConfig(&config, configPath); err != nil {
			return nil, fmt.Errorf("varsayılan config dosyası oluşturulamadı: %v", err)
		}
		
		log.Printf("Varsayılan config dosyası oluşturuldu: %s", configPath)
		return &config, nil
	}

	// Read existing config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config dosyası okunamadı: %v", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("config dosyası parse edilemedi: %v", err)
	}

	log.Printf("Config dosyası yüklendi: %s", configPath)
	return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config, configPath string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("config dizini oluşturulamadı: %v", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("config JSON'a dönüştürülemedi: %v", err)
	}

	// Write to file
	if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("config dosyası yazılamadı: %v", err)
	}

	return nil
}

// AutoImportData automatically imports data files based on configuration
func AutoImportData(config *Config) error {
	if !config.Import.Enabled {
		log.Println("Otomatik import devre dışı")
		return nil
	}

	log.Printf("Otomatik import başlatılıyor: %s", config.Import.DataDirectory)

	// Check if data directory exists
	if _, err := os.Stat(config.Import.DataDirectory); os.IsNotExist(err) {
		log.Printf("Data dizini bulunamadı: %s", config.Import.DataDirectory)
		return nil
	}

	// Walk through data directory
	err := filepath.Walk(config.Import.DataDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Dosya okuma hatası: %v", err)
			return nil // Continue with other files
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check file extension
		ext := strings.ToLower(filepath.Ext(path))
		validExt := false
		for _, allowedExt := range config.Import.FileExtensions {
			if ext == allowedExt {
				validExt = true
				break
			}
		}

		if !validExt {
			log.Printf("Desteklenmeyen dosya uzantısı atlanıyor: %s", path)
			return nil
		}

		// Check file size
		maxSize := config.Import.MaxFileSize * 1024 * 1024 // Convert MB to bytes
		if info.Size() > maxSize {
			log.Printf("Dosya çok büyük, atlanıyor: %s (%.2f MB)", path, float64(info.Size())/(1024*1024))
			return nil
		}

		// Import file
		log.Printf("Dosya import ediliyor: %s", path)
		if err := ImportMatrixFile(path, config); err != nil {
			log.Printf("Dosya import hatası (%s): %v", path, err)
			return nil // Continue with other files
		}

		log.Printf("Dosya başarıyla import edildi: %s", path)
		return nil
	})

	if err != nil {
		return fmt.Errorf("dizin tarama hatası: %v", err)
	}

	log.Println("Otomatik import tamamlandı")
	return nil
}

// ImportMatrixFile imports a single matrix file
func ImportMatrixFile(filePath string, config *Config) error {
	// Read file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("dosya okunamadı: %v", err)
	}

	// Parse matrix based on file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	var matrices []Matrix
	
	switch ext {
	case ".txt":
		matrices, err = ParseTextMatrix(string(content))
	case ".csv":
		matrices, err = ParseCSVMatrix(string(content))
	case ".json":
		matrices, err = ParseJSONMatrix(string(content))
	default:
		return fmt.Errorf("desteklenmeyen dosya formatı: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("matrix parse hatası: %v", err)
	}

	// Get filename without extension for naming
	filename := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	// Process each matrix
	for i, matrix := range matrices {
		matrixName := fmt.Sprintf("%s_matrix_%d", filename, i+1)
		
		// Check if matrix already exists (if skip_existing is enabled)
		if config.Import.SkipExisting {
			exists, err := MatrixExists(matrixName)
			if err != nil {
				log.Printf("Matrix varlık kontrolü hatası: %v", err)
				continue
			}
			if exists {
				log.Printf("Matrix zaten mevcut, atlanıyor: %s", matrixName)
				continue
			}
		}

		// Save matrix to database
		matrixID, err := SaveMatrixToDB(matrixName, matrix, "auto_import")
		if err != nil {
			log.Printf("Matrix kaydetme hatası (%s): %v", matrixName, err)
			continue
		}

		log.Printf("Matrix kaydedildi: %s (ID: %d)", matrixName, matrixID)

		// Auto calculate algorithms if enabled
		if config.Import.AutoCalculate {
			for _, algorithm := range config.Import.Algorithms {
				if err := CalculateAndSaveAlgorithm(matrixID, matrix, algorithm); err != nil {
					log.Printf("Algoritma hesaplama hatası (%s, %s): %v", matrixName, algorithm, err)
				} else {
					log.Printf("Algoritma hesaplandı: %s -> %s", matrixName, algorithm)
				}
			}
		}
	}

	return nil
}

// ParseTextMatrix parses text format matrix
func ParseTextMatrix(content string) ([]Matrix, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("boş dosya")
	}

	var matrix Matrix
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Split by spaces or tabs
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		
		matrix = append(matrix, fields)
	}

	if len(matrix) == 0 {
		return nil, fmt.Errorf("geçerli matrix verisi bulunamadı")
	}

	return []Matrix{matrix}, nil
}

// ParseCSVMatrix parses CSV format matrix
func ParseCSVMatrix(content string) ([]Matrix, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("boş dosya")
	}

	var matrix Matrix
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Split by comma
		fields := strings.Split(line, ",")
		for i, field := range fields {
			fields[i] = strings.TrimSpace(field)
		}
		
		matrix = append(matrix, fields)
	}

	if len(matrix) == 0 {
		return nil, fmt.Errorf("geçerli matrix verisi bulunamadı")
	}

	return []Matrix{matrix}, nil
}

// ParseJSONMatrix parses JSON format matrix
func ParseJSONMatrix(content string) ([]Matrix, error) {
	var data struct {
		Matrices []Matrix `json:"matrices"`
		Matrix   Matrix   `json:"matrix"`
	}

	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("JSON parse hatası: %v", err)
	}

	if len(data.Matrices) > 0 {
		return data.Matrices, nil
	}

	if len(data.Matrix) > 0 {
		return []Matrix{data.Matrix}, nil
	}

	return nil, fmt.Errorf("JSON'da matrix verisi bulunamadı")
} 