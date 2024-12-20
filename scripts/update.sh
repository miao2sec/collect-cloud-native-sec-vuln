#!/bin/bash -eu


COMMIT_MSG=$1



if [ -z "$COMMIT_MSG" ]; then
  echo "commit message required"
  exit 1
fi

result=0
./collect -r  || result=$?

if [ $result -ne 0 ]; then
  echo "[Err] Revert changes" >&2
  git reset --hard HEAD
  exit 1
fi

if [[ -n $(git status --porcelain) ]]; then
  git add .
  git commit -m "${COMMIT_MSG}"
  git push
fi