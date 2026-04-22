#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

REMOTE_NAME="${REMOTE_NAME:-origin}"
REMOTE_URL="${REMOTE_URL:-https://github.com/stellhub/stellar.git}"
BRANCH_NAME="${BRANCH_NAME:-main}"
LOG_DIR="${LOG_DIR:-${REPO_ROOT}/logs}"
LOG_FILE="${LOG_FILE:-${LOG_DIR}/auto-pull.log}"
LOCK_FILE="${LOCK_FILE:-${LOG_DIR}/auto-pull.lock}"

mkdir -p "${LOG_DIR}"

log() {
  printf '[%s] %s\n' "$(date '+%Y-%m-%d %H:%M:%S')" "$1" >> "${LOG_FILE}"
}

exec 9>"${LOCK_FILE}"

if ! flock -n 9; then
  log "检测到已有自动拉取任务正在执行，本次跳过。"
  exit 0
fi

cd "${REPO_ROOT}"

if ! command -v git >/dev/null 2>&1; then
  log "未检测到 git，任务终止。"
  exit 1
fi

CURRENT_REMOTE_URL="$(git remote get-url "${REMOTE_NAME}" 2>/dev/null || true)"

if [[ -z "${CURRENT_REMOTE_URL}" ]]; then
  log "未找到远程仓库 ${REMOTE_NAME}，任务终止。"
  exit 1
fi

if [[ "${CURRENT_REMOTE_URL}" != "${REMOTE_URL}" ]]; then
  log "远程仓库地址不匹配，当前=${CURRENT_REMOTE_URL}，期望=${REMOTE_URL}，本次跳过。"
  exit 1
fi

CURRENT_BRANCH="$(git branch --show-current)"

if [[ "${CURRENT_BRANCH}" != "${BRANCH_NAME}" ]]; then
  log "当前分支为 ${CURRENT_BRANCH}，不是目标分支 ${BRANCH_NAME}，本次跳过。"
  exit 0
fi

if [[ -n "$(git status --porcelain)" ]]; then
  log "检测到工作区存在未提交改动，为避免冲突，本次跳过 git pull。"
  exit 0
fi

log "开始执行 git fetch ${REMOTE_NAME} ${BRANCH_NAME}。"
git fetch "${REMOTE_NAME}" "${BRANCH_NAME}" >> "${LOG_FILE}" 2>&1

LOCAL_HEAD="$(git rev-parse HEAD)"
REMOTE_HEAD="$(git rev-parse "${REMOTE_NAME}/${BRANCH_NAME}")"

if [[ "${LOCAL_HEAD}" == "${REMOTE_HEAD}" ]]; then
  log "当前已是最新版本，无需更新。"
  exit 0
fi

log "发现新提交，开始执行 git pull --ff-only ${REMOTE_NAME} ${BRANCH_NAME}。"
git pull --ff-only "${REMOTE_NAME}" "${BRANCH_NAME}" >> "${LOG_FILE}" 2>&1
log "git pull 执行完成，仓库已更新到最新版本。"
