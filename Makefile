.PHONY: run
run:
	$(GOENV) go run cmd/main.go $(RUN_ARGS)

.PHONY: migrate
migrate:
	$(GOENV) go run cmd/migrate/main.go $(RUN_ARGS)

.PHONY: migrate-minio
migrate-minio:
	$(GOENV) go run cmd/minio/main.go $(RUN_ARGS)

.PHONY: gendoc
gendoc:
	 swag init -g cmd/main.go 