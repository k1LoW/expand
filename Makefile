default: test

ci: depsdev test sec

test:
	go test ./... -coverprofile=coverage.out -covermode=count

sec:
	gosec ./...

lint:
	golangci-lint run ./...

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@v0.10.2
	go install github.com/Songmu/gocredits/cmd/gocredits@v0.2.0
	go install github.com/securego/gosec/v2/cmd/gosec@v2.11.0

prerelease:
	git pull origin main --tag
	go mod tidy
	ghch -w -N ${VER}
	gocredits . > CREDITS
	git add CHANGELOG.md CREDITS go.mod go.sum
	git commit -m'Bump up version number'
	git tag ${VER}

.PHONY: default test
