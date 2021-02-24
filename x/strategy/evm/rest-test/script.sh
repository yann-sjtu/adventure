adventure evm rest-test deploy -u http://10.0.240.22:26659
# get the contract address
adventure evm rest-test run-erc20-transfer-test -c 0xb4eb8fc3ab329c96a182ccd27671a41d46ba3a68 -u http://10.0.240.22:26659

nohup adventure evm rest-test run-erc20-transfer-test -c 0xb4eb8fc3ab329c96a182ccd27671a41d46ba3a68 -u http://10.0.240.22:26659 > ../adventure_log/rest_test.log 2>&1 &