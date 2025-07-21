#!/bin/bash

# Script to test Prometheus metrics and request tracing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configurations
USERS_SERVICE_URL="http://localhost:8081"
PRODUCTS_SERVICE_URL="http://localhost:8082"
GATEWAY_URL="http://localhost:4000"

# Function to log in color
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
}

info() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] INFO: $1${NC}"
}

# Function to check if a service is running
check_service() {
    local service_name=$1
    local url=$2
    
    if curl -s "$url/healthz" > /dev/null 2>&1; then
        log "$service_name is running on $url"
        return 0
    else
        error "$service_name is not running on $url"
        return 1
    fi
}

# Function to test metrics
test_metrics() {
    local service_name=$1
    local url=$2
    
    info "Testing metrics of $service_name..."
    
    # Check if the /metrics endpoint is available
    if ! curl -s "$url/metrics" > /dev/null 2>&1; then
        error "Endpoint /metrics is not available on $url"
        return 1
    fi
    
    # Check if the metrics contain basic data
    local metrics_output=$(curl -s "$url/metrics")
    
    # Check specific metrics
    if echo "$metrics_output" | grep -q "graphql_requests_total"; then
        log "✓ graphql_requests_total metric found"
    else
        warn "graphql_requests_total metric not found"
    fi
    
    if echo "$metrics_output" | grep -q "graphql_request_duration_seconds"; then
        log "✓ graphql_request_duration_seconds metric found"
    else
        warn "graphql_request_duration_seconds metric not found"
    fi
    
    if echo "$metrics_output" | grep -q "graphql_active_requests"; then
        log "✓ graphql_active_requests metric found"
    else
        warn "graphql_active_requests metric not found"
    fi
    
    # Service-specific metrics
    if [ "$service_name" = "users" ]; then
        if echo "$metrics_output" | grep -q "cache_hits_total"; then
            log "✓ cache_hits_total metric found"
        else
            warn "cache_hits_total metric not found"
        fi
        
        if echo "$metrics_output" | grep -q "cache_misses_total"; then
            log "✓ cache_misses_total metric found"
        else
            warn "cache_misses_total metric not found"
        fi
    fi
    
    if [ "$service_name" = "products" ]; then
        if echo "$metrics_output" | grep -q "semaphore_current"; then
            log "✓ semaphore_current metric found"
        else
            warn "semaphore_current metric not found"
        fi
        
        if echo "$metrics_output" | grep -q "semaphore_max"; then
            log "✓ semaphore_max metric found"
        else
            warn "semaphore_max metric not found"
        fi
    fi
    
    log "Metrics of $service_name tested successfully"
}

# Function to test request tracing
test_tracing() {
    local service_name=$1
    local url=$2
    
    info "Testing request tracing of $service_name..."
    
    # Generate a unique TraceID
    local trace_id=$(uuidgen)
    
    # Make a request with TraceID
    local response=$(curl -s -w "%{http_code}" -H "X-Trace-ID: $trace_id" \
        -H "Content-Type: application/json" \
        -d '{"query": "{ __typename }"}' \
        "$url/query")
    
    local http_code="${response: -3}"
    local response_body="${response%???}"
    
    if [ "$http_code" = "200" ]; then
        log "✓ GraphQL request successful"
    else
        error "GraphQL request failed with code $http_code"
        return 1
    fi
    
    # Check if the TraceID was returned in the header
    local response_trace_id=$(curl -s -I -H "X-Trace-ID: $trace_id" \
        -H "Content-Type: application/json" \
        -d '{"query": "{ __typename }"}' \
        "$url/query" | grep -i "X-Trace-ID" | cut -d' ' -f2 | tr -d '\r')
    
    if [ -n "$response_trace_id" ]; then
        log "✓ TraceID returned in header: $response_trace_id"
    else
        warn "TraceID was not returned in the response header"
    fi
    
    log "Request tracing of $service_name tested successfully"
}

# Function to generate load and test metrics
test_load_metrics() {
    local service_name=$1
    local url=$2
    local iterations=$3
    
    info "Generating load on $service_name ($iterations requests)..."
    
    # Make multiple requests to generate metrics
    for i in $(seq 1 $iterations); do
        local trace_id=$(uuidgen)
        
        curl -s -H "X-Trace-ID: $trace_id" \
            -H "Content-Type: application/json" \
            -d '{"query": "{ __typename }"}' \
            "$url/query" > /dev/null 2>&1 &
    done
    
    # Wait for all requests to finish
    wait
    
    log "Load generated on $service_name"
}

