.PHONY: build-run build bundle-font deps

build-run:
	make build
	./target/linux/steam-server-monitor

build:
	go build -o target/linux/steam-server-monitor main.go
	cp config/config.toml target/linux/

package-linux:
	make build
	cd target/linux && tar zcvf steam-server-monitor.tar.gz config.toml steam-server-monitor && cd -

package-linux-installer:
	fyne package -os linux --release
	mkdir -p target/linux
	mv steam-server-monitor.tar.xz target/linux/

package-windows:
	mkdir -p target/windows
	CC=x86_64-w64-mingw32-gcc fyne package -os windows --release --appID com.comoyi.steam-server-monitor --name target/windows/steam-server-monitor.exe
	cp config/config.toml target/windows/
	cd target/windows && zip steam-server-monitor.zip config.toml steam-server-monitor.exe && cd -

clean:
	rm -rf target

bundle-font:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go <font-file>
	#fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go ~/.local/share/fonts/HarmonyOS_Sans_SC_Regular.ttf

bundle-font-build:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go /usr/local/share/fonts/HarmonyOS_Sans_SC_Regular.ttf

deps:
	go get fyne.io/fyne/v2
	go install fyne.io/fyne/v2/cmd/fyne@latest
