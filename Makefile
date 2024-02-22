GO?=go
DESTDIR?=/usr/local

default: bin/lich

all: bin/lich_linux_amd64 bin/lich_linux_arm64 bin/lich_darwin_amd64 bin/lich_darwin_arm64 bin/lich.deb

cmd/lich/version.txt: always
	echo $$(cat VERSION)-$$(git rev-parse --short HEAD) > $@

# Native binary for the current platform.
bin/lich: cmd/lich/version.txt always
	CGO_ENABLED=0 $(GO) build -mod=vendor -o $@ ./cmd/lich

bin/lich_linux_amd64: cmd/lich/version.txt always
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -mod=vendor -o $@ ./cmd/lich

bin/lich_linux_arm64: cmd/lich/version.txt always
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -mod=vendor -o $@ ./cmd/lich

bin/lich_darwin_amd64: cmd/lich/version.txt always
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -mod=vendor -o $@ ./cmd/lich

bin/lich_darwin_arm64: cmd/lich/version.txt always
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -mod=vendor -o $@ ./cmd/lich

pkg/DEBIAN/control: config/control.template cmd/lich/version.txt
	mkdir -p pkg/DEBIAN
	VERSION=$$(cat cmd/lich/version.txt) envsubst < $< > $@

bin/lich.deb: bin/lich_linux_amd64 config/lich.service config/config_example.json pkg/DEBIAN/control
	mkdir -p pkg/opt/lich pkg/etc/systemd/system/ pkg/DEBIAN
	cp bin/lich_linux_amd64 pkg/opt/lich/lich
	cp config/lich.service pkg/etc/systemd/system/lich.service
	cp config/config_example.json pkg/opt/lich/config_example.json
	dpkg -b pkg bin/lich.deb

run: bin/lich
	bin/lich -config config/config.json

install: bin/lich
	install -Dm755 bin/lich $(DESTDIR)/bin/lich

always:

clean:
	rm -f \
		bin/lich_linux_amd64 \
		bin/lich_linux_arm64 \
		bin/lich_darwin_amd64 \
		bin/lich_darwin_arm64 \
		bin/lich \
		bin/lich.deb

.PHONY: all always clean default run install
