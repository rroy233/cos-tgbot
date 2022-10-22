#!/usr/bin/env bash

#自动编译+后台运行脚本

check_running(){
        PID=$(pgrep -f "${NAME}")
        if [[ -n ${PID} ]]; then
                return 0
        else
                return 1
        fi
}

NAME="cosTgBot"

echo "即将开始编译"
rm -f $NAME
go build -o $NAME
echo "编译完成"
if check_running; then
        echo -e "$NAME (PID ${PID}) 正在运行，已结束进程。"
        kill -9 "${PID}"
fi
echo "即将开始运行"
ulimit -n 51200 >/dev/null 2>&1
nohup ./$NAME > ./sys.log 2>&1 &
sleep 2s

if check_running; then
        echo -e "$NAME 启动成功 !"
else
        echo -e "$NAME 启动失败 !"
fi