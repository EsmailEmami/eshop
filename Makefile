generate-oauth-keys:
	./scripts/generate_keys.sh

swagger-update:
	@swag init --parseDependency

build:
	@GOOS=linux GOARCH=amd64 go build -trimpath github.com/esmailemami/eshop

run:
	@go run . serve

build-run: build
	@./eshop serve

seed:
	@go run . db seed

update-database:
	@go run . migration up