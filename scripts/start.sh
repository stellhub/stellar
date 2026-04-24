#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

cd "${REPO_ROOT}"

if ! command -v node >/dev/null 2>&1; then
  echo "错误: 未检测到 node，请先安装 Node.js 20 或更高版本。" >&2
  exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
  echo "错误: 未检测到 npm，请先安装 npm。" >&2
  exit 1
fi

NODE_MAJOR="$(node -p "process.versions.node.split('.')[0]")"

if [[ "${NODE_MAJOR}" -lt 20 ]]; then
  echo "错误: 当前 Node.js 版本为 $(node -v)，需要 Node.js 20 或更高版本。" >&2
  exit 1
fi

if [[ ! -f "${REPO_ROOT}/package-lock.json" ]]; then
  echo "错误: 未检测到 package-lock.json，无法继续执行安装。" >&2
  exit 1
fi

echo "Node.js 环境校验通过: $(node -v)"
echo "package-lock.json 校验通过。"

echo "开始执行 npm install ..."
npm install

echo "开始执行 npm run docs:build ..."
npm run docs:build

echo "构建完成。"
