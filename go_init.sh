#!/bin/bash
set -e

wget https://studygolang.com/dl/golang/go1.16.12.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.12.linux-amd64.tar.gz
sed -i "1iexport PATH=\$PATH:/usr/local/go/bin:/root/go/bin" ~/.bashrc
source ~/.bashrc
rm -f go1.16.12.linux-amd64.tar.gz

apt-get update
apt install make gcc
make