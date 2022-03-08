.PHONY: build-run build bundle-font

build-run:
	make build
	./bin/steam-server-monitor
build:
	go build -o bin/steam-server-monitor main.go
bundle-font:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go <font-file>
