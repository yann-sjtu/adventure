#!/bin/zsh

adventure monitor reward-plunderer -p ./x/monitor/reward_plunderer/config.toml
# local
nohup adventure monitor reward-plunderer -p ./x/monitor/reward_plunderer/config.toml > ~/log/monitor.log 2>&1 &

tail -f ~/log/monitor.log


#nohup adventure monitor final-top21-shares-control -p ./x/monitor/final_top_21_control/config.toml > /root/monitor/monitor.log 2>&1 &
#tail -f ../monitor/monitor.log


adventure monitor tools v-a