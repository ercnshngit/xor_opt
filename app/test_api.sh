#!/bin/bash

echo "=== XOR Optimization API Test ==="
echo ""

# Test BoyarSLP
echo "Testing BoyarSLP Algorithm..."
curl -s -X POST http://localhost:8080/boyar \
  -H "Content-Type: application/json" \
  -d @test_data.json | jq '.algorithm, .results[0].xor_count, .results[0].depth'

echo ""

# Test Paar
echo "Testing Paar Algorithm..."
curl -s -X POST http://localhost:8080/paar \
  -H "Content-Type: application/json" \
  -d @test_data.json | jq '.algorithm, .results[0].xor_count'

echo ""

# Test SLP Heuristic
echo "Testing SLP Heuristic Algorithm..."
curl -s -X POST http://localhost:8080/slp \
  -H "Content-Type: application/json" \
  -d @test_data.json | jq '.algorithm, .results[0].xor_count'

echo ""
echo "=== Test Complete ===" 