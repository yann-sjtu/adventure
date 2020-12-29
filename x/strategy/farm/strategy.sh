#!/usr/bin/env bash

# okexchain1ntvyep3suq5z7789g7d5dejwzameu08m6gh7yl

# init
#adventure farm allocate-tokens 7000"okt -p "./template/address/farm_test/pooler_total" -n 100
adventure account send -r "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool" -p "./template/address/farm_test/pooler_total" -a 7000"okt
adventure account send -r "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool" -p "./template/address/farm_test/locker_total" -a 1000"okt

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
adventure account send -r "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer" -p "./template/address/farm_test/pooler_total" -a 7000"okt
adventure account send -r "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool" -p "./template/address/farm_test/locker_total" -a 1000"okt

adventure farm pooler issue-token -p "./template/mnemonic/farm_test/pooler"
adventure farm allocate-tokens 10000usdk -p "./template/address/farm_test/pooler_total"

adventure farm pooler create-pair -p "./template/mnemonic/farm_test/pooler"
adventure farm pooler add-liquidity -p "./template/mnemonic/farm_test/pooler"

adventure farm locker allocate-tokens-to-lockers-from-all-poolers -p "./template/mnemonic/farm_test/pooler" -l "./template/address/farm_test/locker_total"

# strategy
adventure farm pooler strategy-pooler -p "./template/mnemonic/farm_test/pooler"
adventure farm locker strategy-lock-unlock -p "./template/mnemonic/farm_test/locker"
