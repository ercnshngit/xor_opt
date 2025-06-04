package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Matrix represents a binary matrix as strings
type Matrix [][]string

// Request represents the API request format
type Request struct {
	Matrices []Matrix `json:"matrices"`
}

// Response represents the API response format
type Response struct {
	Algorithm string      `json:"algorithm"`
	Results   []AlgResult `json:"results"`
	Error     string      `json:"error,omitempty"`
}

// AlgResult represents result for one matrix
type AlgResult struct {
	MatrixIndex int      `json:"matrix_index"`
	XorCount    int      `json:"xor_count"`
	Program     []string `json:"program"`
	Depth       int      `json:"depth,omitempty"`
}

// Constants for array sizes - optimized for 4-core 16GB server
const (
	MAX_ARRAY_SIZE = 4000  // Increased for better performance on 16GB RAM
	MAX_ITERATIONS = 50000 // Increased for more thorough calculations
)

// BoyarSLP implementation
type BoyarSLP struct {
	NumInputs    int
	DepthLimit   int
	NumTargets   int
	ProgramSize  int
	Target       []uint64
	Dist         []int
	NDist        []int
	Base         []uint64
	BaseSize     int
	TargetsFound int
	Result       []string
	Depth        []int
	MaxDepth     int
}

func NewBoyarSLP(depthLimit int) *BoyarSLP {
	return &BoyarSLP{
		DepthLimit: depthLimit,
		Target:     make([]uint64, MAX_ARRAY_SIZE),
		Dist:       make([]int, MAX_ARRAY_SIZE),
		NDist:      make([]int, MAX_ARRAY_SIZE),
		Base:       make([]uint64, MAX_ARRAY_SIZE),
		Result:     make([]string, MAX_ARRAY_SIZE),
		Depth:      make([]int, MAX_ARRAY_SIZE),
	}
}

func (b *BoyarSLP) ReadTargetMatrix(matrix Matrix) error {
	b.NumTargets = len(matrix)
	if b.NumTargets == 0 {
		return fmt.Errorf("matris boş")
	}
	if b.NumTargets >= MAX_ARRAY_SIZE {
		return fmt.Errorf("matris çok büyük: %d >= %d", b.NumTargets, MAX_ARRAY_SIZE)
	}
	b.NumInputs = len(matrix[0])
	if b.NumInputs >= MAX_ARRAY_SIZE {
		return fmt.Errorf("matris genişliği çok büyük: %d >= %d", b.NumInputs, MAX_ARRAY_SIZE)
	}

	for i := 0; i < b.NumTargets; i++ {
		var powerOfTwo uint64 = 1
		b.Target[i] = 0
		b.Dist[i] = -1

		for j := 0; j < b.NumInputs; j++ {
			bit, _ := strconv.Atoi(strings.TrimSpace(matrix[i][j]))
			if bit == 1 {
				b.Dist[i]++
				b.Target[i] = b.Target[i] + powerOfTwo
			}
			powerOfTwo = powerOfTwo * 2
		}
	}
	return nil
}

func (b *BoyarSLP) InitBase() error {
	b.TargetsFound = 0
	b.ProgramSize = 0
	b.Result = make([]string, MAX_ARRAY_SIZE)
	b.Base[0] = 1
	b.Depth[0] = 0
	b.MaxDepth = 0

	for i := 1; i < b.NumInputs; i++ {
		if i >= MAX_ARRAY_SIZE {
			return fmt.Errorf("base array overflow: %d >= %d", i, MAX_ARRAY_SIZE)
		}
		b.Base[i] = 2 * b.Base[i-1]
		b.Depth[i] = 0
	}
	b.BaseSize = b.NumInputs

	for i := 0; i < b.NumTargets; i++ {
		if b.Dist[i] == 0 {
			b.TargetsFound++
			for j := 0; j < b.NumInputs; j++ {
				if b.Base[j] == b.Target[i] {
					b.Result = append(b.Result, fmt.Sprintf("y%d = x%d", i, j))
					break
				}
			}
		}
	}
	return nil
}

func (b *BoyarSLP) isTarget(x uint64) bool {
	for i := 0; i < b.NumTargets; i++ {
		if x == b.Target[i] {
			return true
		}
	}
	return false
}

func (b *BoyarSLP) isBase(x uint64) bool {
	if x == 0 {
		return false
	}
	for i := 0; i < b.BaseSize; i++ {
		if x == b.Base[i] {
			return true
		}
	}
	return false
}

func (b *BoyarSLP) max(a, c int) int {
	if a > c {
		return a
	}
	return c
}

func (b *BoyarSLP) reachable(T uint64, K, S int, L uint64) bool {
	if (b.BaseSize-S) < K {
		return false
	}
	if L < 1 {
		return false
	}
	if K == 0 {
		return false
	}
	if K == 1 {
		for i := S; i < b.BaseSize; i++ {
			if T == b.Base[i] && uint64(math.Pow(2, float64(b.Depth[i]))) <= L {
				return true
			}
		}
		return false
	}

	if b.reachable(T^b.Base[S], K-1, S+1, L-uint64(math.Pow(2, float64(b.Depth[S])))) {
		return true
	}
	if b.reachable(T, K, S+1, L) {
		return true
	}
	return false
}

func (b *BoyarSLP) NewDistance(u int, newBase uint64, depthNewBase uint64) int {
	if b.Target[u] == 0 {
		return 0
	}
	if b.isBase(b.Target[u]) || newBase == b.Target[u] {
		return 0
	}
	if b.reachable(b.Target[u]^newBase, b.Dist[u]-1, 0, uint64(math.Pow(2, float64(b.DepthLimit)))-depthNewBase) {
		return b.Dist[u] - 1
	}
	return b.Dist[u]
}

