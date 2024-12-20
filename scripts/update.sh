#!/bin/bash -eu


COMMIT_MSG=$1



if [ -z "$COMMIT_MSG" ]; then
  echo "commit message required"
  exit 1
fi

result=0
./collect -r -c data  || result=$?

if [ $result -ne 0 ]; then
  echo "[Err] Revert changes" >&2
  cd data && git reset --hard HEAD
  exit 1
fi

cd data || exit 1

if [[ -n $(git status --porcelain) ]]; then
  git add .
  git commit -m "${COMMIT_MSG}"
  git push
fi