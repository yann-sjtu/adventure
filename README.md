## Adventure

### 1. 简介

adventure为OKChain的一款强大的交易批量发送工具。使用者可以通过构建自己的配置文件来决定发送交易的类型组合和并发线程数。

adventure依赖于OKChain Go SDK，在编译前请使用go module下载相关依赖。

### 2. 编译与配置
编译：

```shell
$ make adventure
```

配置` config.toml`，一般不做修改。

config.toml 是基础配置文件，存放测试节点：

```yaml
hosts = ["http://127.0.0.1:26657",, ……]
```

配置`tx.json` (样例)

```json
[
    {
        "mnemonic_path":"template/mnemonic/normal_5", //助记词地址1
        "transactions":[
            {
                "type":"issue", //tx类型1
                "args":{
                    "concurrent_num":5, //单次发送tx个数
                    "sleep_time":86400 //发送间隔(s)
                }
            }
        ]
    },
    {
        "mnemonic_path":"template/mnemonic/normal_100", //助记词地址2
        "transactions":[
            {
                "type":"mint", //tx类型1
                "args":{
                    "concurrent_num":1,
                    "sleep_time":20
                }
            },
            {
                "type":"edit", //tx类型2
                "args":{
                    "concurrent_num":1,
                    "sleep_time":15
                }
            }
        ]
    }
]
```

type字段可选填：

```toml
// distribution
WithdrawRewards = "withdraw-rewards"
SetWithdrawAddr = "set-withdraw-addr"

//token
Issue                  = "issue"
Burn                   = "burn"
Mint                   = "mint"
MultiSend              = "multi-send"
Edit                   = "edit"

//dex
List                 = "list"
Deposit              = "deposit"
Withdraw             = "withdraw"

//staking
DelegateVoteUnbond = "delegate_vote_unbond"
Proxy              = "proxy"
```

* 注1：adventure工程下 `template/*.json`已有许多样例，

* 注2：adventure工程下 `template/mnemonic/*`存放助记词，`template/addr/*`存放对应的地址(用于查询、转账等)

### 3. 使用说明

初始化账户：

```shell
# 往 tx.json 中需要的账户转钱
adventure account send --init_amount 1000tokt -m template/mnemonic/normal_100
adventure account send --init_amount 1000tokt -p template/address/normal_100
```

启动：

```shell
nohup adventure start -p template/tx.json > ~/tx.log 2>&1 &
```

### 启动全部交易类型测试
#### 测试账户转账
```shell script
adventure account send -p template/address/captain -a 1000000tokt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/normal_5 -a 1000000tokt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/normal_1000_1 -a 1000000tokt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/normal_1000_2 -a 1000000tokt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/normal_100 -a 1000000tokt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
```
#### 启动
```shell script
nohup adventure start -p template/test_cases/proxy1.json   >> ~/proxy-staking.log 2>&1 &
nohup adventure start -p template/test_cases/staking.json   >> ~/staking.log 2>&1 &
nohup adventure start -p template/test_cases/token-dex-distr.json   >> ~/token-dex-distr.log 2>&1 &
nohup adventure start -p template/test_cases/multi-send.json   >> ~/multi-send.log 2>&1 &
nohup adventure start -p template/test_cases/issue-list.json   >> ~/issue-list.log 2>&1 &
nohup adventure swap loop -p template/mnemonic/normal_100 -g 25 >> ~/swap.log 2>&1 &
```

```shell script
adventure account send -r "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -p "./template/address/farm_test/pooler_total" -a 7000tokt

adventure farm pooler issue-token -p "./template/mnemonic/farm_test/pooler"
adventure farm allocate-tokens 10000usdk -p "./template/address/farm_test/pooler_total"
adventure farm pooler create-pair -p "./template/mnemonic/farm_test/pooler"
adventure farm pooler add-liquidity -p "./template/mnemonic/farm_test/pooler"
adventure farm locker allocate-tokens-to-lockers-from-all-poolers -p "./template/mnemonic/farm_test/pooler" -l "./template/address/farm_test/locker_total"

adventure account send -r "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -p "./template/address/farm_test/locker_total" -a 1000tokt

nohup adventure farm pooler strategy-pooler -p "./template/mnemonic/farm_test/pooler" > ../adventure_log/pooler.log 2>&1 &
nohup adventure farm locker strategy-lock-unlock -p "./template/mnemonic/farm_test/locker" > ../adventure_log/locker.log 2>&1 &
```