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
