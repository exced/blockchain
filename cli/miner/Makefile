build:
	go build -o app
	docker build -t . .
	rm -f app

run-genesis:
	go run main.go --g true
