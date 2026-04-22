#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

HOST="${HOST:-0.0.0.0}"
PORT="${PORT:-5173}"
INSTALL_DEPS="${INSTALL_DEPS:-true}"
LOG_DIR="${LOG_DIR:-${REPO_ROOT}/logs}"
LOG_FILE="${LOG_FILE:-${LOG_DIR}/vitepress.log}"
PID_FILE="${PID_FILE:-${LOG_DIR}/vitepress.pid}"

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

if [[ "${INSTALL_DEPS}" == "true" ]]; then
  if [[ -f package-lock.json ]]; then
    echo "检测到 package-lock.json，执行 npm ci ..."
    npm ci
  else
    echo "未检测到 package-lock.json，执行 npm install ..."
    npm install
  fi
fi

mkdir -p "${LOG_DIR}"

if [[ -f "${PID_FILE}" ]]; then
  OLD_PID="$(cat "${PID_FILE}")"
  if [[ -n "${OLD_PID}" ]] && kill -0 "${OLD_PID}" >/dev/null 2>&1; then
    echo "VitePress 已在后台运行，PID=${OLD_PID}" >&2
    echo "日志文件: ${LOG_FILE}" >&2
    exit 1
  fi
  rm -f "${PID_FILE}"
fi

echo "以后台模式启动 VitePress: http://${HOST}:${PORT}"
echo "日志文件: ${LOG_FILE}"

nohup npm run docs:dev -- --host "${HOST}" --port "${PORT}" >"${LOG_FILE}" 2>&1 &
APP_PID=$!

echo "${APP_PID}" > "${PID_FILE}"

sleep 2

if ! kill -0 "${APP_PID}" >/dev/null 2>&1; then
  echo "启动失败，请检查日志: ${LOG_FILE}" >&2
  rm -f "${PID_FILE}"
  exit 1
fi

echo "启动成功，PID=${APP_PID}"
