#!/usr/bin/env bash

adventure account address -p "./template/mnemonic/farm_test/pooler_1" -f "./template/address/farm_test/pooler_1"
adventure account address -p "./template/mnemonic/farm_test/pooler_2" -f "./template/address/farm_test/pooler_2"
adventure account address -p "./template/mnemonic/farm_test/pooler_3" -f "./template/address/farm_test/pooler_3"
adventure account address -p "./template/mnemonic/farm_test/pooler_4" -f "./template/address/farm_test/pooler_4"
adventure account address -p "./template/mnemonic/farm_test/pooler_5" -f "./template/address/farm_test/pooler_5"
adventure account address -p "./template/mnemonic/farm_test/pooler_6" -f "./template/address/farm_test/pooler_6"
adventure account address -p "./template/mnemonic/farm_test/destructive_locker" -f "./template/address/farm_test/destructive_locker"

# start

# send to all pooler
adventure account send -r "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -p "./template/address/farm_test/pooler_total" -a 7000okt
# send to destructive locker
adventure account send -r "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -p "./template/address/farm_test/destructive_locker" -a 100okt

okexchaincli query account okexchain10q0rk5qnyag7wfvvt7rtphlw589m7frsku8qc9

# issue tokens
adventure farm pooler issue-token -p "./template/mnemonic/farm_test/destructive_pooler"
adventure farm pooler issue-token -p "./template/mnemonic/farm_test/pooler_1"
adventure farm pooler issue-token -p "./template/mnemonic/farm_test/pooler_2"
adventure farm pooler issue-token -p "./template/mnemonic/farm_test/pooler_3"

# create swap pair
adventure farm pooler create-pair -p "./template/mnemonic/farm_test/destructive_pooler"
adventure farm pooler create-pair -p "./template/mnemonic/farm_test/pooler_1"
adventure farm pooler create-pair -p "./template/mnemonic/farm_test/pooler_2"
adventure farm pooler create-pair -p "./template/mnemonic/farm_test/pooler_3"

# add liquidity and get lpt
# waiting
adventure farm pooler add-liquidity -p "./template/mnemonic/farm_test/destructive_pooler"
adventure farm pooler add-liquidity -p "./template/mnemonic/farm_test/pooler_1"
adventure farm pooler add-liquidity -p "./template/mnemonic/farm_test/pooler_2"
adventure farm pooler add-liquidity -p "./template/mnemonic/farm_test/pooler_3"

# create && provide farm pool
# waiting
adventure farm pooler create-provide-farm-pool -p "./template/mnemonic/farm_test/destructive_pooler" -d -w
okexchaincli query farm whitelist
adventure farm pooler create-provide-farm-pool -p "./template/mnemonic/farm_test/pooler_1" -w
adventure farm pooler create-provide-farm-pool -p "./template/mnemonic/farm_test/pooler_2"
adventure farm pooler create-provide-farm-pool -p "./template/mnemonic/farm_test/pooler_3"

okexchaincli query farm pool-num
okexchaincli query farm pools

# strategy destructive pooler

adventure farm pooler strategy-destructive-pooler -p "./template/mnemonic/farm_test/destructive_pooler"

okexchaincli query account okexchain1kg57hzyrjwxeymv2cnpplltecgmvyhxdy88qlv