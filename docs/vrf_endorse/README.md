# 背书策略优化
采用VRF算法动态随机选择背书节点，解决背书交易存在的安全风险和性能瓶颈问题。

## 启用VRF优化
VRF优化特性，是在链码层设置的，在部署链码时可以设置 `vrfEnabled` 是否启用VRF优化特性。

在链码部署步骤中，`approveformyorg`、`commit`阶段需要设置`vrfEnabled`,不设置默认为不启用VRF优化。
```
peer lifecycle chaincode approveformyorg --channelID mychannel --name mycc --version 1 --init-required --package-id ${PACKAGE_ID} --sequence ${VERSION} --waitForEvent --vrfEnabled

peer lifecycle chaincode commit -o orderer.example.com:7050 --channelID mychannel --name mycc $PEER_CONN_PARMS --version 1 --sequence ${VERSION} --init-required --vrfEnabled
```
此外如果要`checkCommitReadiness`也需要设置`vrfEnabled`
```
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name mycc $PEER_CONN_PARMS --version 1 --sequence ${VERSION} --output json --init-required --vrfEnabled
```

## 链码调用
在启用VRF优化后，链码调用方式与正常链码调用没有区别，但内部实现上存在区别：

> 启用VRF优化后，由于VRF是随机选择背书节点，所以可能存在没有背书节点被选中（概率极低），如果没有收集到的背书，则需要以majority背书策略的方式重新发送交易（该过程用户无感知）。
