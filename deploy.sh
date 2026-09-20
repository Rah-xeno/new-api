#!/bin/bash
set -euo pipefail

# 固定项目目录
PROJECT_PATH="/www/newapi2"
LOCK_FILE="/tmp/deploy_newapi2.lock"

# 加文件锁，防止短时间多次push并发执行
exec 200>"${LOCK_FILE}"
if ! flock -n 200 ; then
  echo "已有部署任务正在运行，本次部署终止"
  exit 0
fi

echo "===== 部署开始 $(date) ====="
cd "${PROJECT_PATH}"

echo "1. git pull 更新代码"
git pull

echo "2. 执行构建脚本 sh build.sh"
sh build.sh

echo "3. 重启1Panel编排项目 newapi2"
docker restart "newapi2"

echo "===== 部署完成 $(date) ====="