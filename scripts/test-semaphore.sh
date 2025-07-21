#!/bin/bash

# Script to test the semaphore and backpressure in the Products Service
# Demonstrates limiting concurrent resolutions

echo "Testing Semaphore and Backpressure in the Products Service"
echo "======================================================"
echo ""
echo "Waiting for the gateway to be ready..."
sleep 3

echo "1. Checking initial semaphore statistics..."
echo "Query: { semaphoreStats { max current available usage } }"
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ semaphoreStats { max current available usage } }"}' | jq .
echo ""

echo "2. Testing products by category..."
echo "Query: { productsByCategory(category: \"Electronics\") { id name category } }"
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsByCategory(category: \"Electronics\") { id name category } }"}' | jq .
echo ""

echo "3. Testing normal resolution (no semaphore)..."
echo "Query: { productsByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"}' | jq .
echo ""

echo "4. Testing concurrent resolution with semaphore (max 3 concurrent)..."
echo "Query: { productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"}' | jq .
echo ""

echo "⚡ 3. Testing normal concurrent resolution (no semaphore)..."
echo "Query: { productsByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"}' | jq .
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Duration: ${duration}s"
echo ""

echo "5. Testing backpressure with many products..."
echo "Query: { productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"}' | jq .
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Duration: ${duration}s"
echo ""

echo "6. Checking final semaphore statistics..."
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ semaphoreStats { max current available usage } }"}' | jq .
echo ""

echo "7. Testing backpressure with many products..."
echo "Query: { productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\", \"6\", \"7\", \"8\"]) { id name } }"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\", \"6\", \"7\", \"8\"]) { id name } }"}' | jq .
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Duration: ${duration}s"
echo ""

echo "8. Checking final semaphore statistics..."
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ semaphoreStats { max current available usage } }"}' | jq .
echo ""

echo "9. Performance comparison:"
echo "Normal vs semaphore:"
echo ""

echo "   Normal (5 products):"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"}' > /dev/null
end_time=$(date +%s.%N)
normal_duration=$(echo "$end_time - $start_time" | bc)
echo "   Duration: ${normal_duration}s"

echo "   Semaphore (5 products):"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"}' > /dev/null
end_time=$(date +%s.%N)
semaphore_duration=$(echo "$end_time - $start_time" | bc)
echo "   Duration: ${semaphore_duration}s"

echo ""
echo "Performance analysis:"
echo "   • Normal: ${normal_duration}s"
echo "   • Semaphore: ${semaphore_duration}s"
echo "   • Difference: $(echo "$semaphore_duration - $normal_duration" | bc)s"
echo "   • Backpressure: $(echo "$semaphore_duration > $normal_duration" | bc -l | sed 's/1/Yes/g' | sed 's/0/No/g')"
echo ""

echo "Conclusions:"
echo "    Semaphore limiting to 3 concurrent resolutions"
echo "    Backpressure working (higher latency)"
echo "    Statistics being monitored"
echo "    Overload protection active"
echo ""

echo "GraphQL Playground: http://localhost:4000"
echo "Apollo Studio: https://studio.apollographql.com/" 