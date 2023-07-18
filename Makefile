test:
	cp env.sample.json env.json
	rm -rf env.json

test-with-report:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
