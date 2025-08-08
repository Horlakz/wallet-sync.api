dev:
	@nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go

build:
	@go build -o wallet_sync

run:
	@./wallet_sync

build_run: build run
