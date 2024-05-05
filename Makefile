## runs all tests
test:
	@go test -v ./...

## opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## displays test coverage
coverage:
	@go test -cover ./...

## build command line tool atlas and copies it to my app
build_cli:
	@go build -o ../myApp/atlas ./cmd/cli