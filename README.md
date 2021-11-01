## Adventure

### 1. 简介

### 2. 编译与配置
编译：

```shell
$ make adventure
```

配置` config.toml`，一般不做修改。

config.toml 是基础配置文件，存放测试节点：

* 注1：adventure工程下 `template/*.json`已有许多样例，

* 注2：adventure工程下 `template/mnemonic/*`存放助记词，`template/addr/*`存放对应的地址(用于查询、转账等)

### 3. 使用说明


### 启动全部交易类型测试
#### 测试账户转账
```shell script
adventure account send -p template/address/captain -a 1000000okt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/normal_5 -a 1000000okt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/normal_1000_1 -a 100000okt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/normal_1000_2 -a 10000okt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/proxy_10 -a 10000okt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"
adventure account send -p template/address/normal_100 -a 10000okt -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient"