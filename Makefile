lich_linux: always
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o lich_linux ./cmd/lich

pkg/DEBIAN/control: configs/control.template
	mkdir -p pkg/DEBIAN
	VERSION=$$(cat cmd/lich/VERSION) envsubst < $< > $@

lich.deb: lich_linux configs/lich.service configs/config_example.json pkg/DEBIAN/control
	mkdir -p pkg/opt/lich pkg/etc/systemd/system/ pkg/DEBIAN
	cp lich_linux pkg/opt/lich/lich
	cp configs/lich.service pkg/etc/systemd/system/lich.service
	cp configs/config_example.json pkg/opt/lich/config_example.json
	cp configs/config_example.json pkg/opt/lich/config.json
	dpkg -b pkg lich.deb

always:

.PHONY: always
