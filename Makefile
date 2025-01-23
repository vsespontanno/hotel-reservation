<<<<<<< HEAD
build:
	@go build -o bin/api.exe


run: build
	@./bin/api.exe

test:
=======
build:
	@go build -o bin/api.exe


run: build
	@./bin/api.exe

test:
>>>>>>> 7322db8 (Add initial project structure and user API)
	@go test -v ./...