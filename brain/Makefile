.PHONY: connect
connect:
	ssh brain@192.168.86.100

.PHONY: build
build:
	go build -o bin/server

.PHONY: run
run:
	go run ./*.go

.PHONY: test
test:
	go test ./...

.PHONY: deploy
deploy:
	echo "Building for pi architecture..."
	GOOS=linux GOARCH=arm GOARM=7 go build -o bin/server
	echo "Stopping the server..."
	ssh -t brain@192.168.86.100 'sudo systemctl stop server'
	echo "Building for pi architecture..."
	scp bin/server brain@192.168.86.100:/home/brain/server
	echo "Starting the server again..."
	ssh -t brain@192.168.86.100 'sudo systemctl start server'
