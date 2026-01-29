.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./bin/ghupdater

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/ghupdater-linux-amd64 .

.PHONY: docker-build
docker-build: build-linux
	docker build -t ghupdater:latest .

# Run ghupdater in a Debian container against indiependente/backupper (private).
# Requires GITHUB_TOKEN. Removes the container after run (--rm).
# Output is extracted into ./out .
# Usage: GITHUB_TOKEN=ghp_xxx make run-docker-backupper
.PHONY: run-docker-backupper
run-docker-backupper: docker-build
	@if [ -z "$$GITHUB_TOKEN" ]; then echo "Error: GITHUB_TOKEN must be set (e.g. GITHUB_TOKEN=ghp_xxx make run-docker-backupper)"; exit 1; fi
	mkdir -p out
	docker run --rm \
		-e GITHUB_TOKEN \
		-v "$$(pwd)/out:/out" \
		ghupdater:latest \
		-owner indiependente -repo backupper \
		-archive tar.gz -os linux -arch amd64 \
		-extract /out

.PHONY: coverage
coverage:
	gopherbadger -md="README.md" -png=false

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -cover ./...
