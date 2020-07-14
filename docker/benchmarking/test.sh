#!/bin/bash

tmux kill-session -t resource_recorder; tmux new -s resource_recorder -d; tmux send-keys -t resource_recorder "while true; do (date '+%s' && ps -e -o pcpu,pmem,args --sort=pcpu | grep \"emint\" | grep -v grep | cut -d\" \" -f1-10 | tail) >> resource.log; sleep 1; done" C-m

logger() {
  (date '+%s' && ps -e -o pcpu,pmem,args --sort=pcpu | grep "emint" | grep -v grep | cut -d" " -f1-10 | tail) >> resource.log
  sleep 1
}

go build
bash benchmark_init.sh

sleep 5

./benchmarking s -c $1

sleep 30

tmux kill-session -t resource_recorder