#!/bin/bash

# 循环执行 /app/xray
while true
do
    echo "xray $@"
    /app/xray $@
    echo "xray has been killed, restarting..."
done