
X_APP_VERSION := $(shell cat VERSION)

.PHONY: build-run
build-run:
	make build
	./target/linux/steam-server-monitor

.PHONY: build
build:
	go build -o target/linux/steam-server-monitor main.go
	cp config/config.toml target/linux/

.PHONY: package-linux
package-linux:
	make build
	cd target/linux && tar zcvf steam-server-monitor-$(X_APP_VERSION)-linux.tar.gz config.toml steam-server-monitor && cd -

.PHONY: package-linux-installer
package-linux-installer:
	fyne package -os linux --release
	mkdir -p target/linux
	mv steam-server-monitor.tar.xz target/linux/steam-server-monitor-$(X_APP_VERSION)-linux-installer.tar.xz

.PHONY: package-windows
package-windows:
	mkdir -p target/windows
	CC=x86_64-w64-mingw32-gcc fyne package -os windows --release --appID com.comoyi.steam-server-monitor --name target/windows/steam-server-monitor.exe
	cp config/config.toml target/windows/
	cd target/windows && zip steam-server-monitor-$(X_APP_VERSION)-windows.zip config.toml steam-server-monitor.exe && cd -

.PHONY: package-android
package-android:
	ANDROID_HOME=~/Android/Sdk ANDROID_NDK_HOME=~/Android/Sdk/ndk/23.2.8568313 fyne package --release --target android --appID com.comoyi.steamservermonitor --name steam_server_monitor --appVersion $(X_APP_VERSION)
	mkdir -p target/android
	mv steam_server_monitor.aab target/android/steam_server_monitor_$(X_APP_VERSION).aab
	bundletool build-apks --overwrite --mode=universal --bundle target/android/steam_server_monitor_$(X_APP_VERSION).aab --output target/android/steam_server_monitor_$(X_APP_VERSION).apks --ks ~/.jks/bundle.jks --ks-pass file:.jks-pass --ks-key-alias key1 --key-pass file:.jks-pass
	unzip -o -d target/android target/android/steam_server_monitor_$(X_APP_VERSION).apks
	mv target/android/universal.apk target/android/steam_server_monitor_$(X_APP_VERSION).apk

.PHONY: clean
clean:
	rm -rf target

.PHONY: bundle-font
bundle-font:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go <font-file>
	#fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go ~/.local/share/fonts/HarmonyOS_Sans_SC_Regular.ttf

.PHONY: bundle-font-build
bundle-font-build:
	fyne bundle --package fonts --prefix Resource --name DefaultFont -o fonts/default_font.go /usr/local/share/fonts/HarmonyOS_Sans_SC_Regular.ttf

.PHONY: deps
deps:
	go get fyne.io/fyne/v2
	go install fyne.io/fyne/v2/cmd/fyne@latest
