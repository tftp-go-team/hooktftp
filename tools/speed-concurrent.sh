#!/bin/sh

for i in $(seq 10)
do
    sh tools/speed.sh &
done

wait
