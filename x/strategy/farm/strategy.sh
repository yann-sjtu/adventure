#!/usr/bin/env bash

# okexchain1ntvyep3suq5z7789g7d5dejwzameu08m6gh7yl

# init
#adventure farm allocate-tokens 7000"okt -p "./template/address/farm_test/pooler_total" -n 100
adventure account send -r "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool" -p "./template/address/farm_test/pooler_total" -a 7000okt
adventure account send -r "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool" -p "./template/address/farm_test/locker_total" -a 1000okt

adventure farm pooler issue-token -p "./template/mnemonic/farm_test/pooler"
adventure farm allocate-tokens 10000usdk -p "./template/address/farm_test/pooler_total"

adventure farm pooler create-pair -p "./template/mnemonic/farm_test/pooler"
adventure farm pooler add-liquidity -p "./template/mnemonic/farm_test/pooler"

adventure farm locker allocate-tokens-to-lockers-from-all-poolers -p "./template/mnemonic/farm_test/pooler" -l "./template/address/farm_test/locker_total"

# strategy
adventure farm pooler strategy-pooler -p "./template/mnemonic/farm_test/pooler"
adventure farm locker strategy-lock-unlock -p "./template/mnemonic/farm_test/locker"


# 系统测试
# init
#adventure farm allocate-tokens 7000"okt -p "./template/address/farm_test/pooler_total" -n 100
adventure account send -r "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -p "./template/address/farm_test/pooler_total" -a 7000okt
adventure account send -r "actual assume crew creek furnace water electric fitness stumble usage embark ancient" -p "./template/address/farm_test/locker_total" -a 1000okt

adventure farm pooler issue-token -p "./template/mnemonic/farm_test/pooler"
adventure farm allocate-tokens 10000usdk -p "./template/address/farm_test/pooler_total"

adventure farm pooler create-pair -p "./template/mnemonic/farm_test/pooler"
adventure farm pooler add-liquidity -p "./template/mnemonic/farm_test/pooler"

adventure farm locker allocate-tokens-to-lockers-from-all-poolers -p "./template/mnemonic/farm_test/pooler" -l "./template/address/farm_test/locker_total"

# strategy
adventure farm pooler strategy-pooler -p "./template/mnemonic/farm_test/pooler"
adventure farm locker strategy-lock-unlock -p "./template/mnemonic/farm_test/locker"

nohup adventure farm pooler strategy-pooler -p "./template/mnemonic/farm_test/pooler" > /root/adventure/pooler.log 2>&1 &
tail -f ../adventure/pooler.log

nohup adventure farm locker strategy-lock-unlock -p "./template/mnemonic/farm_test/locker" > /root/adventure/locker.log 2>&1 &
tail -f ../adventure/locker.log
