test:
	go test github.com/exced/simple-blockchain/...

deposit:
	curl -X POST -d '{"From": "me", "To": "you", "Currency": "EXC", "Amount": 2}' http://localhost:3000/deposit

withdraw:
	curl -X POST -d '{"From": "me", "To": "you", "Currency": "EXC", "Amount": 2}' http://localhost:3000/withdraw
