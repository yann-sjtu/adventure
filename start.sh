#!/bin/bash

HOSTS_IN=(
10.0.240.20
10.0.240.21
10.0.240.22
10.0.240.23
10.0.240.24
10.0.240.25
10.0.240.26
10.0.240.27
10.0.240.28
10.0.240.29
10.0.240.30
10.0.240.31
10.0.240.32
10.0.240.33
10.0.240.34
10.0.240.35
10.0.240.36
10.0.240.39
10.0.240.40
10.0.240.41
10.0.240.42
)

# shellcheck disable=SC2068
for host in ${HOSTS_IN[@]}
do
#  scp /Users/green/project/okex/deploy-new/env/mainnet/monitor-aws/agent.yml root@${host}:/root

  echo ${host}
    ssh root@"${host}" << eeooff

#    git config --global credential.helper store
#    rm -rf /root/adventure
#    export https_proxy=http://10.0.66.40:7890 http_proxy=http://10.0.66.40:7890 all_proxy=socks5://10.0.66.40:7891
#    git clone https://github.com/okex/adventure.git -b main-refactor
#    cd /root/adventure
#    mkdir /root/.adventure
#    cp /root/adventure/config.toml /root/.adventure
#    export GOPROXY="http://goproxy.cn"
#    nohup make install >> /root/adventure/make.log 2>&1 &
#    adventure version  | grep commit

#    echo "root soft nofile 1000000
#root hard nofile 1000000
#* soft nofile 1000000
#* hard nofile 1000000" >> /etc/security/limits.conf
#    cat /etc/security/limits.conf | grep 1000000
#    ulimit -a | grep "open files"

#    echo "47.75.105.229 exchaintestrpc.okex.org" >> /etc/hosts
#    cat /etc/hosts | grep exchaintestrpc


#     systemctl start docker
#     docker-compose -f /root/agent.yml up -d

#    ss -a | wc -l
#    ps -ef | grep adventure | grep -v grep | awk '{print $2}'
#    pgrep adventure | xargs kill -s 9

#rm -f /root/adventure/query.log

#nohup adventure bench-query -g 200,200,200,200,200,200,200,200 -t 2 >> /root/adventure/query.log 2>&1 &

    exit
eeooff
done
