#!/bin/bash -eu
export TZ='Asia/Shanghai'

result=0
./collect -r -c "$REPOSITORY"  || result=$?

if [ $result -ne 0 ]; then
  echo "[Err] Revert changes" >&2
  cd "$REPOSITORY" && git reset --hard HEAD
  exit 1
fi

cd "$REPOSITORY" || exit 1

if [[ -n $(git status --porcelain) ]]; then
  git add .
  git commit -m "update at $(date +'%Y-%m-%d %H:%M:%S')"
  git push
fi