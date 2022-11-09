# speedster

#### 介绍

harmony-rpc 仓库主要代码实现功能：

* 基于 RPC的 加速
* Redis 缓存 harmony 的相关接口链上数据，提升接口请求速度
* 基于权重、一致性hash的后端endpoint的负载均衡策略
* 基于fasthttp的连接池

#### 软件架构

* 软件架构，参考doc文件夹内开发相关文档。

—— speeder

｜ —— balancing（负载均衡策略）

｜ —— cache（内存数据库缓存，目前实现了redis， 请求到redis的实现）

｜ —— common（harmony的基础数据结构，但适配rpc时修改过部分结构的数据类型）

｜ —— doc（相关设计文档）

｜ —— etc（配置文件）

｜ —— internal（Restful API server 实现，包括handler实现）

｜ —— logs（日志数据文档，api.yaml配置文件里配置log为file输出）

｜ —— rpc（rpc client实现，请求到endpoint的实现）

｜ —— util（一些辅助函数）

｜ —— interface.go（harmony的rpc接口，暂时未用到，可用于重构）

｜ —— main.go（主入口函数）

#### 安装教程

编译运行：

* `go mod init`
* `go mod tidy`
* `go build main.go`
* `./main -f ./etc/api.yaml`

#### 使用说明

1. rpc接口适配harmony官方的接口请求参数
2. 主要配置文件在 `./etc/api.yaml`文件，启动代码时，需要此配置文件，启动参数参考**编译运行**
3. 代码中有 `todo`、`fixme`的地方是需要后续优化代码 / 结构的地方，可以优化代码逻辑（才疏学浅）
4. 只实现了 Restful 的rpc功能，未实现websocket，可参考代码自行实现
5. 关于连接池，如果考虑拓展性，可自行实现连接池模型，可能会优于fasthttp。也可以考虑用fasthttp的server，或`gnet`替换harmony节点源码的net包，性能应该会有提升
6. Redis里实现了大部分的接口，但有几个地方需要优化

   1. redis的key的timeout时间，可以依据当前epoch、当前block的时间来做key的timeout时间，这样优于目前固定的 2s timeout（harmony共识时间）
   2. redis里的key可以考虑重新优化（数据库优化范畴）
7. Redis主要实现缓存的接口如下


   | 编号 | 支持的Json-RPC接口                               | 是否Redis索引 | 备注               |
   | :----: | -------------------------------------------------- | --------------- | -------------------- |
   |  1  | hmyv2_getBalance                                 | ✅            | 2s超时             |
   |  2  | hmyv2_getBalanceByBlockNumber                    | ✅            |                    |
   |  3  | hmyv2_getStakingTransactionsCount                | ✅            | 2s超时             |
   |  4  | hmyv2_getStakingTransactionsHistory              | ✅            | 代码逻辑需要重新改 |
   |  5  | hmyv2_getTransactionsCount                       | ✅            | 2s超时             |
   |  6  | hmyv2_getTransactionsHistory                     |               |                    |
   |  7  | hmyv2_getBlocks                                  | ✅            |                    |
   |  8  | hmyv2_getBlockByNumber                           | ✅            |                    |
   |  9  | hmyv2_getBlockByHash                             | ✅            |                    |
   |  10  | hmyv2_getBlockSigners                            | ✅            |                    |
   |  11  | hmyv2_getBlockSignersKeys                        | ✅            |                    |
   |  12  | hmyv2_getBlockTransactionCountByNumber           | ✅            |                    |
   |  13  | hmyv2_getBlockTransactionCountByHash             | ✅            |                    |
   |  14  | hmyv2_getHeaderByNumber                          | ✅            |                    |
   |  15  | hmyv2_getLatestChainHeaders                      | ✅            | 2s过期             |
   |  16  | hmyv2_latestHeader                               | ✅            | 2s过期             |
   |  17  | hmyv2_blockNumber                                | ✅            | 2s过期             |
   |  18  | hmyv2_getCirculatingSupply                       |               |                    |
   |  19  | hmyv2_getEpoch                                   |               |                    |
   |  20  | hmyv2_getLastCrossLinks                          |               |                    |
   |  21  | hmyv2_getLeader                                  |               |                    |
   |  22  | hmyv2_gasPrice                                   |               |                    |
   |  23  | hmyv2_getShardingStructure                       |               |                    |
   |  24  | hmyv2_getTotalSupply                             |               |                    |
   |  25  | hmyv2_getValidators                              | ✅            | 2s过期             |
   |  26  | hmyv2_getValidatorKeys                           | ✅            | 2s过期（重新设计） |
   |  27  | [WIP] hmyv2_getCurrentBadBlocks                  |               |                    |
   |  28  | hmyv2_getNodeMetadata                            | ✅            | 2s过期             |
   |  29  | hmyv2_protocolVersion                            |               |                    |
   |  30  | net_peerCount                                    |               |                    |
   |  31  | hmyv2_call                                       |               |                    |
   |  32  | hmyv2_estimateGas                                |               |                    |
   |  33  | hmyv2_getCode                                    | ✅            |                    |
   |  34  | hmyv2_getStorageAt                               |               |                    |
   |  35  | hmyv2_getDelegationsByDelegator                  |               |                    |
   |  36  | hmyv2_getDelegationsByDelegatorByBlockNumber     |               |                    |
   |  37  | hmyv2_getDelegationsByValidator                  |               |                    |
   |  38  | hmyv2_getAllValidatorAddresses                   |               |                    |
   |  39  | hmyv2_getAllValidatorInformation                 |               |                    |
   |  40  | hmyv2_getAllValidatorInformationByBlockNumber    |               |                    |
   |  41  | hmyv2_getElectedValidatorAddresses               |               |                    |
   |  42  | hmyv2_getValidatorInformation                    |               |                    |
   |  43  | hmyv2_getCurrentUtilityMetrics                   |               |                    |
   |  44  | hmyv2_getMedianRawStakeSnapshot                  |               |                    |
   |  45  | hmyv2_getStakingNetworkInfo                      |               |                    |
   |  46  | hmyv2_getSuperCommittees                         |               |                    |
   |  47  | hmyv2_getCXReceiptByHash                         |               |                    |
   |  48  | hmyv2_getPendingCXReceipts                       |               |                    |
   |  49  | hmyv2_resendCx                                   |               |                    |
   |  50  | hmyv2_getPoolStats                               |               |                    |
   |  51  | hmyv2_pendingStakingTransactions                 |               |                    |
   |  52  | hmyv2_pendingTransactions                        |               |                    |
   |  53  | hmyv2_getCurrentStakingErrorSink                 |               |                    |
   |  54  | hmyv2_getStakingTransactionByBlockNumberAndIndex |               |                    |
   |  55  | hmyv2_getStakingTransactionByBlockHashAndIndex   |               |                    |
   |  56  | hmyv2_getStakingTransactionByHash                | ✅            |                    |
   |  57  | hmyv2_sendRawStakingTransaction                  |               |                    |
   |  58  | hmyv2_getCurrentTransactionErrorSink             |               |                    |
   |  59  | hmyv2_getTransactionByBlockHashAndIndex          | ✅            |                    |
   |  60  | hmyv2_getTransactionByBlockNumberAndIndex        | ✅            |                    |
   |  61  | hmyv2_getTransactionByHash                       | ✅            |                    |
   |  62  | hmyv2_getTransactionReceipt                      | ✅            |                    |
   |  63  | hmyv2_sendRawTransaction                         |               |                    |

---

#### 参与贡献

1. Fork 本仓库
2. 新建 Feat_xxx 分支
3. 提交代码
4. 新建 Pull Request
