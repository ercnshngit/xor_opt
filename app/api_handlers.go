package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	titleFilter := r.URL.Query().Get("title")

	// Parse range filter parameters
	hamXorMinStr := r.URL.Query().Get("ham_xor_min")
	hamXorMaxStr := r.URL.Query().Get("ham_xor_max")
	boyarXorMinStr := r.URL.Query().Get("boyar_xor_min")
	boyarXorMaxStr := r.URL.Query().Get("boyar_xor_max")
	paarXorMinStr := r.URL.Query().Get("paar_xor_min")
	paarXorMaxStr := r.URL.Query().Get("paar_xor_max")
	slpXorMinStr := r.URL.Query().Get("slp_xor_min")
	slpXorMaxStr := r.URL.Query().Get("slp_xor_max")

	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse range filter values
	var hamXorMin, hamXorMax, boyarXorMin, boyarXorMax, paarXorMin, paarXorMax, slpXorMin, slpXorMax *int

	if hamXorMinStr != "" {
		if val, err := strconv.Atoi(hamXorMinStr); err == nil {
			hamXorMin = &val
		}
	}
	if hamXorMaxStr != "" {
		if val, err := strconv.Atoi(hamXorMaxStr); err == nil {
			hamXorMax = &val
		}
	}
	if boyarXorMinStr != "" {
		if val, err := strconv.Atoi(boyarXorMinStr); err == nil {
			boyarXorMin = &val
		}
	}
	if boyarXorMaxStr != "" {
		if val, err := strconv.Atoi(boyarXorMaxStr); err == nil {
			boyarXorMax = &val
		}
	}
	if paarXorMinStr != "" {
		if val, err := strconv.Atoi(paarXorMinStr); err == nil {
			paarXorMin = &val
		}
	}
	if paarXorMaxStr != "" {
		if val, err := strconv.Atoi(paarXorMaxStr); err == nil {
			paarXorMax = &val
		}
	}
	if slpXorMinStr != "" {
		if val, err := strconv.Atoi(slpXorMinStr); err == nil {
			slpXorMin = &val
		}
	}
	if slpXorMaxStr != "" {
		if val, err := strconv.Atoi(slpXorMaxStr); err == nil {
			slpXorMax = &val
		}
	}

	matrices, total, err := db.GetMatrices(page, limit, titleFilter, hamXorMin, hamXorMax, boyarXorMin, boyarXorMax, paarXorMin, paarXorMax, slpXorMin, slpXorMax)
	if err != nil {
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

// recalculateHandler recalculates algorithms for a matrix
func recalculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RecalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		return
	}

	// Get the matrix from database
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

	// Recalculate Ham XOR
	newHamXor := calculateHammingXOR(matrix)
	
	// Update Ham XOR in database
	_, err = db.db.Exec("UPDATE matrix_records SET ham_xor_count = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", newHamXor, req.MatrixID)
	if err != nil {
		http.Error(w, "Ham XOR güncellenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var boyarResult, paarResult, slpResult *AlgResult

	// Run requested algorithms
	for _, algorithm := range req.Algorithms {
		switch strings.ToLower(algorithm) {
		case "boyar":
			boyar := NewBoyarSLP(10) // depth limit
			result, err := boyar.Solve(matrix)
			if err == nil {
				boyarResult = &result
			}
		case "paar":
			paar := NewPaarAlgorithm()
			result, err := paar.Solve(matrix)
			if err == nil {
				paarResult = &result
			}
		case "slp":
			slp := NewSLPHeuristic()
			result, err := slp.Solve(matrix)
			if err == nil {
				slpResult = &result
			}
		}
	}

	// Update database with results
	err = db.UpdateMatrixResults(req.MatrixID, boyarResult, paarResult, slpResult)
	if err != nil {
		http.Error(w, "Sonuçlar güncellenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated record
	updatedRecord, err := db.GetMatrixByID(req.MatrixID)
	if err != nil {
		http.Error(w, "Güncellenmiş matris alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedRecord)
}

// parseMatrixFromBinary parses a matrix from its binary string representation
func parseMatrixFromBinary(binaryStr string) (Matrix, error) {
	lines := strings.Split(binaryStr, "\n")
	var matrix Matrix

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove brackets and split by spaces
		line = strings.Trim(line, "[]")
		elements := strings.Fields(line)
		
		if len(elements) > 0 {
			matrix = append(matrix, elements)
		}
	}

	return matrix, nil
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