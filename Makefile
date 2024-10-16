.PHONY: help
help:
	echo 'Hi'

.PHONY: build-run
build-run:
	make build
	./target/linux/steam-server-monitor

.PHONY: build
build:
	go build -o target/linux/steam-server-monitor main.go
	cp config/config.toml target/linux/
