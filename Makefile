generate-oauth-keys:
	./scripts/generate_keys.sh

swagger-update:
	@swag init --parseDependency

run:
	@go run . serve

seed:
	@go run . db seed

update-db:
	@go run . migration up