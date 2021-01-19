
# issue usdk
# usdk-2db
okexchaincli tx token issue -s usdk -n 1000000000 -w usdk --mintable --from turing --fees 0.02okt -y -b block

# send fees to 1000 lockers
adventure farm allocate-tokens 1000usdk-2db -m "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool" -p "./template/address/farm_test/locker_1000"
adventure farm allocate-tokens 1000okt -m "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool" -p "./template/address/farm_test/locker_1000"

# create pool
okexchaincli tx farm create-pool okt-usdk 0.001usdk-2db okt --from turing --fees 0.02okt -y -b block

okexchaincli tx farm provide okt-usdk 10000okt 10 950 --from turing --fees 0.02okt -y -b block

# start strategy
adventure farm locker strategy-single-pool okt-usdk -p "template/mnemonic/farm_test/locker_1000"
