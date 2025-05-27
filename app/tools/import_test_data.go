package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// Read test data file
	data, err := ioutil.ReadFile("web/test_data.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	content := string(data)
	
	// Parse matrices from the content
	matrices := parseTestDataForImport(content)
	
	fmt.Printf("Found %d matrices\n", len(matrices))
	
	// Send each matrix to the API
	for i, matrix := range matrices {
		err := sendMatrixToAPIForImport(matrix.Title, matrix.Matrix)
		if err != nil {
			fmt.Printf("Error sending matrix %d: %v\n", i+1, err)
		} else {
			fmt.Printf("Successfully sent matrix %d: %s\n", i+1, matrix.Title)
		}
	}
}

type TestMatrixImport struct {
	Title  string
	Matrix [][]string
}

func parseTestDataForImport(content string) []TestMatrixImport {
	var matrices []TestMatrixImport
	
	// Split by "-----" separator
	sections := strings.Split(content, "-----")
	
	for _, section := range sections {
		section = strings.TrimSpace(section)
		if section == "" {
			continue
		}
		
		lines := strings.Split(section, "\n")
		if len(lines) < 3 {
			continue
		}
		
		// Extract title (first line)
		title := strings.TrimSpace(lines[0])
		if title == "" {
			continue
		}
		
		// Find matrix data (lines with brackets)
		var matrixLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				matrixLines = append(matrixLines, line)
			}
		}
		
		if len(matrixLines) == 0 {
			continue
		}
		
		// Parse matrix
		var matrix [][]string
		for _, line := range matrixLines {
			// Remove brackets and split by spaces
			line = strings.Trim(line, "[]")
			elements := strings.Fields(line)
			if len(elements) > 0 {
				matrix = append(matrix, elements)
			}
		}
		
		if len(matrix) > 0 {
			matrices = append(matrices, TestMatrixImport{
				Title:  title,
				Matrix: matrix,
			})
		}
	}
	
	return matrices
}

func sendMatrixToAPIForImport(title string, matrix [][]string) error {
	// Prepare request data
	requestData := map[string]interface{}{
		"title":  title,
		"matrix": matrix,
	}
	
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return err
	}
	
	// Send POST request
	resp, err := http.Post("http://localhost:3000/api/matrices/process", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}
	
	return nil
} 