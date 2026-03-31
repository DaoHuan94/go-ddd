include .env
export

MIGRATE=go run -mod=mod -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

migrate-up:
	$(MIGRATE) -path migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path migrations -database "$(DB_URL)" down -all

migrate-create:
	$(MIGRATE) create -ext sql -dir migrations -seq $(name)