func (b *BoyarSLP) TotalDistance(newBase uint64, depthNewBase uint64) int {
	D := 0
	for i := 0; i < b.NumTargets; i++ {
		t := b.NewDistance(i, newBase, depthNewBase)
		b.NDist[i] = t
		D = D + t
	}
	return D
}

func (b *BoyarSLP) EasyMove() bool {
	t := -1
	for i := 0; i < b.NumTargets; i++ {
		if b.Dist[i] == 1 {
			t = i
			break
		}
	}
	if t == -1 {
		return false
	}

	// Array sınır kontrolü
	if b.BaseSize >= MAX_ARRAY_SIZE-1 {
		return false
	}

	newBase := b.Target[t]
	b.Base[b.BaseSize] = newBase
	b.BaseSize++

	depthNewBase := uint64(math.Pow(2, float64(b.DepthLimit)))
	for i := 0; i < b.BaseSize; i++ {
		for j := i + 1; j < b.BaseSize; j++ {
			if (b.Base[i]^b.Base[j]) == b.Base[b.BaseSize-1] {
				newDepth := uint64(math.Pow(2, float64(b.max(b.Depth[i], b.Depth[j])+1)))
				if depthNewBase > newDepth {
					depthNewBase = newDepth
				}
			}
		}
	}

	for u := 0; u < b.NumTargets; u++ {
		b.Dist[u] = b.NewDistance(u, newBase, depthNewBase)
	}
	b.ProgramSize++
	b.TargetsFound++

	// Find which bases created this target
	for i := 0; i < b.BaseSize; i++ {
		for j := i + 1; j < b.BaseSize; j++ {
			if (b.Base[i]^b.Base[j]) == b.Base[b.BaseSize-1] {
				b.Depth[b.BaseSize-1] = b.max(b.Depth[i], b.Depth[j]) + 1
				if b.Depth[b.BaseSize-1] > b.MaxDepth {
					b.MaxDepth = b.Depth[b.BaseSize-1]
				}
				var iStr, jStr string
				if i < b.NumInputs {
					iStr = fmt.Sprintf("x%d", i)
				} else {
					iStr = fmt.Sprintf("t%d", i-b.NumInputs+1)
				}
				if j < b.NumInputs {
					jStr = fmt.Sprintf("x%d", j)
				} else {
					jStr = fmt.Sprintf("t%d", j-b.NumInputs+1)
				}
				b.Result = append(b.Result, fmt.Sprintf("t%d = %s + %s * y%d (%d)", b.ProgramSize, iStr, jStr, t, b.Depth[b.BaseSize-1]))
				return true
			}
		}
	}
	return true
}

func (b *BoyarSLP) PickNewBaseElement() bool {
	// Array sınır kontrolü
	if b.BaseSize >= MAX_ARRAY_SIZE-1 {
		return false
	}

	minDistance := b.BaseSize * b.NumTargets
	oldNorm := 0
	var bestI, bestJ int
	var theBest uint64
	bestDist := make([]int, b.NumTargets)

	for i := 0; i < b.BaseSize-1; i++ {
		if b.Depth[i]+1 >= b.DepthLimit {
			continue
		}
		for j := i + 1; j < b.BaseSize; j++ {
			if b.Depth[j]+1 >= b.DepthLimit {
				continue
			}
			newBase := b.Base[i] ^ b.Base[j]
			if newBase == 0 || b.isBase(newBase) {
				continue
			}

			depthNewBase := uint64(math.Pow(2, float64(b.max(b.Depth[i], b.Depth[j])+1)))
			thisDist := b.TotalDistance(newBase, depthNewBase)

			if thisDist <= minDistance {
				thisNorm := 0
				for k := 0; k < b.NumTargets; k++ {
					d := b.NDist[k]
					thisNorm = thisNorm + d*d
				}

				if thisDist < minDistance || thisNorm > oldNorm {
					bestI = i
					bestJ = j
					theBest = newBase
					copy(bestDist, b.NDist[:b.NumTargets])
					minDistance = thisDist
					oldNorm = thisNorm
				}
			}
		}
	}

	for i := 0; i < b.NumTargets; i++ {
		b.Dist[i] = bestDist[i]
	}

	b.Base[b.BaseSize] = theBest
	b.Depth[b.BaseSize] = b.max(b.Depth[bestI], b.Depth[bestJ]) + 1
	if b.Depth[b.BaseSize] > b.MaxDepth {
		b.MaxDepth = b.Depth[b.BaseSize]
	}
	b.BaseSize++
	b.ProgramSize++

	var iStr, jStr string
	if bestI < b.NumInputs {
		iStr = fmt.Sprintf("x%d", bestI)
	} else {
		iStr = fmt.Sprintf("t%d", bestI-b.NumInputs+1)
	}
	if bestJ < b.NumInputs {
		jStr = fmt.Sprintf("x%d", bestJ)
	} else {
		jStr = fmt.Sprintf("t%d", bestJ-b.NumInputs+1)
	}
	b.Result = append(b.Result, fmt.Sprintf("t%d = %s + %s (%d)", b.ProgramSize, iStr, jStr, b.Depth[b.BaseSize-1]))

	if b.isTarget(theBest) {
		b.TargetsFound++
	}
	return true
}

