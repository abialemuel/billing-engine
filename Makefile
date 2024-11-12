export PKGNAME=billing-engine

help:
	@echo "See Makefile"

generate-mocks:
	@rm -rf mocks
	@mkdir -p mocks/config
	@mkdir -p mocks/src/adapter/mongo
	@mockgen -source=config/config.go -destination=mocks/config/config.go -package=config Config
	@mockgen -source=pkg/user/business/contract/repository.go -destination=mocks/pkg/user/business/contract/repository.go -package=contract Repository

run:
	@go run main.go

run-docker:
	@docker-compose up --build

### Test
test:
	@go test -cover -coverprofile=coverage.cov $$(go list ./... | grep -v /mock/)
	@go tool cover -func coverage.cov
	@go tool cover -html=coverage.cov -o coverage.html
	@rm -f coverage.cov



