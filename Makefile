.PHONY: run-backend run-app test

run-backend:
	cd backend && go run ./cmd/api

run-app:
	cd app && flutter run

test:
	cd backend && go test ./...
	cd app && flutter test
