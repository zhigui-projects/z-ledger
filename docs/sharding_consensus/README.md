# 分片共识
分片共识技术支持不同通道使用不同的共识方式，目前可选的共识有：solo、kafka、etcdraft、sbft(自研)、hotstuff(自研)。

## 环境搭建
测试脚本位置: fabric/scripts/first-network

共识节点配置：fabric/scripts/first-network/configtx.yaml

### 启动网络
```
./byfn.sh up -o etcdraft
```
脚本中生成的配置(genesis.block、channel.tx)是etcdraft共识的，如果您要启动其他共识的网络，请按如下操作：
```
cd $GOPATH/src/github.com/hyperledger/fabric

mkdir -p scripts/bin

go build -o scripts/bin/configtxgen cmd/configtxgen/main.go

cd scripts/first-network

// solo、kafka、etcdraft、sbft、hotstuff
./byfn.sh generate -o solo
./byfn.sh up -o solo
```

### 创建分片网络
在启动网络后，创建新的通道，采用不同的共识
```
// 如果要采用kafka共识，请先启动kafka服务
docker-compose -f docker-compose-kafka.yaml up -d

docker exec -it cli bash

./scripts/solochannel.sh

./scripts/kafkachannel.sh

./scripts/sbftchannel.sh

./scripts/hschannel.sh
```