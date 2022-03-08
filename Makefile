.PHONY: build-run build bundle-font

build-run:
	make build
	./target/linux/steam-server-monitor

build:
	go build -o target/linux/steam-server-monitor main.go
	cp config/config.toml target/linux

build-windows:
	mkdir -p target/windows
	CC=x86_64-w64-mingw32-gcc fyne package -os windows --name target/windows/steam-server-monitor.exe
	cp config/config.toml target/windows

bundle-font:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go <font-file>

deps:
	go get fyne.io/fyne/v2
	go install fyne.io/fyne/v2/cmd/fyne
