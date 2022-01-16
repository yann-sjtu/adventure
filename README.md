# Adventure

## 1. 编译
```shell
make
```

## 2. 操作
### 2.1 初始化账户
```shell
adventure evm batch-transfer 10okt -i ${ip} -s ${private_key} -a ${address_file}
```
* -i: ip地址
  * 必填
  * 支持cosmos端口、eth端口
* -s: 私钥 (不是助记词)
  * 必填
  * 对应地址，拥有足够的okt
* -a: 账户地址文件路径
  * 选填，如果为空，代码默认内置2000个固定账户
  * 0x地址格式

### 2.2 压力测试

#### 2.2.1 转账
```shell
adventure evm bench transfer --ips "http://54.249.243.203:26657","http://54.150.103.141:26657","http://18.178.126.121:26657","http://35.73.6.150:26657" -c 100 -t 1 
```

#### 2.2.2 压力测试合约
```shell
adventure evm bench operate --ips "http://54.249.243.203:26657","http://54.150.103.141:26657","http://18.178.126.121:26657","http://35.73.6.150:26657" -c 100 -t 1  --opts 1,1,1,1,1 --times 1 --contract 0x6cc0277c979325800294774d7ae478A96B824271 --id 0 
```

#### 2.2.3 测试网uniswap挖卖提
```shell

```

#### 2.2.4 测试网查询
```shell
 adventure evm bench query --ips https://exchaintestrpc.okex.org -t 1000 -o 1,1,1,1,1,1,1,1,1
```