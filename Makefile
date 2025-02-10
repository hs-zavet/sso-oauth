DB_URL=postgresql://postgres:postgres@localhost:5555/postgres?sslmode=disable
OPENAPI_GENERATOR := java -jar ./openapi-generator-cli.jar
CONFIG_FILE := ./config_local.yaml

generate-models:
	find $(RESOURCES_DIR) -type f ! \( -name "resources_types.go" -o -name "links.go" \) -delete
	swagger-cli bundle $(API_SRC) --outfile $(API_BUNDLED) --type yaml

	$(OPENAPI_GENERATOR) generate \
		-i $(API_BUNDLED) -g go \
		-o $(OUTPUT_DIR) \
		--additional-properties=packageName=resources

	mkdir -p $(RESOURCES_DIR)
	find $(OUTPUT_DIR) -name '*.go' -exec mv {} $(RESOURCES_DIR)/ \;
	find $(RESOURCES_DIR) -type f -name "*_test.go" -delete

start-docs:
	 http-server .

generate-sqlc:
	sqlc generate

migrate-up:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o main main.go
	KV_VIPER_FILE=$(CONFIG_FILE) ./main migrate up

migrate-down:
	migrate -path internal/data/sql/repositories/migrations -database $(DB_URL) -verbose down

run-server:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o main main.go
	KV_VIPER_FILE=$(CONFIG_FILE) ./main run service