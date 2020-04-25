1. Start simple consensus network
```
./byfn.sh up -o etcdraft
```
If you want to start solo, kafka or sbft consensus, please do the following
```
cd $GOPATH/src/github.com/hyperledger/fabric
mkdir -p scripts/bin
go build -o scripts/bin/configtxgen cmd/configtxgen/main.go
./byfn.sh generate -o solo
./byfn.sh up -o solo
```

2. Start multiple consensus network
```
./byfn.sh up -o etcdraft

docker-compose -f docker-compose-kafka.yaml up -d

docker exec -it cli bash

./scripts/solochannel.sh

./scripts/kafkachannel.sh

./scripts/sbftchannel.sh
```