func (b *BoyarSLP) Solve(matrix Matrix) (AlgResult, error) {
	err := b.ReadTargetMatrix(matrix)
	if err != nil {
		return AlgResult{}, err
	}
	
	err = b.InitBase()
	if err != nil {
		return AlgResult{}, err
	}

	iterations := 0
	for b.TargetsFound < b.NumTargets && iterations < MAX_ITERATIONS {
		if !b.EasyMove() {
			if !b.PickNewBaseElement() {
				break // Array sınırına ulaşıldı
			}
		}
		iterations++
	}

	var program []string
	for _, res := range b.Result {
		if res != "" {
			program = append(program, res)
		}
	}

	return AlgResult{
		XorCount: b.ProgramSize,
		Program:  program,
		Depth:    b.MaxDepth,
	}, nil
}

// Paar Algorithm implementation
type PaarAlgorithm struct {
	NumInputs int
	Dim       int
}

func NewPaarAlgorithm() *PaarAlgorithm {
	return &PaarAlgorithm{}
}

func (p *PaarAlgorithm) ReadTargetMatrix(matrix Matrix) error {
	p.Dim = len(matrix)
	if p.Dim == 0 {
		return fmt.Errorf("matris boş")
	}
	if p.Dim >= MAX_ARRAY_SIZE {
		return fmt.Errorf("matris çok büyük: %d >= %d", p.Dim, MAX_ARRAY_SIZE)
	}
	p.NumInputs = len(matrix[0])
	if p.NumInputs >= MAX_ARRAY_SIZE {
		return fmt.Errorf("matris genişliği çok büyük: %d >= %d", p.NumInputs, MAX_ARRAY_SIZE)
	}
	return nil
}

func (p *PaarAlgorithm) InitBase() error {
	// PAAR algoritması için özel init işlemi
	return nil
}

func (p *PaarAlgorithm) hammingWeight(input uint64) int {
	count := 0
	for input != 0 {
		count += int(input & 1)
		input >>= 1
	}
	return count
}

func (p *PaarAlgorithm) Solve(matrix Matrix) (AlgResult, error) {
	err := p.ReadTargetMatrix(matrix)
	if err != nil {
		return AlgResult{}, err
	}
	
	err = p.InitBase()
	if err != nil {
		return AlgResult{}, err
	}

	// Convert matrix to uint64 array (column-wise like C++)
	inputMatrix := make([]uint64, p.Dim+200)
	for i := 0; i < p.Dim; i++ {
		var val uint64 = 0
		for j := 0; j < p.NumInputs; j++ {
			bit, _ := strconv.Atoi(strings.TrimSpace(matrix[j][i]))
			if bit == 1 {
				val |= (1 << (p.Dim - 1 - j))
			}
		}
		inputMatrix[i] = val & ((1 << p.Dim) - 1)
	}

	xorCount := 0
	numberOfColumns := p.Dim
	var program []string

	// Compute naive xor count
	for i := 0; i < p.Dim; i++ {
		xorCount += p.hammingWeight(inputMatrix[i])
	}
	xorCount -= p.Dim

	for {
		hwMax := 0
		var iMax, jMax int

		for i := 0; i < numberOfColumns; i++ {
			for j := i + 1; j < numberOfColumns; j++ {
				tmp := inputMatrix[i] & inputMatrix[j]
				hw := p.hammingWeight(tmp)
				if hw > hwMax {
					hwMax = hw
					iMax = i
					jMax = j
				}
			}
		}

		if hwMax <= 1 {
			break
		}

		newColumn := inputMatrix[iMax] & inputMatrix[jMax]
		inputMatrix[numberOfColumns] = newColumn
		inputMatrix[iMax] = (newColumn ^ ((1 << p.Dim) - 1)) & inputMatrix[iMax]
		inputMatrix[jMax] = (newColumn ^ ((1 << p.Dim) - 1)) & inputMatrix[jMax]
		xorCount -= (hwMax - 1)
		numberOfColumns++
		program = append(program, fmt.Sprintf("x%d = x%d + x%d", numberOfColumns-1, iMax, jMax))
	}

	// Generate output equations
	for i := 0; i < p.Dim; i++ {
		var equation strings.Builder
		equation.WriteString(fmt.Sprintf("y%d = ", i))
		first := true
		for j := 0; j < numberOfColumns; j++ {
			if (inputMatrix[j] & (1 << (p.Dim - 1 - i))) != 0 {
				if !first {
					equation.WriteString(" + ")
				}
				equation.WriteString(fmt.Sprintf("x%d", j))
				first = false
			}
		}
		if !first {
			program = append(program, equation.String())
		}
	}

	return AlgResult{
		XorCount: xorCount,
		Program:  program,
	}, nil
}

// SLP Heuristic implementation
type SLPHeuristic struct {
	NumInputs    int
	NumTargets   int
	XorCount     int
	Target       []uint64
	Dist         []int
	NDist        []int
	Base         []uint64
	Program      []string
	BaseSize     int
	TargetsFound int
}

func NewSLPHeuristic() *SLPHeuristic {
	return &SLPHeuristic{
		Target:  make([]uint64, MAX_ARRAY_SIZE),
		Dist:    make([]int, MAX_ARRAY_SIZE),
		NDist:   make([]int, MAX_ARRAY_SIZE),
		Base:    make([]uint64, MAX_ARRAY_SIZE),
		Program: make([]string, MAX_ARRAY_SIZE),
	}
}

