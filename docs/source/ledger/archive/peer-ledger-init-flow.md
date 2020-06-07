# Peer 账本初始化流程

> 作者：彭洪伟（MaxPeng）
>
> 邮件：pengisgood@gmail.com
>
> 最后更新时间：2020年05月19日

``` mermaid
sequenceDiagram

participant start
participant peer
participant ledger_mgmt
participant kv_ledger_provider
participant kv_ledger
participant store
participant hybrid_blockstore_provider
participant hybrid_blockstore
participant blockfile_mgr

start ->> ledger_mgmt: ledgermgmt.NewLedgerMgr()
ledger_mgmt ->> kv_ledger_provider: kvledger.NewProvider()
kv_ledger_provider ->> kv_ledger_provider: initLedgerStorageProvider()
kv_ledger_provider ->> store: ledgerstorage.NewProvider()
store ->> hybrid_blockstore_provider:  hybridblkstorage.NewProvider()
hybrid_blockstore_provider ->> store: blkstorage.BlockStoreProvider
store ->> kv_ledger_provider: *ledgerstorage.Provider
kv_ledger_provider ->> ledger_mgmt:  *kvledger.Provider
ledger_mgmt ->> start: *ledgermgmt.LedgerMgr
start ->> start: initGossipService()
start ->> peer: peerInstance.Initialize()
peer ->> ledger_mgmt: LedgerMgr.OpenLedger()
ledger_mgmt ->> kv_ledger_provider: ledgerProvider.Open()
kv_ledger_provider ->> kv_ledger_provider: openInternal()
kv_ledger_provider ->> store: ledgerStoreProvider.Open()
store ->> hybrid_blockstore_provider: blkStoreProvider.OpenBlockStore()
hybrid_blockstore_provider ->> hybrid_blockstore:newHybridBlockStore()
hybrid_blockstore ->> blockfile_mgr: newBlockfileMgr()
blockfile_mgr ->> hybrid_blockstore: *hybridBlockfileMgr
hybrid_blockstore ->> hybrid_blockstore_provider: *hybridBlockStore
hybrid_blockstore_provider ->> store: blkstorage.BlockStore
store ->> kv_ledger_provider: *Store
kv_ledger_provider ->> kv_ledger: newKVLedger() 
kv_ledger ->> kv_ledger_provider: *kvLedger
kv_ledger_provider ->> ledger_mgmt: ledger.PeerLedger
ledger_mgmt ->> peer: ledger.PeerLedger
peer ->> peer: createChannel()
note over peer: initialize gossip and archive serive
peer ->> peer: initChannel()
```