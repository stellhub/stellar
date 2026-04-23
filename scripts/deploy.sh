#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

REMOTE_NAME="${REMOTE_NAME:-origin}"
BRANCH_NAME="${BRANCH_NAME:-main}"
PUBLISH_DIR="${PUBLISH_DIR:-/var/www/stellar}"
SYNC_NGINX_CONFIG="${SYNC_NGINX_CONFIG:-true}"
NGINX_CONF_PATH="${NGINX_CONF_PATH:-/etc/nginx/nginx.conf}"
LOG_DIR="${LOG_DIR:-${REPO_ROOT}/logs}"
LOG_FILE="${LOG_FILE:-${LOG_DIR}/deploy.log}"
LOCK_FILE="${LOCK_FILE:-${LOG_DIR}/deploy.lock}"

mkdir -p "${LOG_DIR}"

log() {
  printf '[%s] %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$1" | tee -a "${LOG_FILE}"
}

require_command() {
  local command_name="$1"
  if ! command -v "${command_name}" >/dev/null 2>&1; then
    log "未检测到命令 ${command_name}，部署终止。"
    exit 1
  fi
}

exec 9>"${LOCK_FILE}"

if ! flock -n 9; then
  log "检测到已有部署任务正在执行，本次跳过。"
  exit 0
fi

cd "${REPO_ROOT}"

require_command git
require_command node
require_command npm
require_command rsync
require_command sudo
require_command flock

if ! sudo -n true >/dev/null 2>&1; then
  log "当前用户无法免密执行 sudo，部署终止。"
  exit 1
fi

NODE_MAJOR="$(node -p "process.versions.node.split('.')[0]")"
if [[ "${NODE_MAJOR}" -lt 20 ]]; then
  log "当前 Node.js 版本为 $(node -v)，需要 Node.js 20 或更高版本。"
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  log "检测到工作区存在未提交改动，部署终止。"
  exit 1
fi

log "开始获取远端分支 ${REMOTE_NAME}/${BRANCH_NAME}。"
git fetch "${REMOTE_NAME}" "${BRANCH_NAME}" >> "${LOG_FILE}" 2>&1

LOCAL_HEAD="$(git rev-parse HEAD)"
REMOTE_HEAD="$(git rev-parse "${REMOTE_NAME}/${BRANCH_NAME}")"

if [[ "${LOCAL_HEAD}" != "${REMOTE_HEAD}" ]]; then
  log "发现新提交，开始执行 git pull --ff-only。"
  git pull --ff-only "${REMOTE_NAME}" "${BRANCH_NAME}" >> "${LOG_FILE}" 2>&1
else
  log "当前代码已是最新版本。"
fi

if [[ -f package-lock.json ]]; then
  log "检测到 package-lock.json，开始执行 npm ci。"
  npm ci >> "${LOG_FILE}" 2>&1
else
  log "未检测到 package-lock.json，开始执行 npm install。"
  npm install >> "${LOG_FILE}" 2>&1
fi

log "开始构建 VitePress 静态站点。"
npm run docs:build >> "${LOG_FILE}" 2>&1

log "开始同步构建产物到 ${PUBLISH_DIR}。"
sudo mkdir -p "${PUBLISH_DIR}"
sudo rsync -av --delete "${REPO_ROOT}/docs/.vitepress/dist/" "${PUBLISH_DIR}/" >> "${LOG_FILE}" 2>&1

if [[ "${SYNC_NGINX_CONFIG}" == "true" ]]; then
  log "开始校验仓库内 Nginx 配置。"
  sudo nginx -t -c "${REPO_ROOT}/nginx.conf" >> "${LOG_FILE}" 2>&1

  log "开始同步 nginx.conf 到 ${NGINX_CONF_PATH}。"
  sudo cp "${REPO_ROOT}/nginx.conf" "${NGINX_CONF_PATH}"
fi

log "开始校验当前 Nginx 配置。"
sudo nginx -t >> "${LOG_FILE}" 2>&1

log "开始重载 Nginx。"
sudo systemctl reload nginx >> "${LOG_FILE}" 2>&1

log "部署完成。"
