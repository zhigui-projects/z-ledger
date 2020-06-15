# Ledger Archive Example

> 作者：彭洪伟（MaxPeng）
>
> 邮件：pengisgood@gmail.com
>
> 最后更新时间：2020年06月15日

在本主题中，我们将介绍以下几个方面：

* [配置参数](#配置参数)
* [启动步骤](#启动步骤)
* [归档系统链码](#归档系统链码)



## 配置参数

目前已经实现了和 HDFS、IPFS 的对接，根据自己的需要，在 `core.yaml` 中开启归档模块并配置相应的参数。HDFS 和 IPFS 的配置上稍微有些区别，下面会详细介绍各项参数。

```yaml
ledger:
  archive:
    # 是否开启归档模块 
    enable: false
    # 分布式存储的类型，目前可选项为："hdfs","ipfs"
    type: ipfs
    # 若分布式存储系统为hdfs，则需要提供如下配置
    hdfsConfig:
      # HDFS的namenode地址列表
      nameNodes:
        - 1.2.3.4:9000
      # HDFS的用户，客户端会以该用户身份执行操作
      user: hdfs
      # true: 通过hostname连接datanode，需要在docker-compose中提供DNS配置；false: 通过IP连接datanode
      useDatanodeHostname: false
    # 若分布式存储系统为ipfs，则需要提供如下配置
    ipfsConfig:
      # IPFS的API地址
      url: 1.2.3.4:5001
```

若 `useDatanodeHostname` 配置为 `true`，则需在 Peer 的 docker-compose 文件中提供如下的 DNS 配置（注意：请替换为真实的DNS）：

```yaml
version: '2'

services:
  peer-base:
    image: hyperledger/fabric-peer:$IMAGE_TAG
    dns:
      - 1.2.3.4
      - 5.6.7.8
      - 8.8.8.8
      - 8.8.4.4

```



## 启动步骤

若是选择 HDFS，则需要先启动 HDFS 服务，并确保相应的IP和端口，Peer节点能够正常访问，然后启动Peer节点。启动Peer节点的过程中，会初始化归档服务，并且会尝试连接HDFS服务。如果出现错误，请留意日志中的错误信息，需要注意的是归档模块的初始化失败并不会的导致Peer节点的启动失败。在 HDFS 服务出现故障并恢复后，有可能需要重启Peer节点。

同理，若是选择 IPFS，也需要让IPFS服务在Peer节点之前启动，并确保相应的IP和端口，Peer节点能够正常访问。

迁移后的账本文件所存储的位置默认和Peer节点上的存储目录保持一致。需要特别注意的一点是，由于IPFS是基于内容寻址的，所以如果直接和IPFS里的block或者是object交互会不那么容易，所以目前采用的是IPFS的MFS文件系统，这样就能比较方便地、像操作Peer节点上的目录一样操作IPFS中的文件了。



## 归档系统链码

和其他的系统链码一样，ASCC 系统链码会在 Peer 节点启动的时候自动运行，提供接口来触发归档的执行。类似 QSCC 的实现，用户可以直接通过调用 ASCC 的方法来执行归档。目前支持按照交易的时间戳归档，精确到天。后续可以考虑扩展其他的归档规则，比如区块高度、存储大小等。

等待Peer节点正常启动后，ASCC系统链码也就正常启动了。这时候，就可以通过命令行或者其他方式像调用普通链码一样调用ASCC了。下面以命令行的方式演示，这里需要注意的是，ASCC系统链码需要以组织**管理员**的身份调用，否则会校验权限不通过。

```bash
$ docker exec -it cli sh
$ peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt -C mychannel -n ascc -c '{"Args":["archiveByDate","mychannel","2020-06-15"]}'
```

