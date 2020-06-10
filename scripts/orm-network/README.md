1.编译tools
```
cd $GOPATH/src/github.com/hyperledger/fabric
mkdir -p scripts/bin
go build -o scripts/bin/configtxgen cmd/configtxgen/main.go
go build -o scripts/bin/cryptogen cmd/cryptogen/main.go
```

2.默认是sqlite3数据库
```
cd scripts/orm-network
./byfn.sh up 
```

3.切换到mysql数据库
```
cd scripts/orm-network
./byfn.sh up -s mysqldb
```
