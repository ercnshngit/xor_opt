package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// SaveMatrixRequest represents the request to save a matrix
type SaveMatrixRequest struct {
	Title  string `json:"title"`
	Group  string `json:"group,omitempty"`
	Matrix Matrix `json:"matrix"`
}

// GetMatricesResponse represents the response for getting matrices
type GetMatricesResponse struct {
	Matrices   []*MatrixRecord `json:"matrices"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	TotalPages int             `json:"total_pages"`
}

// RecalculateRequest represents the request to recalculate algorithms
type RecalculateRequest struct {
	MatrixID   int      `json:"matrix_id"`
	Algorithms []string `json:"algorithms"` // ["boyar", "paar", "slp"]
}

// BulkRecalculateRequest represents the request to recalculate algorithms for multiple matrices
type BulkRecalculateRequest struct {
	Algorithms []string `json:"algorithms"` // ["boyar", "paar", "slp"]
	Limit      int      `json:"limit"`      // Maximum number of matrices to process
}

// BulkRecalculateResponse represents the response for bulk recalculation
type BulkRecalculateResponse struct {
	ProcessedCount int `json:"processed_count"`
	TotalCount     int `json:"total_count"`
	Message        string `json:"message"`
}

// BulkInverseRequest represents the request for bulk inverse calculation
type BulkInverseRequest struct {
	MaxSmallestXor     int  `json:"max_smallest_xor"`
	SkipExisting       bool `json:"skip_existing"`
	CalculateAlgorithms bool `json:"calculate_algorithms"`
	BatchSize          int  `json:"batch_size"`
}

// BulkInverseResponse represents the response for bulk inverse calculation
type BulkInverseResponse struct {
	ProcessedCount int                    `json:"processed_count"`
	SuccessCount   int                    `json:"success_count"`
	ErrorCount     int                    `json:"error_count"`
	Results        []BulkInverseResult    `json:"results"`
	Message        string                 `json:"message"`
}

// BulkInverseResult represents a single result in bulk inverse calculation
type BulkInverseResult struct {
	OriginalID    int    `json:"original_id"`
	OriginalTitle string `json:"original_title"`
	InverseID     *int   `json:"inverse_id,omitempty"`
	InverseTitle  string `json:"inverse_title,omitempty"`
	Status        string `json:"status"` // "success", "error", "skipped"
	Message       string `json:"message"`
}

// GetInversePairsResponse represents the response for inverse pairs
type GetInversePairsResponse struct {
	Pairs      []*InversePair `json:"pairs"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// MissingAlgorithmsResponse represents the response for missing algorithms count
type MissingAlgorithmsResponse struct {
	MissingCount int `json:"missing_count"`
	TotalCount   int `json:"total_count"`
	Message      string `json:"message"`
}

// saveMatrixHandler saves a matrix to the database
func saveMatrixHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req SaveMatrixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Başlık gerekli", http.StatusBadRequest)
		return
	}

	if len(req.Matrix) == 0 {
		http.Error(w, "Matris gerekli", http.StatusBadRequest)
		return
	}

	record, err := db.SaveMatrix(req.Title, req.Matrix, req.Group)
	if err != nil {
		http.Error(w, "Matris kaydedilemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(record)
}

// getMatricesHandler retrieves matrices with pagination and filtering
func getMatricesHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	w.Header().Set("Content-Type", "application/json")
	
	// Add cache headers for better performance
	w.Header().Set("Cache-Control", "public, max-age=30")
	w.Header().Set("ETag", fmt.Sprintf("\"%d\"", time.Now().Unix()/30)) // 30 second cache

	// Parse query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	titleFilter := r.URL.Query().Get("title")
	groupFilter := r.URL.Query().Get("group")

	// Parse range filters
	var hamXorMin, hamXorMax, boyarXorMin, boyarXorMax, paarXorMin, paarXorMax, slpXorMin, slpXorMax *int

	if val := r.URL.Query().Get("ham_xor_min"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			hamXorMin = &parsed
		}
	}

	if val := r.URL.Query().Get("ham_xor_max"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			hamXorMax = &parsed
		}
	}

	if val := r.URL.Query().Get("boyar_xor_min"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			boyarXorMin = &parsed
		}
	}

	if val := r.URL.Query().Get("boyar_xor_max"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			boyarXorMax = &parsed
		}
	}

	if val := r.URL.Query().Get("paar_xor_min"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			paarXorMin = &parsed
		}
	}

	if val := r.URL.Query().Get("paar_xor_max"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			paarXorMax = &parsed
		}
	}

	if val := r.URL.Query().Get("slp_xor_min"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			slpXorMin = &parsed
		}
	}

	if val := r.URL.Query().Get("slp_xor_max"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			slpXorMax = &parsed
		}
	}

	log.Printf("📊 [API] GetMatrices request: page=%d, limit=%d, title_filter='%s', group_filter='%s'", page, limit, titleFilter, groupFilter)

	matrices, total, err := db.GetMatrices(page, limit, titleFilter, groupFilter, hamXorMin, hamXorMax, boyarXorMin, boyarXorMax, paarXorMin, paarXorMax, slpXorMin, slpXorMax)
	if err != nil {
		log.Printf("❌ [API] GetMatrices error: %v", err)
		http.Error(w, "Matrisler alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (total + limit - 1) / limit

	response := GetMatricesResponse{
		Matrices:   matrices,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	duration := time.Since(startTime)
	log.Printf("✅ [API] GetMatrices completed in %v: returned %d matrices (total: %d)", duration, len(matrices), total)

	json.NewEncoder(w).Encode(response)
}

// getMatrixHandler retrieves a single matrix by ID
func getMatrixHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Geçersiz ID", http.StatusBadRequest)
		return
	}

	record, err := db.GetMatrixByID(id)
	if err != nil {
		http.Error(w, "Matris alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if record == nil {
		http.Error(w, "Matris bulunamadı", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(record)
}

// recalculateHandler recalculates algorithms for a specific matrix
func recalculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RecalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	if req.MatrixID <= 0 {
		http.Error(w, "Geçersiz matris ID", http.StatusBadRequest)
		return
	}

	// Default algorithms if not specified
	if len(req.Algorithms) == 0 {
		req.Algorithms = []string{"boyar", "paar", "slp", "sbp"}
	}

	// Get matrix from database
	record, err := db.GetMatrixByID(req.MatrixID)
	if err != nil {
		http.Error(w, "Matris alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if record == nil {
		http.Error(w, "Matris bulunamadı", http.StatusNotFound)
		return
	}

	// Parse matrix from binary string
	matrix, err := parseMatrixFromBinary(record.MatrixBinary)
	if err != nil {
		http.Error(w, "Matris parse edilemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Run algorithms in background
	go func() {
		startTime := time.Now()
		log.Printf("🚀 [RECALCULATE] Matris %d için algoritma hesaplama başlatıldı", req.MatrixID)

		var boyarResult, paarResult, slpResult, sbpResult *AlgResult

		// Run requested algorithms
		for _, algorithm := range req.Algorithms {
			algorithmStartTime := time.Now()
			switch strings.ToLower(algorithm) {
			case "boyar":
				log.Printf("⏳ [BOYAR] Matris %d için Boyar algoritması başlatıldı", req.MatrixID)
				boyar := NewBoyarSLP(10)
				if result, err := boyar.Solve(matrix); err == nil {
					boyarResult = &result
					duration := time.Since(algorithmStartTime)
					log.Printf("✅ [BOYAR] Matris %d - Boyar algoritması tamamlandı: %.2f saniye (XOR count: %d)", 
						req.MatrixID, duration.Seconds(), result.XorCount)
				} else {
					duration := time.Since(algorithmStartTime)
					log.Printf("❌ [BOYAR] Matris %d - Boyar algoritması hatası (%.2f saniye): %v", req.MatrixID, duration.Seconds(), err)
				}
			case "paar":
				log.Printf("⏳ [PAAR] Matris %d için Paar algoritması başlatıldı", req.MatrixID)
				paar := NewPaarAlgorithm()
				if result, err := paar.Solve(matrix); err == nil {
					paarResult = &result
					duration := time.Since(algorithmStartTime)
					log.Printf("✅ [PAAR] Matris %d - Paar algoritması tamamlandı: %.2f saniye (XOR count: %d)", 
						req.MatrixID, duration.Seconds(), result.XorCount)
				} else {
					duration := time.Since(algorithmStartTime)
					log.Printf("❌ [PAAR] Matris %d - Paar algoritması hatası (%.2f saniye): %v", req.MatrixID, duration.Seconds(), err)
				}
			case "slp":
				log.Printf("⏳ [SLP] Matris %d için SLP algoritması başlatıldı", req.MatrixID)
				slp := NewSLPHeuristic()
				if result, err := slp.Solve(matrix); err == nil {
					slpResult = &result
					duration := time.Since(algorithmStartTime)
					log.Printf("✅ [SLP] Matris %d - SLP algoritması tamamlandı: %.2f saniye (XOR count: %d)", 
						req.MatrixID, duration.Seconds(), result.XorCount)
				} else {
					duration := time.Since(algorithmStartTime)
					log.Printf("❌ [SLP] Matris %d - SLP algoritması hatası (%.2f saniye): %v", req.MatrixID, duration.Seconds(), err)
				}
			case "sbp":
				log.Printf("⏳ [SBP] Matris %d için SBP algoritması başlatıldı", req.MatrixID)
				sbp := NewSBPAlgorithm(10)
				if result, err := sbp.Solve(matrix); err == nil {
					sbpResult = &result
					duration := time.Since(algorithmStartTime)
					log.Printf("✅ [SBP] Matris %d - SBP algoritması tamamlandı: %.2f saniye (XOR count: %d)", 
						req.MatrixID, duration.Seconds(), result.XorCount)
				} else {
					duration := time.Since(algorithmStartTime)
					log.Printf("❌ [SBP] Matris %d - SBP algoritması hatası (%.2f saniye): %v", req.MatrixID, duration.Seconds(), err)
				}
			}
		}

		// Update database with results
		dbUpdateStartTime := time.Now()
		if len(req.Algorithms) == 1 && strings.ToLower(req.Algorithms[0]) == "sbp" && sbpResult != nil {
			// Only SBP requested, use specific SBP update function
			err = db.UpdateSBPResults(req.MatrixID, sbpResult)
		} else {
			// Multiple algorithms or non-SBP, use full update function
			err = db.UpdateMatrixResultsWithSBP(req.MatrixID, boyarResult, paarResult, slpResult, sbpResult)
		}
		dbUpdateDuration := time.Since(dbUpdateStartTime)
		totalDuration := time.Since(startTime)
		
		if err != nil {
			log.Printf("❌ [RECALCULATE] Matris %d - Algoritma sonuçları güncellenemedi (DB güncelleme: %.2f saniye, Toplam: %.2f saniye): %v", 
				req.MatrixID, dbUpdateDuration.Seconds(), totalDuration.Seconds(), err)
		} else {
			log.Printf("🎉 [RECALCULATE] Matris %d algoritmaları tamamlandı! (DB güncelleme: %.2f saniye, Toplam: %.2f saniye)", 
				req.MatrixID, dbUpdateDuration.Seconds(), totalDuration.Seconds())
		}
	}()

	// Return updated record
	updatedRecord, err := db.GetMatrixByID(req.MatrixID)
	if err != nil {
		http.Error(w, "Güncellenmiş matris alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedRecord)
}

// processAndSaveMatrixHandler processes a matrix with all algorithms and saves to database
func processAndSaveMatrixHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req SaveMatrixRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Başlık gerekli", http.StatusBadRequest)
		return
	}

	if len(req.Matrix) == 0 {
		http.Error(w, "Matris gerekli", http.StatusBadRequest)
		return
	}

	// Save matrix first
	record, err := db.SaveMatrix(req.Title, req.Matrix, req.Group)
	if err != nil {
		http.Error(w, "Matris kaydedilemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Run all algorithms
	var boyarResult, paarResult, slpResult, sbpResult *AlgResult

	// Boyar algorithm
	boyar := NewBoyarSLP(10)
	if result, err := boyar.Solve(req.Matrix); err == nil {
		boyarResult = &result
	}

	// Paar algorithm
	paar := NewPaarAlgorithm()
	if result, err := paar.Solve(req.Matrix); err == nil {
		paarResult = &result
	}

	// SLP algorithm
	slp := NewSLPHeuristic()
	if result, err := slp.Solve(req.Matrix); err == nil {
		slpResult = &result
	}

	// SBP algorithm
	sbp := NewSBPAlgorithm(10)
	if result, err := sbp.Solve(req.Matrix); err == nil {
		sbpResult = &result
	}

	// Update database with results
	err = db.UpdateMatrixResultsWithSBP(record.ID, boyarResult, paarResult, slpResult, sbpResult)
	if err != nil {
		http.Error(w, "Sonuçlar güncellenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated record
	updatedRecord, err := db.GetMatrixByID(record.ID)
	if err != nil {
		http.Error(w, "Güncellenmiş matris alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedRecord)
}

// bulkRecalculateHandler recalculates algorithms for matrices without algorithm results
func bulkRecalculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req BulkRecalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	// Default algorithms if not specified
	if len(req.Algorithms) == 0 {
		req.Algorithms = []string{"boyar", "paar", "slp", "sbp"}
	}

	// Default limit if not specified
	if req.Limit <= 0 {
		req.Limit = 100
	}

	// Get matrices without algorithm results
	matrices, err := db.GetMatricesWithoutAlgorithms(req.Limit)
	if err != nil {
		http.Error(w, "Matrisler alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(matrices) == 0 {
		response := BulkRecalculateResponse{
			ProcessedCount: 0,
			TotalCount:     0,
			Message:        "Algoritma hesaplanacak matris bulunamadı",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Process matrices in background
	go func() {
		bulkStartTime := time.Now()
		log.Printf("🚀 [BULK-RECALCULATE] Toplu hesaplama başlatıldı - %d matris işlenecek", len(matrices))
		
		for i, matrix := range matrices {
			matrixStartTime := time.Now()
			log.Printf("⏳ [BULK-RECALCULATE] Matris %d/%d (ID: %d) işleniyor: %s", i+1, len(matrices), matrix.ID, matrix.Title)
			
			// Parse matrix from binary string
			matrixData, err := parseMatrixFromBinary(matrix.MatrixBinary)
			if err != nil {
				log.Printf("❌ [BULK-RECALCULATE] Matris parse hatası (ID %d): %v", matrix.ID, err)
				continue
			}

			// Recalculate Ham XOR
			hamXorStartTime := time.Now()
			newHamXor := calculateHammingXOR(matrixData)
			hamXorDuration := time.Since(hamXorStartTime)
			
			// Update Ham XOR in database
			_, err = db.db.Exec("UPDATE matrix_records SET ham_xor_count = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", newHamXor, matrix.ID)
			if err != nil {
				log.Printf("❌ [BULK-RECALCULATE] Ham XOR güncellenemedi (ID %d): %v", matrix.ID, err)
			} else {
				log.Printf("✅ [HAM-XOR] Matris %d - Ham XOR hesaplandı: %.3f saniye (Değer: %d)", 
					matrix.ID, hamXorDuration.Seconds(), newHamXor)
			}

			var boyarResult, paarResult, slpResult, sbpResult *AlgResult

			// Run requested algorithms
			for _, algorithm := range req.Algorithms {
				algorithmStartTime := time.Now()
				switch strings.ToLower(algorithm) {
				case "boyar":
					log.Printf("⏳ [BOYAR] Matris %d için Boyar algoritması başlatıldı", matrix.ID)
					boyar := NewBoyarSLP(10)
					if result, err := boyar.Solve(matrixData); err == nil {
						boyarResult = &result
						duration := time.Since(algorithmStartTime)
						log.Printf("✅ [BOYAR] Matris %d - Boyar algoritması tamamlandı: %.2f saniye (XOR count: %d)", 
							matrix.ID, duration.Seconds(), result.XorCount)
					} else {
						duration := time.Since(algorithmStartTime)
						log.Printf("❌ [BOYAR] Matris %d - Boyar algoritması hatası (%.2f saniye): %v", matrix.ID, duration.Seconds(), err)
					}
				case "paar":
					log.Printf("⏳ [PAAR] Matris %d için Paar algoritması başlatıldı", matrix.ID)
					paar := NewPaarAlgorithm()
					if result, err := paar.Solve(matrixData); err == nil {
						paarResult = &result
						duration := time.Since(algorithmStartTime)
						log.Printf("✅ [PAAR] Matris %d - Paar algoritması tamamlandı: %.2f saniye (XOR count: %d)", 
							matrix.ID, duration.Seconds(), result.XorCount)
					} else {
						duration := time.Since(algorithmStartTime)
						log.Printf("❌ [PAAR] Matris %d - Paar algoritması hatası (%.2f saniye): %v", matrix.ID, duration.Seconds(), err)
					}
				case "slp":
					log.Printf("⏳ [SLP] Matris %d için SLP algoritması başlatıldı", matrix.ID)
					slp := NewSLPHeuristic()
					if result, err := slp.Solve(matrixData); err == nil {
						slpResult = &result
						duration := time.Since(algorithmStartTime)
						log.Printf("✅ [SLP] Matris %d - SLP algoritması tamamlandı: %.2f saniye (XOR count: %d)", 
							matrix.ID, duration.Seconds(), result.XorCount)
					} else {
						duration := time.Since(algorithmStartTime)
						log.Printf("❌ [SLP] Matris %d - SLP algoritması hatası (%.2f saniye): %v", matrix.ID, duration.Seconds(), err)
					}
				case "sbp":
					log.Printf("⏳ [SBP] Matris %d için SBP algoritması başlatıldı", matrix.ID)
					sbp := NewSBPAlgorithm(10)
					if result, err := sbp.Solve(matrixData); err == nil {
						sbpResult = &result
						duration := time.Since(algorithmStartTime)
						log.Printf("✅ [SBP] Matris %d - SBP algoritması tamamlandı: %.2f saniye (XOR count: %d)", 
							matrix.ID, duration.Seconds(), result.XorCount)
					} else {
						duration := time.Since(algorithmStartTime)
						log.Printf("❌ [SBP] Matris %d - SBP algoritması hatası (%.2f saniye): %v", matrix.ID, duration.Seconds(), err)
					}
				}
			}

			// Update database with results
			dbUpdateStartTime := time.Now()
			err = db.UpdateMatrixResultsWithSBP(matrix.ID, boyarResult, paarResult, slpResult, sbpResult)
			dbUpdateDuration := time.Since(dbUpdateStartTime)
			matrixTotalDuration := time.Since(matrixStartTime)
			
			if err != nil {
				log.Printf("❌ [BULK-RECALCULATE] Matris %d - Algoritma sonuçları güncellenemedi (DB güncelleme: %.2f saniye, Matris toplam: %.2f saniye): %v", 
					matrix.ID, dbUpdateDuration.Seconds(), matrixTotalDuration.Seconds(), err)
			} else {
				log.Printf("🎉 [BULK-RECALCULATE] Matris %d algoritmaları tamamlandı! (DB güncelleme: %.2f saniye, Matris toplam: %.2f saniye)", 
					matrix.ID, dbUpdateDuration.Seconds(), matrixTotalDuration.Seconds())
			}
		}
		bulkTotalDuration := time.Since(bulkStartTime)
		log.Printf("🏁 [BULK-RECALCULATE] Toplu algoritma hesaplama tamamlandı: %d matris işlendi (Toplam süre: %.2f saniye, Ortalama: %.2f saniye/matris)", 
			len(matrices), bulkTotalDuration.Seconds(), bulkTotalDuration.Seconds()/float64(len(matrices)))
	}()

	response := BulkRecalculateResponse{
		ProcessedCount: 0,
		TotalCount:     len(matrices),
		Message:        fmt.Sprintf("%d matris için algoritma hesaplama başlatıldı", len(matrices)),
	}
	json.NewEncoder(w).Encode(response)
}

// calculateInverseHandler calculates and saves the inverse of a matrix
func calculateInverseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Geçersiz ID formatı", http.StatusBadRequest)
		return
	}

	// Calculate and save inverse matrix
	inverseRecord, err := db.SaveMatrixInverse(id)
	if err != nil {
		http.Error(w, "Ters matris hesaplanamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(inverseRecord)
}

// bulkInverseHandler handles bulk inverse matrix calculation
func bulkInverseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req BulkInverseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	if req.MaxSmallestXor <= 0 {
		http.Error(w, "Geçersiz smallest XOR değeri", http.StatusBadRequest)
		return
	}

	// Default batch size
	if req.BatchSize <= 0 {
		req.BatchSize = 10
	}

	log.Printf("🔄 [BULK-INVERSE] Toplu ters alma başlatıldı - Max XOR: %d, Batch: %d, Skip existing: %v", 
		req.MaxSmallestXor, req.BatchSize, req.SkipExisting)

	// Get matrices with smallest XOR less than the specified value
	matrices, err := db.GetMatricesForBulkInverse(req.MaxSmallestXor, req.SkipExisting)
	if err != nil {
		http.Error(w, "Matrisler alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(matrices) == 0 {
		response := BulkInverseResponse{
			ProcessedCount: 0,
			SuccessCount:   0,
			ErrorCount:     0,
			Results:        []BulkInverseResult{},
			Message:        "Belirtilen kriterlere uygun matris bulunamadı",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("📊 [BULK-INVERSE] %d matris bulundu, işlem başlıyor", len(matrices))

	// Process matrices in background
	go func() {
		results := make([]BulkInverseResult, 0, len(matrices))
		successCount := 0
		errorCount := 0

		for i, matrix := range matrices {
			log.Printf("🔄 [BULK-INVERSE] Matris %d/%d işleniyor (ID: %d): %s", 
				i+1, len(matrices), matrix.ID, matrix.Title)

			result := BulkInverseResult{
				OriginalID:    matrix.ID,
				OriginalTitle: matrix.Title,
			}

			// Check if inverse already exists (double check)
			if req.SkipExisting && matrix.InverseMatrixID != nil {
				result.Status = "skipped"
				result.Message = "Ters matris zaten mevcut"
				results = append(results, result)
				continue
			}

			// Calculate inverse matrix
			inverseRecord, err := db.SaveMatrixInverse(matrix.ID)
			if err != nil {
				log.Printf("❌ [BULK-INVERSE] Matris %d için ters hesaplanamadı: %v", matrix.ID, err)
				result.Status = "error"
				result.Message = err.Error()
				errorCount++
			} else {
				log.Printf("✅ [BULK-INVERSE] Matris %d için ters hesaplandı (ID: %d)", matrix.ID, inverseRecord.ID)
				result.Status = "success"
				result.InverseID = &inverseRecord.ID
				result.InverseTitle = inverseRecord.Title
				result.Message = "Ters matris başarıyla hesaplandı"
				successCount++

				// Calculate algorithms if requested
				if req.CalculateAlgorithms {
					log.Printf("🧮 [BULK-INVERSE] Matris %d için algoritma hesaplamaları başlatılıyor", inverseRecord.ID)
					
					// Parse inverse matrix
					inverseMatrix, parseErr := parseMatrixFromBinary(inverseRecord.MatrixBinary)
					if parseErr != nil {
						log.Printf("❌ [BULK-INVERSE] Ters matris parse edilemedi (ID %d): %v", inverseRecord.ID, parseErr)
					} else {
						// Run algorithms in background
						go func(matrixID int, matrix Matrix, title string) {
							var boyarResult, paarResult, slpResult *AlgResult

							// Boyar SLP
							if boyar := NewBoyarSLP(10); boyar != nil {
								if result, err := boyar.Solve(matrix); err == nil {
									boyarResult = &result
									log.Printf("✅ [BULK-INVERSE-BOYAR] %s tamamlandı - XOR: %d", title, result.XorCount)
								} else {
									log.Printf("❌ [BULK-INVERSE-BOYAR] %s hatası: %v", title, err)
								}
							}

							// Paar Algorithm
							if paar := NewPaarAlgorithm(); paar != nil {
								if result, err := paar.Solve(matrix); err == nil {
									paarResult = &result
									log.Printf("✅ [BULK-INVERSE-PAAR] %s tamamlandı - XOR: %d", title, result.XorCount)
								} else {
									log.Printf("❌ [BULK-INVERSE-PAAR] %s hatası: %v", title, err)
								}
							}

							// SLP Heuristic
							if slp := NewSLPHeuristic(); slp != nil {
								if result, err := slp.Solve(matrix); err == nil {
									slpResult = &result
									log.Printf("✅ [BULK-INVERSE-SLP] %s tamamlandı - XOR: %d", title, result.XorCount)
								} else {
									log.Printf("❌ [BULK-INVERSE-SLP] %s hatası: %v", title, err)
								}
							}

							// Update database
							if updateErr := db.UpdateMatrixResults(matrixID, boyarResult, paarResult, slpResult, nil); updateErr != nil {
								log.Printf("❌ [BULK-INVERSE-UPDATE] %s sonuçları güncellenemedi: %v", title, updateErr)
							} else {
								log.Printf("✅ [BULK-INVERSE-UPDATE] %s sonuçları güncellendi", title)
							}
						}(inverseRecord.ID, inverseMatrix, inverseRecord.Title)
					}
				}
			}

			results = append(results, result)

			// Small delay between operations
			time.Sleep(100 * time.Millisecond)
		}

		log.Printf("🎉 [BULK-INVERSE] Toplu ters alma tamamlandı - Başarılı: %d, Hatalı: %d", 
			successCount, errorCount)
	}()

	// Return immediate response
	response := BulkInverseResponse{
		ProcessedCount: 0,
		SuccessCount:   0,
		ErrorCount:     0,
		Results:        []BulkInverseResult{},
		Message:        fmt.Sprintf("%d matris için ters alma işlemi başlatıldı", len(matrices)),
	}

	json.NewEncoder(w).Encode(response)
}

// getInversePairsHandler handles requests for inverse matrix pairs
func getInversePairsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	groupFilter := r.URL.Query().Get("group")
	maxCombinedXorStr := r.URL.Query().Get("max_combined_xor")
	sortOrder := r.URL.Query().Get("sort")

	// Default values
	page := 1
	limit := 25

	// Parse page
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse max combined XOR
	var maxCombinedXor *int
	if maxCombinedXorStr != "" {
		if val, err := strconv.Atoi(maxCombinedXorStr); err == nil && val > 0 {
			maxCombinedXor = &val
		}
	}

	// Default sort order
	if sortOrder == "" {
		sortOrder = "combined_asc"
	}

	log.Printf("📊 [INVERSE-PAIRS] Ters matris çiftleri isteniyor - Sayfa: %d, Limit: %d, Grup: '%s', Max Combined XOR: %v, Sıralama: %s", 
		page, limit, groupFilter, maxCombinedXor, sortOrder)

	// Get inverse pairs from database
	pairs, total, err := db.GetInversePairs(page, limit, groupFilter, maxCombinedXor, sortOrder)
	if err != nil {
		log.Printf("❌ [INVERSE-PAIRS] Ters matris çiftleri alınamadı: %v", err)
		http.Error(w, "Ters matris çiftleri alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	response := GetInversePairsResponse{
		Pairs:      pairs,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	log.Printf("✅ [INVERSE-PAIRS] %d ters matris çifti döndürüldü (Toplam: %d, Sayfa: %d/%d)", 
		len(pairs), total, page, totalPages)

	json.NewEncoder(w).Encode(response)
}

// sbpHandler handles SBP algorithm requests
func sbpHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[SBP] İstek başladı - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		log.Printf("[SBP] OPTIONS isteği işlendi")
		return
	}

	if r.Method != "POST" {
		log.Printf("[SBP] HATA: Geçersiz method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Matrices [][][]string `json:"matrices"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[SBP] HATA: JSON decode hatası: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	log.Printf("[SBP] %d matris alındı", len(request.Matrices))

	var results []map[string]interface{}
	for i, matrix := range request.Matrices {
		log.Printf("[SBP] Matris %d işleniyor (%dx%d)", i+1, len(matrix), len(matrix[0]))
		
		if len(matrix) == 0 {
			log.Printf("[SBP] HATA: Matris %d boş", i+1)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        "Empty matrix",
				"xor_count":    0,
				"depth":        0,
				"program":      []string{},
			})
			continue
		}

		sbp := NewSBPAlgorithm(10) // depth limit 10
		err := sbp.ReadTargetMatrix(matrix)
		if err != nil {
			log.Printf("[SBP] HATA: Matris %d okuma hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"depth":        0,
				"program":      []string{},
			})
			continue
		}

		err = sbp.InitBase()
		if err != nil {
			log.Printf("[SBP] HATA: Matris %d init hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"depth":        0,
				"program":      []string{},
			})
			continue
		}

		result, err := sbp.Solve(matrix)
		if err != nil {
			log.Printf("[SBP] HATA: Matris %d solve hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"depth":        0,
				"program":      []string{},
			})
			continue
		}

		log.Printf("[SBP] Matris %d başarıyla işlendi - XOR: %d, Derinlik: %d", i+1, result.XorCount, result.Depth)
		results = append(results, map[string]interface{}{
			"matrix_index": i,
			"xor_count":    result.XorCount,
			"depth":        result.Depth,
			"program":      result.Program,
		})
	}

	duration := time.Since(startTime)
	log.Printf("[SBP] İstek tamamlandı - Süre: %v, Sonuç sayısı: %d", duration, len(results))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"algorithm": "SBP",
		"results":   results,
	})
}

// getMissingAlgorithmsHandler returns the count of matrices with missing algorithm results
func getMissingAlgorithmsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	missingCount, totalCount, err := db.GetMissingAlgorithmsCount()
	if err != nil {
		log.Printf("❌ [API] GetMissingAlgorithmsCount error: %v", err)
		http.Error(w, "Eksik algoritma sayısı alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := MissingAlgorithmsResponse{
		MissingCount: missingCount,
		TotalCount:   totalCount,
		Message:      fmt.Sprintf("%d matris için algoritma sonuçları eksik", missingCount),
	}

	json.NewEncoder(w).Encode(response)
} 