func (s *SLPHeuristic) ReadTargetMatrix(matrix Matrix) error {
	s.NumTargets = len(matrix)
	if s.NumTargets == 0 {
		return fmt.Errorf("matris boş")
	}
	if s.NumTargets >= MAX_ARRAY_SIZE {
		return fmt.Errorf("matris çok büyük: %d >= %d", s.NumTargets, MAX_ARRAY_SIZE)
	}
	s.NumInputs = len(matrix[0])
	if s.NumInputs >= MAX_ARRAY_SIZE {
		return fmt.Errorf("matris genişliği çok büyük: %d >= %d", s.NumInputs, MAX_ARRAY_SIZE)
	}

	for i := 0; i < s.NumTargets; i++ {
		var powerOfTwo uint64 = 1
		s.Target[i] = 0
		s.Dist[i] = -1

		for j := 0; j < s.NumInputs; j++ {
			bit, _ := strconv.Atoi(strings.TrimSpace(matrix[i][j]))
			if bit == 1 {
				s.Dist[i]++
				s.Target[i] = s.Target[i] + powerOfTwo
			}
			powerOfTwo = powerOfTwo * 2
		}
	}
	return nil
}

func (s *SLPHeuristic) InitBase() error {
	s.TargetsFound = 0
	s.Base[0] = 1
	s.Program[0] = "x0"
	for i := 1; i < s.NumInputs; i++ {
		if i >= MAX_ARRAY_SIZE {
			return fmt.Errorf("base array overflow: %d >= %d", i, MAX_ARRAY_SIZE)
		}
		s.Base[i] = 2 * s.Base[i-1]
		s.Program[i] = fmt.Sprintf("x%d", i)
	}
	s.BaseSize = s.NumInputs

	for i := 0; i < s.NumTargets; i++ {
		if s.Dist[i] == 0 {
			s.TargetsFound++
		}
	}
	return nil
}

func (s *SLPHeuristic) isTarget(x uint64) bool {
	for i := 0; i < s.NumTargets; i++ {
		if x == s.Target[i] {
			return true
		}
	}
	return false
}

func (s *SLPHeuristic) isBase(x uint64) bool {
	if x == 0 {
		return false
	}
	for i := 0; i < s.BaseSize; i++ {
		if x == s.Base[i] {
			return true
		}
	}
	return false
}

func (s *SLPHeuristic) reachable(T uint64, K, S int) bool {
	if (s.BaseSize-S) < K {
		return false
	}
	if K == 0 {
		return false
	}
	if K == 1 {
		for i := S; i < s.BaseSize; i++ {
			if T == s.Base[i] {
				return true
			}
		}
		return false
	}

	if s.reachable(T^s.Base[S], K-1, S+1) {
		return true
	}
	if s.reachable(T, K, S+1) {
		return true
	}
	return false
}

func (s *SLPHeuristic) NewDistance(u int, newBase uint64) int {
	if s.isBase(s.Target[u]) || newBase == s.Target[u] {
		return 0
	}
	if s.reachable(s.Target[u]^newBase, s.Dist[u]-1, 0) {
		return s.Dist[u] - 1
	}
	return s.Dist[u]
}

func (s *SLPHeuristic) TotalDistance(newBase uint64) int {
	D := 0
	for i := 0; i < s.NumTargets; i++ {
		t := s.NewDistance(i, newBase)
		s.NDist[i] = t
		D = D + t
	}
	return D
}

func (s *SLPHeuristic) EasyMove() bool {
	t := -1
	for i := 0; i < s.NumTargets; i++ {
		if s.Dist[i] == 1 {
			t = i
			break
		}
	}
	if t == -1 {
		return false
	}

	// Array sınır kontrolü
	if s.BaseSize >= MAX_ARRAY_SIZE-1 {
		return false
	}

	newBase := s.Target[t]
	for u := 0; u < s.NumTargets; u++ {
		s.Dist[u] = s.NewDistance(u, newBase)
	}

	s.Base[s.BaseSize] = newBase

	// Find which lines in Base caused this
	var a, b string
	for i := 0; i < s.BaseSize; i++ {
		for j := i + 1; j < s.BaseSize; j++ {
			if (s.Base[i] ^ s.Base[j]) == s.Target[t] {
				a = strings.Split(s.Program[i], " ")[0]
				b = strings.Split(s.Program[j], " ")[0]
				break
			}
		}
	}

	s.Program[s.BaseSize] = fmt.Sprintf("y%d = %s + %s", t, a, b)
	s.BaseSize++
	s.XorCount++
	s.TargetsFound++
	return true
}

func (s *SLPHeuristic) PickNewBaseElement() bool {
	// Array sınır kontrolü
	if s.BaseSize >= MAX_ARRAY_SIZE-1 {
		return false
	}

	minDistance := s.BaseSize * s.NumTargets
	oldNorm := 0
	var bestI, bestJ int
	var theBest uint64
	bestDist := make([]int, s.NumTargets)

	for i := 0; i < s.BaseSize-1; i++ {
		for j := i + 1; j < s.BaseSize; j++ {
			newBase := s.Base[i] ^ s.Base[j]
			if newBase == 0 || s.isBase(newBase) {
				continue
			}

			thisDist := s.TotalDistance(newBase)
			if thisDist <= minDistance {
				thisNorm := 0
				for k := 0; k < s.NumTargets; k++ {
					d := s.NDist[k]
					thisNorm = thisNorm + d*d
				}

				if thisDist < minDistance || thisNorm > oldNorm {
					bestI = i
					bestJ = j
					theBest = newBase
					copy(bestDist, s.NDist[:s.NumTargets])
					minDistance = thisDist
					oldNorm = thisNorm
				}
			}
		}
	}

	for i := 0; i < s.NumTargets; i++ {
		s.Dist[i] = bestDist[i]
	}

	s.Base[s.BaseSize] = theBest
	a := strings.Split(s.Program[bestI], " ")[0]
	b := strings.Split(s.Program[bestJ], " ")[0]
	s.Program[s.BaseSize] = fmt.Sprintf("t%d = %s + %s", s.XorCount, a, b)
	s.BaseSize++
	s.XorCount++

	if s.isTarget(theBest) {
		s.TargetsFound++
	}
	return true
}

