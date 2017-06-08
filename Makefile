test:
	go test github.com/exced/simple-blockchain/...

deposit:
	curl -X POST -d '{"from": "me", "to": "you", "currency": "EXC", "amount": 2}' http://localhost:3000/deposit

withdraw:
	curl -X POST -d '{"from": "me", "to": "you", "currency": "EXC", "amount": 2}' http://localhost:3000/withdraw

signin:
	curl -X POST -d '{"key": "admin", "password": "admin"}' http://localhost:3000/signin
