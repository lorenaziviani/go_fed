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
