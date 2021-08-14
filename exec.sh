while :
do
   nohup adventure bench-make-tx --concurrency 150 --sleepTime 3 --rpc-chainid exchain-65 --rpc-hosts https://exchaintesttmrpc.okex.org --privkeyPath /root/adventure/template/privkey/priv_aa --abiPath /root/adventure/template/contract/AttackOperate.abi --contractAddress 0x2a928b71dbaf29b27029259b7ca2acdab4a92bc7 > tx-a.log &
   sleep 1080
   ps -ef | grep adventure | awk '{print $2}' | xargs kill -9

   nohup adventure bench-make-tx --concurrency 150 --sleepTime 3 --rpc-chainid exchain-65 --rpc-hosts https://exchaintesttmrpc.okex.org --privkeyPath /root/adventure/template/privkey/priv_ab --abiPath /root/adventure/template/contract/AttackOperate.abi --contractAddress 0x2a928b71dbaf29b27029259b7ca2acdab4a92bc7 > tx-b.log &
   sleep 1080
   ps -ef | grep adventure | awk '{print $2}' | xargs kill -9

   nohup adventure bench-make-tx --concurrency 150 --sleepTime 3 --rpc-chainid exchain-65 --rpc-hosts https://exchaintesttmrpc.okex.org --privkeyPath /root/adventure/template/privkey/priv_ac --abiPath /root/adventure/template/contract/AttackOperate.abi --contractAddress 0x2a928b71dbaf29b27029259b7ca2acdab4a92bc7 > tx-c.log &
   sleep 1080
   ps -ef | grep adventure | awk '{print $2}' | xargs kill -9

   nohup adventure bench-make-tx --concurrency 150 --sleepTime 3 --rpc-chainid exchain-65 --rpc-hosts https://exchaintesttmrpc.okex.org --privkeyPath /root/adventure/template/privkey/priv_ad --abiPath /root/adventure/template/contract/AttackOperate.abi --contractAddress 0x2a928b71dbaf29b27029259b7ca2acdab4a92bc7 > tx-d.log &
   sleep 1080
   ps -ef | grep adventure | awk '{print $2}' | xargs kill -9

   nohup adventure bench-make-tx --concurrency 150 --sleepTime 3 --rpc-chainid exchain-65 --rpc-hosts https://exchaintesttmrpc.okex.org --privkeyPath /root/adventure/template/privkey/priv_ae --abiPath /root/adventure/template/contract/AttackOperate.abi --contractAddress 0x2a928b71dbaf29b27029259b7ca2acdab4a92bc7 > tx-e.log &
   sleep 1080
   ps -ef | grep adventure | awk '{print $2}' | xargs kill -9

done