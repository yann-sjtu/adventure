#!/bin/zsh

adventure monitor shares-control -p ./x/monitor/shares-control/config.toml

nohup adventure monitor shares-control -p ./x/monitor/shares-control/config.toml > /root/monitor/monitor.log 2>&1 &

tail -f ../monitor/monitor.log