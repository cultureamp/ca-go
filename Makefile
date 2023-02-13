.PHONY: mod
mod:
	go mod download

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: test-ci
test-ci:
	mkdir artifacts
	go test -race ./... -covermode=atomic -coverprofile=artifacts/count.out
	go tool cover -func=artifacts/count.out | tee artifacts/coverage.out

# ensures that `go mod tidy` has been run after any dependency changes
.PHONY: ensure-deps
ensure-deps: mod
	@go mod tidy
	@git diff --exit-code
