#!/bin/bash

# Test queries for the Gofed Federation Gateway
# Executes example queries via curl

GATEWAY_URL="http://localhost:4000"

echo "Testing Gofed Federation Gateway"
echo "=================================="
echo ""

# Function to execute query
execute_query() {
    local query_name="$1"
    local query="$2"
    
    echo "$query_name"
    echo "Query: $query"
    echo "Result:"
    curl -s -X POST "$GATEWAY_URL/" \
        -H "Content-Type: application/json" \
        -d "{\"query\": \"$query\"}" | jq .
    echo ""
    echo "----------------------------------------"
    echo ""
}

# Wait for gateway to be ready
echo "⏳ Waiting for gateway to be ready..."
sleep 3

# Test 1: Basic query of users
execute_query "1. Basic query of users" "{ users { id name email } }"

# Test 2: Basic query of products
execute_query "2. Basic query of products" "{ products { id name description price category } }"

# Test 3: Federated query - products with owner
execute_query "3. Federated query - products with owner" "{ products { id name owner { id name email } } }"

# Test 4: Federated query - users and products together
execute_query "4. Federated query - users and products together" "{ users { id name } products { id name owner { id name } } }"

# Test 5: Query by ID
execute_query "5. Query by ID" "{ user(id: \"1\") { id name email } }"

# Test 6: Query by ID with owner
execute_query "6. Query by ID with owner" "{ product(id: \"1\") { id name description price owner { id name email } } }"

# Test 7: Complex federated query
execute_query "7. Query complexa federada" "{ users { id name } products { id name price category owner { id name } } }"

# Test 8: Concurrent query - multiple users (WaitGroup + Channels)
echo "8. Concurrent query - multiple users (WaitGroup + Channels)"
echo "Query: { usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name email } }"
echo "Result:"
start_time=$(date +%s.%N)
curl -s -X POST "$GATEWAY_URL/" \
    -H "Content-Type: application/json" \
    -d '{"query": "{ usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name email } }"}' | jq .
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Duration: ${duration}s"
echo ""
echo "----------------------------------------"
echo ""

# Test 9: Concurrent query - performance test (8 users)
echo "9. Concurrent query - performance test (8 users)"
echo "Query: { usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\", \"6\", \"7\", \"8\"]) { id name email } }"
echo "Result:"
start_time=$(date +%s.%N)
curl -s -X POST "$GATEWAY_URL/" \
    -H "Content-Type: application/json" \
    -d '{"query": "{ usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\", \"6\", \"7\", \"8\"]) { id name email } }"}' | jq .
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Duration: ${duration}s"
echo ""
echo "----------------------------------------"
echo ""

# Test 10: Federated query with concurrent owners
echo "10. Federated query with concurrent owners"
echo "Query: { products { id name owner { id name email } } usersByIds(ids: [\"1\", \"2\", \"3\"]) { id name email } }"
echo "Result:"
start_time=$(date +%s.%N)
curl -s -X POST "$GATEWAY_URL/" \
    -H "Content-Type: application/json" \
    -d '{"query": "{ products { id name owner { id name email } } usersByIds(ids: [\"1\", \"2\", \"3\"]) { id name email } }"}' | jq .
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Duration: ${duration}s"
echo ""
echo "----------------------------------------"
echo ""

echo "All tests completed!"
echo ""
echo "GraphQL Playground available at: $GATEWAY_URL"
echo "Apollo Studio: https://studio.apollographql.com/"
echo ""
echo "Concurrent resolutions implemented:"
echo "   • WaitGroup for synchronization"
echo "   • Channels for communication"
echo "   • Context for cancellation"
echo "   • Timeout of 5 seconds"
echo "   • Simulated latency of 100ms" 