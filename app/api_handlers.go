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

// getGroupsHandler retrieves all unique group names
func getGroupsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Add cache headers for better performance
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minute cache for groups
	
	groups, err := db.GetGroups()
	if err != nil {
		log.Printf("❌ [API] GetGroups error: %v", err)
		http.Error(w, "Gruplar alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ [API] GetGroups completed: returned %d groups", len(groups))
	json.NewEncoder(w).Encode(groups)
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
		req.Algorithms = []string{"boyar", "paar", "slp"}
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
		log.Printf("Matris %d için algoritma hesaplama başlatıldı", req.MatrixID)

		var boyarResult, paarResult, slpResult *AlgResult

		// Run requested algorithms
		for _, algorithm := range req.Algorithms {
			switch strings.ToLower(algorithm) {
			case "boyar":
				boyar := NewBoyarSLP(10)
				if result, err := boyar.Solve(matrix); err == nil {
					boyarResult = &result
				} else {
					log.Printf("Boyar algoritması hatası (ID %d): %v", req.MatrixID, err)
				}
			case "paar":
				paar := NewPaarAlgorithm()
				if result, err := paar.Solve(matrix); err == nil {
					paarResult = &result
				} else {
					log.Printf("Paar algoritması hatası (ID %d): %v", req.MatrixID, err)
				}
			case "slp":
				slp := NewSLPHeuristic()
				if result, err := slp.Solve(matrix); err == nil {
					slpResult = &result
				} else {
					log.Printf("SLP algoritması hatası (ID %d): %v", req.MatrixID, err)
				}
			}
		}

		// Update database with results
		err = db.UpdateMatrixResults(req.MatrixID, boyarResult, paarResult, slpResult)
		if err != nil {
			log.Printf("Algoritma sonuçları güncellenemedi (ID %d): %v", req.MatrixID, err)
		} else {
			log.Printf("Matris %d algoritmaları tamamlandı", req.MatrixID)
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
	var boyarResult, paarResult, slpResult *AlgResult

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

	// Update database with results
	err = db.UpdateMatrixResults(record.ID, boyarResult, paarResult, slpResult)
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
		req.Algorithms = []string{"boyar", "paar", "slp"}
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
		for i, matrix := range matrices {
			log.Printf("Toplu hesaplama: Matris %d/%d (ID: %d) işleniyor...", i+1, len(matrices), matrix.ID)
			
			// Parse matrix from binary string
			matrixData, err := parseMatrixFromBinary(matrix.MatrixBinary)
			if err != nil {
				log.Printf("Matris parse hatası (ID %d): %v", matrix.ID, err)
				continue
			}

			// Recalculate Ham XOR
			newHamXor := calculateHammingXOR(matrixData)
			
			// Update Ham XOR in database
			_, err = db.db.Exec("UPDATE matrix_records SET ham_xor_count = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", newHamXor, matrix.ID)
			if err != nil {
				log.Printf("Ham XOR güncellenemedi (ID %d): %v", matrix.ID, err)
			}

			var boyarResult, paarResult, slpResult *AlgResult

			// Run requested algorithms
			for _, algorithm := range req.Algorithms {
				switch strings.ToLower(algorithm) {
				case "boyar":
					boyar := NewBoyarSLP(10)
					if result, err := boyar.Solve(matrixData); err == nil {
						boyarResult = &result
					} else {
						log.Printf("Boyar algoritması hatası (ID %d): %v", matrix.ID, err)
					}
				case "paar":
					paar := NewPaarAlgorithm()
					if result, err := paar.Solve(matrixData); err == nil {
						paarResult = &result
					} else {
						log.Printf("Paar algoritması hatası (ID %d): %v", matrix.ID, err)
					}
				case "slp":
					slp := NewSLPHeuristic()
					if result, err := slp.Solve(matrixData); err == nil {
						slpResult = &result
					} else {
						log.Printf("SLP algoritması hatası (ID %d): %v", matrix.ID, err)
					}
				}
			}

			// Update database with results
			err = db.UpdateMatrixResults(matrix.ID, boyarResult, paarResult, slpResult)
			if err != nil {
				log.Printf("Algoritma sonuçları güncellenemedi (ID %d): %v", matrix.ID, err)
			} else {
				log.Printf("Matris %d algoritmaları tamamlandı", matrix.ID)
			}
		}
		log.Printf("Toplu algoritma hesaplama tamamlandı: %d matris işlendi", len(matrices))
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

// getMatrixByIdHandler retrieves a specific matrix by ID
func getMatrixByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Geçersiz ID formatı", http.StatusBadRequest)
		return
	}

	matrix, err := db.GetMatrixByID(id)
	if err != nil {
		http.Error(w, "Matris bulunamadı: "+err.Error(), http.StatusNotFound)
		return
	}

	// Get duplicate matrices
	duplicates, err := db.GetDuplicateMatrices(id)
	if err != nil {
		log.Printf("Duplicate matrisler alınamadı: %v", err)
		duplicates = []*DuplicateMatrixRecord{} // Empty slice if error
	}

	// Create response with duplicates
	response := struct {
		*MatrixRecord
		Duplicates []*DuplicateMatrixRecord `json:"duplicates"`
	}{
		MatrixRecord: matrix,
		Duplicates:   duplicates,
	}

	json.NewEncoder(w).Encode(response)
}

// getDuplicateMatricesHandler returns all duplicate matrices
func getDuplicateMatricesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	duplicatesMap, err := db.GetAllDuplicateMatrices()
	if err != nil {
		http.Error(w, "Duplicate matrisler alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert map to a more structured response
	type DuplicateGroup struct {
		OriginalID  int                      `json:"original_id"`
		Duplicates  []*DuplicateMatrixRecord `json:"duplicates"`
		Count       int                      `json:"count"`
	}

	var response []DuplicateGroup
	for originalID, duplicates := range duplicatesMap {
		if len(duplicates) > 0 { // Only include if there are duplicates
			response = append(response, DuplicateGroup{
				OriginalID: originalID,
				Duplicates: duplicates,
				Count:      len(duplicates),
			})
		}
	}

	json.NewEncoder(w).Encode(response)
} 