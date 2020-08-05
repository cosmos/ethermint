#!/bin/bash

tmux kill-session -t resource_recorder 1> /dev/null; tmux new -s resource_recorder -d; tmux send-keys -t resource_recorder "while true; do (date '+%s' && ps -e -o pcpu,pmem,args --sort=pcpu | grep \"geth\" | grep -v grep | cut -d\" \" -f1-3 | tail) >> resource.log; sleep 1; done" C-m

bash init_geth.sh 2> /dev/null

sleep 10

../benchmarking s -c $1 -p geth

sleep 30
tmux kill-session -t resource_recorder

START=$(cat start.txt)
END=$(cat end.txt)
../benchmarking a -s $START -e $END -p geth

rm start.txt end.txt