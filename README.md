# Blockchain - Cryptocurrency

## Miner Protocol
* Connect to other peer and receive minimal info to start serving and mining
* Connect to network and listen and serve

## Client
* withdraw
* deposit
* Look at blockchain

## Example

### Cli

* Depose 100 cryptocurrency to account stored in file "private.pem".
Deposit from bank is currently fake since transactions are not verified.

* Withdraw 20 cryptocurrency to Thomas, using account stored in file "private.pem"

```bash
~cli go run main.go -i "./private.pem" deposit HSBC 100
~cli go run main.go -i "./private.pem" withdraw Thomas 20
```

### Miner

* The first miner is a special case since it does not connect to an existing peer.
* Peer listen http on port 3001, use blockchain file blockchain1.bc to load, and store blockchain, connects to peer
at port 3000 and use "private1.pem" file as its account. This account is needed to reward successful miner.

```bash
~miner go run main.go
~miner go run main.go -p 3001 -b "./blockchain1.bc" 3000
~miner go run main.go -p 3002 -b "./blockchain2.bc" 3001
```
