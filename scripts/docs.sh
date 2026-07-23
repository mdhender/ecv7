#!/bin/bash

if [ "${1}" == build ]; then
  cd docs || exit 2
  hugo build
  exit 0
fi

echo "usage: scripts/docs.sh [serve]"
exit 2
if [ "${1}" == serve ]; then
  cd docs || exit 2
  hugo serve
  exit 0
fi

echo "usage: scripts/docs.sh [serve]"
exit 2