func (s *SLPHeuristic) Solve(matrix Matrix) (AlgResult, error) {
	err := s.ReadTargetMatrix(matrix)
	if err != nil {
		return AlgResult{}, err
	}
	
	err = s.InitBase()
	if err != nil {
		return AlgResult{}, err
	}
	s.XorCount = 0

	iterations := 0
	for s.TargetsFound < s.NumTargets && iterations < MAX_ITERATIONS {
		if !s.EasyMove() {
			if !s.PickNewBaseElement() {
				break // Array sınırına ulaşıldı
			}
		}
		iterations++
	}

	var program []string
	for j := 0; j < s.XorCount; j++ {
		if s.NumInputs+j < MAX_ARRAY_SIZE && s.Program[s.NumInputs+j] != "" {
			program = append(program, s.Program[s.NumInputs+j])
		}
	}

	return AlgResult{
		XorCount: s.XorCount,
		Program:  program,
	}, nil
}

// SBP Algorithm implementation (based on the original SBP code)
type SBPAlgorithm struct {
	NumInputs    int
	DepthLimit   int
	NumTargets   int
	ProgramSize  int
	Target       []uint64
	Dist         []int
	NDist        []int
	Base         []uint64
	BaseSize     int
	TargetsFound int
	Result       []string
	Depth        []int
	MaxDepth     int
	MaxDist      int
}

func NewSBPAlgorithm(depthLimit int) *SBPAlgorithm {
	return &SBPAlgorithm{
		DepthLimit: depthLimit,
		Target:     make([]uint64, MAX_ARRAY_SIZE),
		Dist:       make([]int, MAX_ARRAY_SIZE),
		NDist:      make([]int, MAX_ARRAY_SIZE),
		Base:       make([]uint64, MAX_ARRAY_SIZE),
		Result:     make([]string, MAX_ARRAY_SIZE),
		Depth:      make([]int, MAX_ARRAY_SIZE),
	}
}

func (s *SBPAlgorithm) ReadTargetMatrix(matrix Matrix) error {
	s.NumTargets = len(matrix)
	if s.NumTargets == 0 {
		return fmt.Errorf("matris boş")
	}
	if s.NumTargets >= MAX_ARRAY_SIZE {
		return fmt.Errorf("matris çok büyük: %d >= %d", s.NumTargets, MAX_ARRAY_SIZE)
	}
	s.NumInputs = len(matrix[0])
	if s.NumInputs >= MAX_ARRAY_SIZE {
		return fmt.Errorf("matris genişliği çok büyük: %d >= %d", s.NumInputs, MAX_ARRAY_SIZE)
	}

	s.MaxDist = 0
	for i := 0; i < s.NumTargets; i++ {
		var powerOfTwo uint64 = 1
		s.Target[i] = 0
		s.Dist[i] = -1

		for j := 0; j < s.NumInputs; j++ {
			bit, _ := strconv.Atoi(strings.TrimSpace(matrix[i][j]))
			if bit == 1 {
				s.Dist[i]++
				s.Target[i] = s.Target[i] + powerOfTwo
			}
			powerOfTwo = powerOfTwo * 2
		}
		if s.Dist[i] > s.MaxDist {
			s.MaxDist = s.Dist[i]
		}
	}
	return nil
}

func (s *SBPAlgorithm) InitBase() error {
	s.TargetsFound = 0
	s.ProgramSize = 0
	s.Result = make([]string, MAX_ARRAY_SIZE)
	s.Base[0] = 1
	s.Depth[0] = 0
	s.MaxDepth = 0

	for i := 1; i < s.NumInputs; i++ {
		if i >= MAX_ARRAY_SIZE {
			return fmt.Errorf("base array overflow: %d >= %d", i, MAX_ARRAY_SIZE)
		}
		s.Base[i] = 2 * s.Base[i-1]
		s.Depth[i] = 0
	}
	s.BaseSize = s.NumInputs

	for i := 0; i < s.NumTargets; i++ {
		if s.Dist[i] == 0 {
			s.TargetsFound++
			for j := 0; j < s.NumInputs; j++ {
				if s.Base[j] == s.Target[i] {
					s.Result = append(s.Result, fmt.Sprintf("y%d = x%d", i, j))
					break
				}
			}
		}
	}
	return nil
}

func (s *SBPAlgorithm) isTarget(x uint64) bool {
	for i := 0; i < s.NumTargets; i++ {
		if x == s.Target[i] {
			return true
		}
	}
	return false
}

func (s *SBPAlgorithm) isBase(x uint64) bool {
	if x == 0 {
		return false
	}
	for i := 0; i < s.BaseSize; i++ {
		if x == s.Base[i] {
			return true
		}
	}
	return false
}

func (s *SBPAlgorithm) max(a, b int) int {
	if s.Depth[a] > s.Depth[b] {
		return s.Depth[a]
	}
	return s.Depth[b]
}

func (s *SBPAlgorithm) reachable(T uint64, K, S int, L uint64) bool {
	if (s.BaseSize-S) < K {
		return false
	}
	if L < 1 {
		return false
	}
	if K == 0 {
		return false
	}
	if K == 1 {
		for i := S; i < s.BaseSize; i++ {
			if T == s.Base[i] && uint64(math.Pow(2, float64(s.Depth[i]))) <= L {
				return true
			}
		}
		return false
	}

	if s.reachable(T^s.Base[S], K-1, S+1, L-uint64(math.Pow(2, float64(s.Depth[S])))) {
		return true
	}
	if s.reachable(T, K, S+1, L) {
		return true
	}
	return false
}

