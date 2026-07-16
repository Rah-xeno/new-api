#!/bin/bash
set -e
cd "$(dirname "$0")"

echo "=== 构建用户端前端 (web) ==="
cd web
bun install --frozen-lockfile
cd default
DISABLE_ESLINT_PLUGIN='true' VITE_REACT_APP_VERSION=$(cat ../../VERSION) bun run build
cd ../..

echo "=== 交叉编译 Go 后端 (Linux amd64) ==="
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o new-api main.go

echo ""
echo "=== 构建完成 ==="
echo "二进制文件: $(pwd)/new-api"
echo ""
echo "部署步骤:"
echo "1. 将 new-api 上传到 1Panel 服务器"
echo "2. 在服务器上运行: chmod +x new-api && ./new-api --log-dir ./logs"
echo "   (推荐用 1Panel 的 supervisor 来管理进程)"