# Function to show current metrics
show_current_metrics() {
    local service_name=$1
    local url=$2
    
    info "Current metrics of $service_name:"
    
    local metrics_output=$(curl -s "$url/metrics")
    
    # Extract specific values
    local requests_total=$(echo "$metrics_output" | grep "graphql_requests_total" | grep -v "#" | head -1 | sed 's/.*} //')
    local active_requests=$(echo "$metrics_output" | grep "graphql_active_requests" | grep -v "#" | head -1 | sed 's/.*} //')
    
    echo "  - Total requests: $requests_total"
    echo "  - Active requests: $active_requests"
    
    if [ "$service_name" = "users" ]; then
        local cache_hits=$(echo "$metrics_output" | grep "cache_hits_total" | grep -v "#" | head -1 | sed 's/.*} //')
        local cache_misses=$(echo "$metrics_output" | grep "cache_misses_total" | grep -v "#" | head -1 | sed 's/.*} //')
        echo "  - Cache hits: $cache_hits"
        echo "  - Cache misses: $cache_misses"
    fi
    
    if [ "$service_name" = "products" ]; then
        local semaphore_current=$(echo "$metrics_output" | grep "semaphore_current" | grep -v "#" | head -1 | sed 's/.*} //')
        local semaphore_max=$(echo "$metrics_output" | grep "semaphore_max" | grep -v "#" | head -1 | sed 's/.*} //')
        echo "  - Semaphore current: $semaphore_current"
        echo "  - Semaphore max: $semaphore_max"
    fi
}

# Main function
main() {
    log "Starting metrics and tracing tests..."
    
    # Check if the services are running
    log "Checking if the services are running..."
    
    if ! check_service "Users Service" "$USERS_SERVICE_URL"; then
        error "Users Service is not available. Execute 'make run-users' first."
        exit 1
    fi
    
    if ! check_service "Products Service" "$PRODUCTS_SERVICE_URL"; then
        error "Products Service is not available. Execute 'make run-products' first."
        exit 1
    fi
    
    if ! check_service "Gateway" "$GATEWAY_URL"; then
        error "Gateway is not available. Execute 'make run-gateway' first."
        exit 1
    fi
    
    log "All services are running!"
    
    # Test initial metrics
    log "=== Testing initial metrics ==="
    test_metrics "users" "$USERS_SERVICE_URL"
    test_metrics "products" "$PRODUCTS_SERVICE_URL"
    
    # Test request tracing
    log "=== Testing request tracing ==="
    test_tracing "users" "$USERS_SERVICE_URL"
    test_tracing "products" "$PRODUCTS_SERVICE_URL"
    
    # Show metrics before load
    log "=== Metrics before load ==="
    show_current_metrics "users" "$USERS_SERVICE_URL"
    show_current_metrics "products" "$PRODUCTS_SERVICE_URL"
    
    # Generate load
    log "=== Generating load to test metrics ==="
    test_load_metrics "users" "$USERS_SERVICE_URL" 10
    test_load_metrics "products" "$PRODUCTS_SERVICE_URL" 10
    
    # Wait a bit for the metrics to update
    sleep 2
    
    # Show metrics after load
    log "=== Metrics after load ==="
    show_current_metrics "users" "$USERS_SERVICE_URL"
    show_current_metrics "products" "$PRODUCTS_SERVICE_URL"
    
    # Test specific queries for cache and semaphore
    log "=== Testing specific queries ==="
    
    # Test cache of users service
    info "Testing cache of Users Service..."
    for i in {1..3}; do
        curl -s -H "Content-Type: application/json" \
            -d '{"query": "{ users { id name email } }"}' \
            "$USERS_SERVICE_URL/query" > /dev/null 2>&1
    done
    
    # Test semaphore of products service
    info "Testing semaphore of Products Service..."
    for i in {1..3}; do
        curl -s -H "Content-Type: application/json" \
            -d '{"query": "{ productsWithSemaphore(ids: [\"1\", \"2\"]) { id name } }"}' \
            "$PRODUCTS_SERVICE_URL/query" > /dev/null 2>&1
    done
    
    # Wait 
    sleep 2
    
    # Show final metrics
    log "=== Final metrics ==="
    show_current_metrics "users" "$USERS_SERVICE_URL"
    show_current_metrics "products" "$PRODUCTS_SERVICE_URL"
    
    log "Metrics and tracing tests completed successfully!"
    log ""
    log "To view complete metrics:"
    log "  curl $USERS_SERVICE_URL/metrics"
    log "  curl $PRODUCTS_SERVICE_URL/metrics"
    log ""
    log "To test tracing manually:"
    log "  curl -H 'X-Trace-ID: $(uuidgen)' -H 'Content-Type: application/json' -d '{\"query\": \"{ __typename }\"}' $USERS_SERVICE_URL/query"
}

# Execute main function
main "$@" 