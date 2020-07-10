# IPFS Private Network Setup

> 作者：彭洪伟（MaxPeng）
>
> 邮件：pengisgood@gmail.com
>
> 最后更新时间：2020年07月10日

IPFS 作为通过内容寻址并以点对点传输的分布式文件系统，本身并不提供在多个存储节点间内容复制的功能。如果需要支持在多个存储节点之间像HDFS一样进行内容复制，可以参考 [IPFS Cluster](https://cluster.ipfs.io/) 的方案，IPFS Cluster的相关内容不在本文的讨论范围之内。

本文试验的环境为Ubuntu 16.04.6 LTS。下面的步骤描述了如何一步一步搭建两个节点的私有IPFS网络。

## 端口

| 端口 | 作用                           |
| :--: | ------------------------------ |
| 4001 | 用于P2P发现其他的节点          |
| 5001 | API server                     |
| 8080 | Gateway server，浏览器中会用到 |

## 安装Golang环境

Golang的环境需要在每个节点上都安装，步骤如下：

1. 更新包和依赖

```shell
sudo apt-get update
sudo apt-get -y upgrade
```

2. 安装golang

Golang的最新安装包可以从这里下载：https://golang.org/dl/

```shell
wget https://golang.org/dl/go1.14.4.linux-amd64.tar.gz
sudo tar -xvf go1.14.4.linux-amd64.tar.gz
sudo mv go /usr/local
```

3. 设置环境变量

在`.bashrc`中设置`GOPATH`， `GOROOT`， `PATH`：

```shell
mkdir $HOME/gopath
sudo vim $HOME/.bashrc

# 添加以下几行到.bashrc中	
export GOROOT=/usr/local/go
export GOPATH=$HOME/gopath
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

4. 更新环境变量并检查golang的版本

```shell
source ~/.bashrc
go version
```

## 安装IPFS

IPFS需要在每个节点上都安装，步骤如下：

1. 安装IPFS

IPFS的最新安装包可以从这里下载：https://dist.ipfs.io/#go-ipfs

```shell
wget https://dist.ipfs.io/go-ipfs/v0.6.0/go-ipfs_v0.6.0_linux-amd64.tar.gz
tar xvfz go-ipfs_v0.6.0_linux-amd64.tar.gz
sudo mv go-ipfs/ipfs /usr/local/bin/ipfs
ipfs init
ipfs version
```

2. 创建私有网络

需要提前创建`swarm.key`，然后拷贝到每个节点上。

```shell
echo -e "/key/swarm/psk/1.0.0/\n/base16/\n`tr -dc 'a-f0-9' < /dev/urandom | head -c64`" > ~/.ipfs/swarm.key
```

3. 设置Bootstrap节点

其中Peer ID可以通过`ipfs id`查看。

```shell
# 删除所有的bootstrap节点
ipfs bootstrap rm --all

# 给每个节点添加一个指定的bootstrap节点，假设为10.0.1.19
ipfs bootstrap add /ip4/10.0.1.19/tcp/4001/ipfs/QmRtS4uyvkXt83gSZ1YCWkDbYKRLwgzANKBvtxw8fqGbzE
```

修改`.bashrc`中的环境变量，设置为私有网络：

```shell
vim ~/.bashrc

export LIBP2P_FORCE_PNET=1

source ~/.bashrc
```

4. 修改IPFS服务绑定的IP

默认都是绑定`127.0.0.1`，这样其他的节点就访问不到了，这里我们都改为该节点的内网IP。

```shell
vim ~/.ipfs/config
```

将配置中的IP修改为本机的内网IP（可以通过`ifconfig`查看，这里假设为：10.0.1.19）：

```json
"Addresses": {
    "API": "/ip4/10.0.1.19/tcp/5001",
    "Announce": [],
    "Gateway": "/ip4/10.0.1.19/tcp/8080",
    "NoAnnounce": [],
    "Swarm": [
      "/ip4/0.0.0.0/tcp/4001",
      "/ip6/::/tcp/4001",
      "/ip4/0.0.0.0/udp/4001/quic",
      "/ip6/::/udp/4001/quic"
    ]
  }
```

5. 以后台服务启动IPFS节点

创建如下`systemd` service

``` shell
sudo vim /etc/systemd/system/ipfs.service
```

内容如下：

```shell
[Unit]
 Description=IPFS Daemon
 After=syslog.target network.target remote-fs.target nss-lookup.target
 [Service]
 Type=simple
 ExecStart=/usr/local/bin/ipfs daemon --enable-namesys-pubsub
 User=root
 [Install]
 WantedBy=multi-user.target
```

启动后台服务

```shell
sudo systemctl daemon-reload
sudo systemctl enable ipfs
sudo systemctl start ipfs
sudo systemctl status ipfs
```

如果最后看到如下的日志，怎说明服务已经正常启动。

```
● ipfs.service - IPFS Daemon
   Loaded: loaded (/etc/systemd/system/ipfs.service; enabled; vendor preset: enabled)
   Active: active (running) since Fri 2020-07-10 15:27:42 CST; 6s ago
 Main PID: 19911 (ipfs)
   CGroup: /system.slice/ipfs.service
           └─19911 /usr/local/bin/ipfs daemon --enable-namesys-pubsub

Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: Swarm listening on /ip4/127.0.0.1/tcp/4001
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: Swarm listening on /ip6/::1/tcp/4001
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: Swarm listening on /p2p-circuit
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: Swarm announcing /ip4/10.0.1.19/tcp/4001
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: Swarm announcing /ip4/127.0.0.1/tcp/4001
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: Swarm announcing /ip6/::1/tcp/4001
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: API server listening on /ip4/10.0.1.19/tcp/5001
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: WebUI: http://10.0.1.19:5001/webui
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: Gateway (readonly) server listening on /ip4/10.0.1.19/tcp/8080
Jul 10 15:27:42 iZj6cisuup16tg9zo25rtyZ ipfs[19911]: Daemon is ready
```

