test-with-report:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down