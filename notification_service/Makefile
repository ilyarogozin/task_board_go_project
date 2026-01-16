POSTGRES_URL=postgres://postgres:postgres@localhost:5436/notification_db?sslmode=disable

migrate-up:
	migrate -path ./migrations -database "$(POSTGRES_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(POSTGRES_URL)" down 1