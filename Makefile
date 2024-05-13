watch:
	@wgo run ./cmd/main.go all
build:
	@go build -o ./bin/main ./cmd/main.go


# A function called that will be used to run the build program but will also forward the arguments to the program
run:
	@go run ./cmd/main.go all $(filter-out $@,$(MAKECMDGOALS))

migrateup:
	dbmate up