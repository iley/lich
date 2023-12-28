default: lich

all: lich_linux_amd64 lich_linux_arm64 lich_darwin_amd64 lich_darwin_amd64 lich.deb

cmd/lich/version.txt: always
	echo $$(cat VERSION)-$$(git rev-parse --short HEAD) > $@

# Native binary for the current platform.
lich: cmd/lich/version.txt always
	CGO_ENABLED=0 go build -mod=vendor -o lich ./cmd/lich

lich_linux_amd64: cmd/lich/version.txt always
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o $@ ./cmd/lich

lich_linux_arm64: cmd/lich/version.txt always
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod=vendor -o $@ ./cmd/lich

lich_darwin_amd64: cmd/lich/version.txt always
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod=vendor -o $@ ./cmd/lich

lich_darwin_arm64: cmd/lich/version.txt always
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -mod=vendor -o $@ ./cmd/lich

pkg/DEBIAN/control: config/control.template cmd/lich/version.txt
	mkdir -p pkg/DEBIAN
	VERSION=$$(cat cmd/lich/version.txt) envsubst < $< > $@

lich.deb: lich_linux_amd64 config/lich.service config/config_example.json pkg/DEBIAN/control
	mkdir -p pkg/opt/lich pkg/etc/systemd/system/ pkg/DEBIAN
	cp lich_linux_amd64 pkg/opt/lich/lich
	cp config/lich.service pkg/etc/systemd/system/lich.service
	cp config/config_example.json pkg/opt/lich/config_example.json
	dpkg -b pkg lich.deb

run: lich
	./lich -config config/config.json

always:

clean:
	rm -f lich_linux_amd64 lich lich.deb

.PHONY: all always clean default run
