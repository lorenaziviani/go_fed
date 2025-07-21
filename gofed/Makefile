run-users:
	cd services/users && go run main.go

run-products:
	cd services/products && go run main.go

run-gateway:
	cd gateway && npm install && npm start

run-all:
	@echo "Starting all services..."
	@echo "Users service on port 8081"
	@echo "Products service on port 8082"
	@echo "Gateway on port 4000"
	@echo ""
	@echo "Make sure to run services in separate terminals:"
	@echo "Terminal 1: make run-users"
	@echo "Terminal 2: make run-products"
	@echo "Terminal 3: make run-gateway"
	@echo ""
	@echo "Or use Docker Compose: docker-compose up"

docker-up:
	docker-compose up

docker-down:
	docker-compose down

docker-build:
	docker-compose build

test-queries:
	@echo "Running federation tests..."
	./scripts/test-queries.sh

test-semaphore:
	@echo "Testing semaphore and backpressure..."
	./scripts/test-semaphore.sh

test-race-conditions:
	@echo "Testing race conditions and cache..."
	./scripts/test-race-conditions.sh

test-federation:
	@echo "Testing federation with owner queries..."
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ products { id name owner { id name email } } }"}' | jq .

test-complex:
	@echo "Testing complex federation query..."
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ users { id name } products { id name price category owner { id name } } }"}' | jq .

test-concurrent:
	@echo "Testing concurrent user resolution..."
	@echo "Query: usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"])"
	@start_time=$$(date +%s.%N); \
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name email } }"}' | jq .; \
	end_time=$$(date +%s.%N); \
	duration=$$(echo "$$end_time - $$start_time" | bc); \
	echo "Duration: $${duration}s"

test-performance:
	@echo "Testing performance with 8 concurrent users..."
	@echo "Query: usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\", \"6\", \"7\", \"8\"])"
	@start_time=$$(date +%s.%N); \
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ usersByIds(ids: [\"1\", \"2\", \"3\", \"4\", \"5\", \"6\", \"7\", \"8\"]) { id name email } }"}' | jq .; \
	end_time=$$(date +%s.%N); \
	duration=$$(echo "$$end_time - $$start_time" | bc); \
	echo "Duration: $${duration}s"

test-semaphore-performance:
	@echo "Testing semaphore performance..."
	@echo "Query: productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"])"
	@start_time=$$(date +%s.%N); \
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ productsWithSemaphore(ids: [\"1\", \"2\", \"3\", \"4\", \"5\"]) { id name } }"}' | jq .; \
	end_time=$$(date +%s.%N); \
	duration=$$(echo "$$end_time - $$start_time" | bc); \
	echo "Duration: $${duration}s"

test-semaphore-stats:
	@echo "Checking semaphore statistics..."
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ semaphoreStats { max current available usage } }"}' | jq .

test-cache:
	@echo "Testing user cache..."
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ usersFromCache { id name email } }"}' | jq .

test-metrics:
	@echo "Testing Prometheus metrics and request tracing..."
	./scripts/test-metrics.sh

show-users-metrics:
	@echo "Users Service Metrics:"
	curl -s http://localhost:8081/metrics | grep -E "(graphql_requests_total|graphql_active_requests|cache_hits_total|cache_misses_total)" | head -10

show-products-metrics:
	@echo "Products Service Metrics:"
	curl -s http://localhost:8082/metrics | grep -E "(graphql_requests_total|graphql_active_requests|semaphore_current|semaphore_max)" | head -10

show-all-metrics:
	@echo "=== Users Service Metrics ==="
	@make show-users-metrics
	@echo ""
	@echo "=== Products Service Metrics ==="
	@make show-products-metrics

test-tracing:
	@echo "Testing request tracing with TraceID..."
	@trace_id=$$(uuidgen); \
	echo "Generated TraceID: $$trace_id"; \
	curl -s -H "X-Trace-ID: $$trace_id" \
		-H "Content-Type: application/json" \
		-d '{"query": "{ __typename }"}' \
		http://localhost:8081/query | jq .; \
	echo "Check logs for TraceID: $$trace_id"

load-test-metrics:
	@echo "Generating load to test metrics..."
	@for i in {1..20}; do \
		trace_id=$$(uuidgen); \
		curl -s -H "X-Trace-ID: $$trace_id" \
			-H "Content-Type: application/json" \
			-d '{"query": "{ __typename }"}' \
			http://localhost:8081/query > /dev/null & \
		curl -s -H "X-Trace-ID: $$trace_id" \
			-H "Content-Type: application/json" \
			-d '{"query": "{ __typename }"}' \
			http://localhost:8082/query > /dev/null & \
	done; \
	wait; \
	echo "Load test completed. Check metrics endpoints."

test-cache-stats:
	@echo "Checking cache statistics..."
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ cacheStats { size maxSize ttl } }"}' | jq .

test-race-simulation:
	@echo "Simulating race condition (may cause issues)..."
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ simulateRaceCondition { success message duration } }"}' | jq .

test-safe-access:
	@echo "Testing safe access simulation..."
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ simulateSafeAccess { success message duration } }"}' | jq .

test-race:
	@echo "Running race detection tests..."
	cd services/users && go test -race ./...

test-benchmark:
	@echo "Running performance benchmarks..."
	cd services/users && go test -bench . -benchmem ./graph/

test-benchmark-detail:
	@echo "Running detailed benchmarks..."
	./scripts/benchmark.sh

test-cache-benchmark:
	@echo "Running cache performance benchmarks..."
	cd services/users && go test -bench=BenchmarkCache -benchmem ./graph/

test-all:
	@echo "Running all tests..."
	@echo "1. Race detection..."
	@$(MAKE) test-race
	@echo ""
	@echo "2. Performance benchmarks..."
	@$(MAKE) test-benchmark
	@echo ""
	@echo "3. Cache benchmarks..."
	@$(MAKE) test-cache-benchmark
	@echo ""
	@echo "4. Federation tests..."
	@$(MAKE) test-queries
	@echo ""
	@echo "5. Semaphore tests..."
	@$(MAKE) test-semaphore
	@echo ""
	@echo "6. Race condition tests..."
	@$(MAKE) test-race-conditions

health-check:
	@echo "Checking service health..."
	@echo "Users service:"
	curl -s http://localhost:8081/healthz | jq .
	@echo "Products service:"
	curl -s http://localhost:8082/healthz | jq .
	@echo "Gateway introspection:"
	curl -s -X POST http://localhost:4000/ \
		-H "Content-Type: application/json" \
		-d '{"query": "{ __schema { types { name } } }"}' | jq '.data.__schema.types | length'

clean:
	@echo "Cleaning up..."
	docker-compose down -v
	docker system prune -f
	@echo "Cleanup complete!"
