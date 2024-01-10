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
	~/go/bin/swag init -g cmd/main.go --parseDependency --parseInternal

.PHONY: genitems
genitems:
	$(GOENV) python3 cmd/random_items/random_items.py $(RUN_ARGS)

.PHONY: insgenit
insgenit:
	$(GOENV) go run cmd/random_items/main.go $(RUN_ARGS)