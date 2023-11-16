generate-oauth-keys:
	./scripts/generate_keys.sh

swagger-update:
	@swag init --parseDependency

build:
	@GOOS=linux GOARCH=amd64 go build -trimpath github.com/esmailemami/eshop

run: build
	@./eshop serve

seed:
	@go run . db seed

update-database:
	@go run . migration up

deploy:
	liara deploy --disks uploads:/uploads --app eshop-bak --port 8080