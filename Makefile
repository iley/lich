default: lich_linux

all: lich_linux lich.deb

lich_linux: always
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o lich_linux ./cmd/lich

pkg/DEBIAN/control: config/control.template
	mkdir -p pkg/DEBIAN
	VERSION=$$(cat cmd/lich/VERSION) envsubst < $< > $@

lich.deb: lich_linux config/lich.service config/config_example.json pkg/DEBIAN/control
	mkdir -p pkg/opt/lich pkg/etc/systemd/system/ pkg/DEBIAN
	cp lich_linux pkg/opt/lich/lich
	cp config/lich.service pkg/etc/systemd/system/lich.service
	cp config/config_example.json pkg/opt/lich/config_example.json
	dpkg -b pkg lich.deb

always:

clean:
	rm -f lich_linux lich.deb

.PHONY: all always clean default