func (s *SBPAlgorithm) NewDistance(u int, newBase uint64, depthNewBase uint64) int {
	if s.Target[u] == 0 {
		return 0
	}
	if s.isBase(s.Target[u]) || newBase == s.Target[u] {
		return 0
	}
	if s.reachable(s.Target[u]^newBase, s.Dist[u]-1, 0, uint64(math.Pow(2, float64(s.DepthLimit)))-depthNewBase) {
		return s.Dist[u] - 1
	}
	return s.Dist[u]
}

func (s *SBPAlgorithm) TotalDistance(newBase uint64, depthNewBase uint64) int {
	D := 0
	for i := 0; i < s.NumTargets; i++ {
		t := s.NewDistance(i, newBase, depthNewBase)
		s.NDist[i] = t
		D = D + t
	}
	return D
}

func (s *SBPAlgorithm) EasyMove() bool {
	t := -1
	for i := 0; i < s.NumTargets; i++ {
		if s.Dist[i] == 1 {
			t = i
			break
		}
	}
	if t == -1 {
		return false
	}

	// Array sınır kontrolü
	if s.BaseSize >= MAX_ARRAY_SIZE-1 {
		return false
	}

	newBase := s.Target[t]
	s.Base[s.BaseSize] = newBase
	s.BaseSize++

	depthNewBase := uint64(math.Pow(2, float64(s.DepthLimit)))
	for i := 0; i < s.BaseSize; i++ {
		for j := i + 1; j < s.BaseSize; j++ {
			if (s.Base[i]^s.Base[j]) == s.Base[s.BaseSize-1] {
				newDepth := uint64(math.Pow(2, float64(s.max(i, j)+1)))
				if depthNewBase > newDepth {
					depthNewBase = newDepth
				}
			}
		}
	}

	for u := 0; u < s.NumTargets; u++ {
		s.Dist[u] = s.NewDistance(u, newBase, depthNewBase)
	}
	s.ProgramSize++
	s.TargetsFound++

	// Find which bases created this target
	for i := 0; i < s.BaseSize; i++ {
		for j := i + 1; j < s.BaseSize; j++ {
			if (s.Base[i]^s.Base[j]) == s.Base[s.BaseSize-1] {
				s.Depth[s.BaseSize-1] = s.max(i, j) + 1
				if s.Depth[s.BaseSize-1] > s.MaxDepth {
					s.MaxDepth = s.Depth[s.BaseSize-1]
				}
				var iStr, jStr string
				if i < s.NumInputs {
					iStr = fmt.Sprintf("x%d", i)
				} else {
					iStr = fmt.Sprintf("t%d", i-s.NumInputs+1)
				}
				if j < s.NumInputs {
					jStr = fmt.Sprintf("x%d", j)
				} else {
					jStr = fmt.Sprintf("t%d", j-s.NumInputs+1)
				}
				s.Result = append(s.Result, fmt.Sprintf("t%d = %s + %s * y%d (%d)", s.ProgramSize, iStr, jStr, t, s.Depth[s.BaseSize-1]))
				return true
			}
		}
	}
	return true
}

func (s *SBPAlgorithm) PickNewBaseElement() bool {
	// Array sınır kontrolü
	if s.BaseSize >= MAX_ARRAY_SIZE-1 {
		return false
	}

	minDistance := s.BaseSize * s.NumTargets
	oldNorm := 0
	var bestI, bestJ int
	var theBest uint64
	bestDist := make([]int, s.NumTargets)

	for i := 0; i < s.BaseSize-1; i++ {
		if s.Depth[i]+1 >= s.DepthLimit {
			continue
		}
		for j := i + 1; j < s.BaseSize; j++ {
			if s.Depth[j]+1 >= s.DepthLimit {
				continue
			}
			newBase := s.Base[i] ^ s.Base[j]
			if newBase == 0 || s.isBase(newBase) {
				continue
			}

			depthNewBase := uint64(math.Pow(2, float64(s.max(i, j)+1)))
			thisDist := s.TotalDistance(newBase, depthNewBase)

			if thisDist <= minDistance {
				thisNorm := 0
				for k := 0; k < s.NumTargets; k++ {
					d := s.NDist[k]
					thisNorm = thisNorm + d*d
				}

				if thisDist < minDistance || thisNorm > oldNorm {
					bestI = i
					bestJ = j
					theBest = newBase
					copy(bestDist, s.NDist[:s.NumTargets])
					minDistance = thisDist
					oldNorm = thisNorm
				}
			}
		}
	}

	for i := 0; i < s.NumTargets; i++ {
		s.Dist[i] = bestDist[i]
	}

	s.Base[s.BaseSize] = theBest
	s.Depth[s.BaseSize] = s.max(bestI, bestJ) + 1
	if s.Depth[s.BaseSize] > s.MaxDepth {
		s.MaxDepth = s.Depth[s.BaseSize]
	}
	s.BaseSize++
	s.ProgramSize++

	var iStr, jStr string
	if bestI < s.NumInputs {
		iStr = fmt.Sprintf("x%d", bestI)
	} else {
		iStr = fmt.Sprintf("t%d", bestI-s.NumInputs+1)
	}
	if bestJ < s.NumInputs {
		jStr = fmt.Sprintf("x%d", bestJ)
	} else {
		jStr = fmt.Sprintf("t%d", bestJ-s.NumInputs+1)
	}
	s.Result = append(s.Result, fmt.Sprintf("t%d = %s + %s (%d)", s.ProgramSize, iStr, jStr, s.Depth[s.BaseSize-1]))

	if s.isTarget(theBest) {
		s.TargetsFound++
	}
	return true
}

