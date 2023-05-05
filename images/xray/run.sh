#!/bin/bash


PARAMETER=$@

# 如果 ARGS 中包含 listen 字符串
if [[ $PARAMETER =~ "listen" ]]; then
    while true
    do
        echo "xray $PARAMETER"
        sleep 1
        /app/xray $PARAMETER
        echo "xray has been killed, restarting..."
    done
else
  echo "xray $PARAMETER"
  /app/xray PARAMETER
fi

