watch:
	@wgo run ./cmd/main.go all

build:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/quran-macos ./cmd/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o  ./bin/quran-linux ./cmd/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o  ./bin/quran-windows.exe ./cmd/main.go

run:
	./bin/quran-linux all

migrateup:
	dbmate up