func (s *SBPAlgorithm) Solve(matrix Matrix) (AlgResult, error) {
	err := s.ReadTargetMatrix(matrix)
	if err != nil {
		return AlgResult{}, err
	}
	
	// Check if depth limit is exceeded
	if s.MaxDist+1 > int(math.Pow(2, float64(s.DepthLimit))) {
		return AlgResult{}, fmt.Errorf("depth limit exceeded")
	}
	
	err = s.InitBase()
	if err != nil {
		return AlgResult{}, err
	}

	iterations := 0
	threshold := 1000 // SBP threshold
	for s.TargetsFound < s.NumTargets && iterations < MAX_ITERATIONS {
		if s.ProgramSize > threshold {
			break
		}
		if !s.EasyMove() {
			if !s.PickNewBaseElement() {
				break // Array sınırına ulaşıldı
			}
		}
		iterations++
	}

	var program []string
	for _, res := range s.Result {
		if res != "" {
			program = append(program, res)
		}
	}

	return AlgResult{
		XorCount: s.ProgramSize,
		Program:  program,
		Depth:    s.MaxDepth,
	}, nil
}

// API Handlers
func boyarHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[BOYAR] İstek başladı - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		log.Printf("[BOYAR] OPTIONS isteği işlendi")
		return
	}

	if r.Method != "POST" {
		log.Printf("[BOYAR] HATA: Geçersiz method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Matrices [][][]string `json:"matrices"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[BOYAR] HATA: JSON decode hatası: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	log.Printf("[BOYAR] %d matris alındı", len(request.Matrices))

	var results []map[string]interface{}
	for i, matrix := range request.Matrices {
		log.Printf("[BOYAR] Matris %d işleniyor (%dx%d)", i+1, len(matrix), len(matrix[0]))
		
		if len(matrix) == 0 {
			log.Printf("[BOYAR] HATA: Matris %d boş", i+1)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        "Empty matrix",
				"xor_count":    0,
				"depth":        0,
				"program":      []string{},
			})
			continue
		}

		boyar := NewBoyarSLP(10)
		err := boyar.ReadTargetMatrix(matrix)
		if err != nil {
			log.Printf("[BOYAR] HATA: Matris %d okuma hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"depth":        0,
				"program":      []string{},
			})
			continue
		}

		err = boyar.InitBase()
		if err != nil {
			log.Printf("[BOYAR] HATA: Matris %d init hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"depth":        0,
				"program":      []string{},
			})
			continue
		}

		result, err := boyar.Solve(matrix)
		if err != nil {
			log.Printf("[BOYAR] HATA: Matris %d solve hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"depth":        0,
				"program":      []string{},
			})
			continue
		}

		log.Printf("[BOYAR] Matris %d başarıyla işlendi - XOR: %d, Derinlik: %d", i+1, result.XorCount, result.Depth)
		results = append(results, map[string]interface{}{
			"matrix_index": i,
			"xor_count":    result.XorCount,
			"depth":        result.Depth,
			"program":      result.Program,
		})
	}

	duration := time.Since(startTime)
	log.Printf("[BOYAR] İstek tamamlandı - Süre: %v, Sonuç sayısı: %d", duration, len(results))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"algorithm": "BoyarSLP",
		"results":   results,
	})
}

func paarHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[PAAR] İstek başladı - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		log.Printf("[PAAR] OPTIONS isteği işlendi")
		return
	}

	if r.Method != "POST" {
		log.Printf("[PAAR] HATA: Geçersiz method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Matrices [][][]string `json:"matrices"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[PAAR] HATA: JSON decode hatası: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	log.Printf("[PAAR] %d matris alındı", len(request.Matrices))

	var results []map[string]interface{}
	for i, matrix := range request.Matrices {
		log.Printf("[PAAR] Matris %d işleniyor (%dx%d)", i+1, len(matrix), len(matrix[0]))
		
		if len(matrix) == 0 {
			log.Printf("[PAAR] HATA: Matris %d boş", i+1)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        "Empty matrix",
				"xor_count":    0,
				"program":      []string{},
			})
			continue
		}

		paar := NewPaarAlgorithm()
		err := paar.ReadTargetMatrix(matrix)
		if err != nil {
			log.Printf("[PAAR] HATA: Matris %d okuma hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"program":      []string{},
			})
			continue
		}

		err = paar.InitBase()
		if err != nil {
			log.Printf("[PAAR] HATA: Matris %d init hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"program":      []string{},
			})
			continue
		}

		result, err := paar.Solve(matrix)
		if err != nil {
			log.Printf("[PAAR] HATA: Matris %d solve hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"program":      []string{},
			})
			continue
		}

		log.Printf("[PAAR] Matris %d başarıyla işlendi - XOR: %d", i+1, result.XorCount)
		results = append(results, map[string]interface{}{
			"matrix_index": i,
			"xor_count":    result.XorCount,
			"program":      result.Program,
		})
	}

	duration := time.Since(startTime)
	log.Printf("[PAAR] İstek tamamlandı - Süre: %v, Sonuç sayısı: %d", duration, len(results))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"algorithm": "PAAR",
		"results":   results,
	})
}

func slpHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	log.Printf("[SLP] İstek başladı - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		log.Printf("[SLP] OPTIONS isteği işlendi")
		return
	}

	if r.Method != "POST" {
		log.Printf("[SLP] HATA: Geçersiz method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Matrices [][][]string `json:"matrices"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[SLP] HATA: JSON decode hatası: %v", err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid JSON: " + err.Error(),
		})
		return
	}

	log.Printf("[SLP] %d matris alındı", len(request.Matrices))

	var results []map[string]interface{}
	for i, matrix := range request.Matrices {
		log.Printf("[SLP] Matris %d işleniyor (%dx%d)", i+1, len(matrix), len(matrix[0]))
		
		if len(matrix) == 0 {
			log.Printf("[SLP] HATA: Matris %d boş", i+1)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        "Empty matrix",
				"xor_count":    0,
				"program":      []string{},
			})
			continue
		}

		slp := NewSLPHeuristic()
		err := slp.ReadTargetMatrix(matrix)
		if err != nil {
			log.Printf("[SLP] HATA: Matris %d okuma hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"program":      []string{},
			})
			continue
		}

		err = slp.InitBase()
		if err != nil {
			log.Printf("[SLP] HATA: Matris %d init hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"program":      []string{},
			})
			continue
		}

		result, err := slp.Solve(matrix)
		if err != nil {
			log.Printf("[SLP] HATA: Matris %d solve hatası: %v", i+1, err)
			results = append(results, map[string]interface{}{
				"matrix_index": i,
				"error":        err.Error(),
				"xor_count":    0,
				"program":      []string{},
			})
			continue
		}

		log.Printf("[SLP] Matris %d başarıyla işlendi - XOR: %d", i+1, result.XorCount)
		results = append(results, map[string]interface{}{
			"matrix_index": i,
			"xor_count":    result.XorCount,
			"program":      result.Program,
		})
	}

	duration := time.Since(startTime)
	log.Printf("[SLP] İstek tamamlandı - Süre: %v, Sonuç sayısı: %d", duration, len(results))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"algorithm": "SLP",
		"results":   results,
	})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("=== XOR Optimizasyon Backend Başlatılıyor ===")
	
	// Load configuration
	configPath := "./config.json"
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatal("Config yüklenemedi:", err)
	}
	
	log.Printf("Config yüklendi: %+v", config)
	
	// Initialize database
	if err := InitDatabase(config); err != nil {
		log.Fatal("Veritabanı başlatılamadı:", err)
	}
	defer db.Close()

	// Auto import data if enabled
	if config.Import.ProcessOnStart {
		log.Println("Başlangıçta otomatik import başlatılıyor...")
		if err := AutoImportData(config); err != nil {
			log.Printf("Otomatik import hatası: %v", err)
		}
	}

	// Create router
	r := mux.NewRouter()

	// CORS middleware
	corsOptions := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}
	if !config.Server.EnableCORS {
		corsOptions.AllowedOrigins = []string{config.Server.Host}
	}
	c := cors.New(corsOptions)

	// Static dosyalar için handler
	staticDir := config.Server.StaticDir
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/index.html")
	})

	// Original API endpoints
	r.HandleFunc("/boyar", boyarHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/paar", paarHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/slp", slpHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/sbp", sbpHandler).Methods("POST", "OPTIONS")

	// New database API endpoints
	r.HandleFunc("/api/matrices", getMatricesHandler).Methods("GET")
	r.HandleFunc("/api/matrices", saveMatrixHandler).Methods("POST")
	r.HandleFunc("/api/matrices/{id:[0-9]+}", getMatrixHandler).Methods("GET")
	r.HandleFunc("/api/matrices/{id:[0-9]+}/inverse", calculateInverseHandler).Methods("POST")
	r.HandleFunc("/api/matrices/process", processAndSaveMatrixHandler).Methods("POST")
	r.HandleFunc("/api/matrices/recalculate", recalculateHandler).Methods("POST")
	r.HandleFunc("/api/matrices/bulk-recalculate", bulkRecalculateHandler).Methods("POST")
	r.HandleFunc("/api/matrices/bulk-inverse", bulkInverseHandler).Methods("POST")
	r.HandleFunc("/api/matrices/missing-algorithms", getMissingAlgorithmsHandler).Methods("GET")
	r.HandleFunc("/api/inverse-pairs", getInversePairsHandler).Methods("GET")

	// Config API endpoints
	r.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
	}).Methods("GET")
	
	r.HandleFunc("/api/config/import", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "POST" {
			if err := AutoImportData(config); err != nil {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Import işlemi başarıyla tamamlandı",
			})
		}
	}).Methods("POST")

	port := config.Server.Port
	log.Printf("Server starting on %s", port)
	log.Printf("Web Interface: http://%s%s", config.Server.Host, port)
	log.Printf("Config dosyası: %s", configPath)
	log.Printf("Data dizini: %s", config.Import.DataDirectory)
	log.Printf("Otomatik import: %v", config.Import.Enabled)
	log.Printf("API Endpoints:")
	log.Printf("  POST /boyar - BoyarSLP algorithm")
	log.Printf("  POST /paar  - Paar algorithm")
	log.Printf("  POST /slp   - SLP Heuristic algorithm")
	log.Printf("  POST /sbp   - SBP algorithm")
	log.Printf("  GET  /api/matrices - Get matrices with pagination")
	log.Printf("  POST /api/matrices - Save matrix")
	log.Printf("  GET  /api/matrices/{id} - Get matrix by ID")
	log.Printf("  POST /api/matrices/process - Process and save matrix")
	log.Printf("  POST /api/matrices/recalculate - Recalculate algorithms")
	log.Printf("  POST /api/matrices/bulk-recalculate - Bulk recalculate algorithms")
	log.Printf("  POST /api/matrices/bulk-inverse - Bulk inverse")
	log.Printf("  GET  /api/matrices/missing-algorithms - Get missing algorithms")
	log.Printf("  GET  /api/inverse-pairs - Get inverse matrix pairs")
	log.Printf("  GET  /api/config - Get current configuration")
	log.Printf("  POST /api/config/import - Trigger manual import")
	log.Printf("=== Backend hazır, istekleri bekleniyor ===")

	handler := c.Handler(r)
	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatal("Server başlatılamadı:", err)
	}
}
