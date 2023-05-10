VERSION    :=$(shell cat .version)
YAML_FILES :=$(shell find . ! -path "./vendor/*" ! -path "./deployment/*" -type f -regex ".*y*ml" -print)
REG_URI    :=us-west1-docker.pkg.dev/s3cme1/vul
DB_BUCKET  :=vuln-db-dumps
NOW        ?=$(shell date +%s)

all: help

.PHONY: version
version: ## Prints the current version
	@echo $(VERSION)

.PHONY: init
init: ## Initializes all dependancies
	terraform -chdir=./deployment init

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
	go test -count=1 -race -covermode=atomic -coverprofile=cover.out ./...

.PHONY: lint
lint: lint-go lint-yaml ## Lints the entire project 
	@echo "Completed Go and YAML lints"

.PHONY: lint
lint-go: ## Lints the entire project using go 
	golangci-lint -c .golangci.yaml run

.PHONY: lint-yaml
lint-yaml: ## Runs yamllint on all yaml files (brew install yamllint)
	yamllint -c .yamllint $(YAML_FILES)

.PHONY: build
build: tidy ## Builds CLI binary
	mkdir -p ./bin
	CGO_ENABLED=0 go build -trimpath \
	-ldflags="-w -s -X main.version=$(RELEASE_VERSION) \
	-extldflags '-static'" -mod vendor \
	-o bin/server internal/cmd/main.go

.PHONY: image
image: ## Builds container image
	KO_DOCKER_REPO=$(REG_URI)/vul \
    GOFLAGS="-ldflags=-X=main.version=$(VERSION)" \
    ko build internal/cmd/main.go --bare --tags $(VERSION)

.PHONY: vulncheck
vulncheck: ## Checks for soource vulnerabilities
	govulncheck -test ./...

.PHONY: deploy
deploy: ## Applies Terraform deployment
	terraform -chdir=./deployment apply -auto-approve

.PHONY: release
release: test lint tag ## Runs test, lint, and tag before release
	@echo "Triggered image build/publish for: $(RELEASE_VERSION)"
	tools/gh/wait-for-publish-to-finish
	tools/tf/apply-if-img-exists "$(REG_URI)/vul" "$(RELEASE_VERSION)"

.PHONY: server
server: ## Runs uncompiled app 
	go run internal/cmd/main.go

.PHONY: db
db: ## Runs postgres DB as a container
	docker run \
		-d \
		--name postgres \
		-e POSTGRES_USER=vimp \
		-e POSTGRES_PASSWORD=test \
		-p 5432:5432 \
		-v $(PWD)/data:/var/lib/postgresql/data \
		postgres

.PHONY: dbrestre
dbrestre: ## Restores Cloud SQL DB locally
	gcloud sql export sql db gs://$(DB_BUCKET)/$(NOW).gz -d vimp
	gsutil cp gs://$(DB_BUCKET)/$(NOW).gz tools/db/dump/
	gzip -d tools/db/dump/$(NOW).gz
	PGPASSWORD=test psql -h localhost -U vimp -d vimp < tools/db/sql/pre-restore.sql
	PGPASSWORD=test psql -h localhost -U vimp -d vimp < tools/db/dump/$(NOW)
	PGPASSWORD=test psql -h localhost -U vimp -d vimp < tools/db/sql/post-restore.sql
	gsutil rm gs://$(DB_BUCKET)/$(NOW).gz
	# rm tools/db/dump/$(NOW)

.PHONY: dbconn
dbconn: ## Connect to remote db
	PGPASSWORD=test psql -h localhost -U vimp -d vimp

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
