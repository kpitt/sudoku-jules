# Download external Go dependencies during setup.
go mod download
go mod tidy

# Install `golangci-lint` for linting commands
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.10.1
