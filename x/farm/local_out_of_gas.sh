#!/usr/bin/env bash

okexchaincli tx send richer okexchain10q0rk5qnyag7wfvvt7rtphlw589m7frsku8qc9 899900000okt -y -b block --fees 0.002okt

adventure farm allocate-tokens 3000000okt -p "./template/address/farm_test/local_pooler_out_of_gas"

adventure farm issue-tokens -p "./template/mnemonic/farm_test/local_pooler_out_of_gas"
