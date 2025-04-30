DB_HOST ?= localhost
DB_PORT ?= 3306
DB_USER ?= root
DB_PASSWORD ?= password
DB_NAME ?= mckinney_go_notes_db

MIGRATE=migrate -path ./migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)"

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-steps:
	$(MIGRATE) down $(steps)

migrate-force:
	$(MIGRATE) force $(version)

migrate-version:
	$(MIGRATE) version