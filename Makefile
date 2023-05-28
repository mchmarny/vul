VERSION    :=$(shell cat .version)
YAML_FILES :=$(shell find . ! -path "./vendor/*" ! -path "./deployment/*" -type f -regex ".*\.yaml" -print)

all: help

.PHONY: version
version: ## Prints the current version
	@echo $(VERSION)

.PHONY: init
init: ## Initializes all dependancies
	terraform -chdir=./deployment/demo init

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependancies 
	go mod tidy
	go mod vendor

.PHONY: upgrade
upgrade: ## Upgrades all dependancies 
	go get -d -u ./...
	go mod tidy
	go mod vendor

.PHONY: test
test: tidy ## Runs unit tests
	go test -short -count=1 -race -covermode=atomic -coverprofile=cover.out ./...

.PHONY: lint
lint: lint-go lint-yaml ## Lints the entire project 
	@echo "Completed Go and YAML lints"

.PHONY: lint
lint-go: ## Lints the entire project using go 
	golangci-lint -c .golangci.yaml run

.PHONY: lint-yaml
lint-yaml: ## Runs yamllint on all yaml files (brew install yamllint)
	yamllint -c .yamllint $(YAML_FILES)

.PHONY: app
app: tidy ## Builds CLI binary
	mkdir -p ./bin
	CGO_ENABLED=0 go build -trimpath \
	-ldflags="-w -s -X main.version=$(RELEASE_VERSION) \
	-extldflags '-static'" -mod vendor \
	-o bin/server internal/cmd/app/main.go

.PHONY: importer
importer: tidy ## Builds CLI binary
	mkdir -p ./bin
	CGO_ENABLED=0 go build -trimpath \
	-ldflags="-w -s -X main.version=$(RELEASE_VERSION) \
	-extldflags '-static'" -mod vendor \
	-o bin/importer internal/cmd/importer/main.go

.PHONY: vulncheck
vulncheck: ## Checks for soource vulnerabilities
	govulncheck -test ./...

.PHONY: deployment
deployment: ## Applies Terraform deployment
	terraform -chdir=./deployment/demo apply -auto-approve

.PHONY: release
release: test lint tag ## Runs test, lint, and tag before release
	@echo "Triggered image build/publish for: $(VERSION)"
	tool/gh/wait-for-publish-to-finish $(VERSION)
	tool/tf/apply-if-img-exists
	tool/api/e2e

.PHONY: db
db: ## Runs postgres DB as a container
	@echo "Running postgres DB as a container"
	docker run \
		-d \
		--name postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=test \
		-e POSTGRES_DB=vul \
		-p 5432:5432 \
		-v $(PWD)/data:/var/lib/postgresql/data \
		postgres

.PHONY: dbrestore
dbrestore: dbless db ## Restores Cloud SQL DB locally
	@echo "Restoring Cloud SQL DB to locally"
	tool/db/restore

.PHONY: dbconn
dbconn: ## Connect to remote db
	@echo "Connecting to local DB"
	PGPASSWORD=test psql -h localhost -U postgres -d vul

.PHONY: dbless
dbless: ## Stops and remvoes previously run postgres DB container 
	docker stop /postgres
	docker remove /postgres

.PHONY: tag
tag: ## Creates release tag 
	git tag -s -m "version bump to $(VERSION)" $(VERSION)
	git push origin $(VERSION)

.PHONY: tagless
tagless: ## Delete the current release tag 
	git tag -d $(VERSION)
	git push --delete origin $(VERSION)

.PHONY: clean
clean: ## Cleans bin and temp directories
	go clean
	rm -fr ./vendor
	rm -fr ./bin

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
