lich_linux: always
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o lich_linux ./cmd/lich

always:

.PHONY: always
