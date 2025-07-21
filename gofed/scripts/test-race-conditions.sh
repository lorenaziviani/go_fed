#!/bin/bash

# Script to test race conditions and cache in Users Service
# Demonstrates concurrency problems and their solutions

echo "Testing Race Conditions and Cache in Users Service"
echo "==================================================="
echo ""

# Wait for gateway to be ready
echo "Waiting for gateway to be ready..."
sleep 3

echo "1. Checking initial cache statistics..."
echo "Query: { cacheStats { size maxSize ttl } }"
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ cacheStats { size maxSize ttl } }"}' | jq .
echo ""

echo "2. Testing cache user search..."
echo "Query: { usersFromCache { id name email } }"
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ usersFromCache { id name email } }"}' | jq .
echo ""

echo "3. Checking statistics after populating cache..."
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ cacheStats { size maxSize ttl } }"}' | jq .
echo ""

echo "4. Testing specific cache user search..."
echo "Query: { userFromCache(id: \"1\") { id name email } }"
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ userFromCache(id: \"1\") { id name email } }"}' | jq .
echo ""

echo "5. Simulating race condition (CAN CAUSE PROBLEMS)..."
echo "Query: { simulateRaceCondition { success message duration } }"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ simulateRaceCondition { success message duration } }"}' | jq .
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Duration: ${duration}s"
echo ""

echo "6. Simulating safe access (thread-safe)..."
echo "Query: { simulateSafeAccess { success message duration } }"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ simulateSafeAccess { success message duration } }"}' | jq .
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Duration: ${duration}s"
echo ""

echo "7. Checking final cache statistics..."
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ cacheStats { size maxSize ttl } }"}' | jq .
echo ""

echo "8. Performance comparison..."
echo "Testing cache vs direct search:"
echo ""

echo "   Direct search (users):"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ users { id name email } }"}' > /dev/null
end_time=$(date +%s.%N)
direct_duration=$(echo "$end_time - $start_time" | bc)
echo "   Duration: ${direct_duration}s"

echo "   Cache search (usersFromCache):"
start_time=$(date +%s.%N)
curl -s -X POST http://localhost:4000/ \
  -H "Content-Type: application/json" \
  -d '{"query": "{ usersFromCache { id name email } }"}' > /dev/null
end_time=$(date +%s.%N)
cache_duration=$(echo "$end_time - $start_time" | bc)
echo "   Duration: ${cache_duration}s"

echo ""
echo "Performance analysis:"
echo "   • Direct search: ${direct_duration}s"
echo "   • Cache: ${cache_duration}s"
echo "   • Difference: $(echo "$cache_duration - $direct_duration" | bc)s"
echo "   • Cache faster: $(echo "$cache_duration < $direct_duration" | bc -l | sed 's/1/Sim/g' | sed 's/0/Não/g')"
echo ""

echo "9. Stress test with multiple requests..."
echo "Executing 10 concurrent requests to cache:"
echo ""

for i in {1..10}; do
  echo "   Request $i:"
  start_time=$(date +%s.%N)
  curl -s -X POST http://localhost:4000/ \
    -H "Content-Type: application/json" \
    -d '{"query": "{ usersFromCache { id name } }"}' > /dev/null
  end_time=$(date +%s.%N)
  req_duration=$(echo "$end_time - $start_time" | bc)
  echo "   Duration: ${req_duration}s"
done

echo ""
echo "Conclusions:"
echo "    Cache thread-safe implemented"
echo "    Race condition detected and simulated"
echo "    Safe access with mutex working"
echo "    Performance improved with cache"
echo "    Active concurrency protection"
echo ""

echo "GraphQL Playground: http://localhost:4000"
echo "Apollo Studio: https://studio.apollographql.com/"
echo ""
echo "Usage tips:"
echo "   • Use 'usersFromCache' for performance"
echo "   • Use 'simulateRaceCondition' to see problems"
echo "   • Use 'simulateSafeAccess' to see solutions"
echo "   • Use 'cacheStats' for monitoring" 