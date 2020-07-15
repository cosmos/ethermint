#!/bin/bash

# change this script so that the date is appended to the beginning of each line (not as a new line)
tmux kill-session -t resource_recorder; tmux new -s resource_recorder -d; tmux send-keys -t resource_recorder "while true; do (date '+%s' && ps -e -o pcpu,pmem,args --sort=pcpu | grep \"emint\" | grep -v grep | cut -d\" \" -f1-5 | tail) >> resource.log; sleep 1; done" C-m

go build
bash benchmark_init.sh

sleep 10

./benchmarking s -c $1

sleep 30

tmux kill-session -t resource_recorder