# Gossip Archive Meta Info

> 作者：彭洪伟（MaxPeng）
>
> 邮件：pengisgood@gmail.com
>
> 最后更新时间：2020年05月19日

``` mermaid
sequenceDiagram

participant start
participant peer
participant gossip_service
participant archive_service
participant gossip
participant ledger_mgmt
participant kv_ledger
participant hybrid_blockstore
participant blockfile_mgr

start ->> peer: peerInstance.Initialize()
peer ->> peer: p.createChannel()
peer ->> gossip_service: p.GossipService.InitializeChannel()
gossip_service ->> archive_service: archiveservice.New()
archive_service ->> gossip: NewGossipImpl()
gossip ->> archive_service: ArchiveGossip
archive_service ->> archive_service: a.handleMessages()
note over archive_service: go routine
archive_service ->> archive_service: a.gossip.Accept()
note over archive_service: waiting for gossip message
archive_service ->> archive_service: a.handleMessage()
archive_service ->> archive_service: a.getLedger()
archive_service ->> ledger_mgmt: a.ledgerMgr.GetOpenedLedger()
ledger_mgmt ->> archive_service: *ledger.PeerLedger
archive_service ->> kv_ledger: ledger.UpdateArchiveMetaInfo()
kv_ledger ->> hybrid_blockstore: l.blockStore.UpdateArchiveMetaInfo()
hybrid_blockstore ->> blockfile_mgr: store.fileMgr.saveArchiveMetaInfo()
hybrid_blockstore ->> blockfile_mgr: store.fileMgr.updateArchiveMetaInfo()
```

