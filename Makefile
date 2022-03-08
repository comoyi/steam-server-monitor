build-run:
	make build
	./bin/steam-server-monitor
build:
	go build -o bin/steam-server-monitor main.go
