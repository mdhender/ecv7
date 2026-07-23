#!/usr/bin/env bash
#
# docs.sh — update the docs checkout and rebuild the live site, on the server.
#
# Installed at /opt/ecv7/scripts/docs.sh. Pulls the latest commit into the ecv7
# checkout and rebuilds Hugo's public/ directory, which the web server serves
# directly (see deploy/nginx.conf or deploy/Caddyfile).
#
# The rebuild happens in place on the live public/ directory, so a browser
# loading a page mid-build may briefly hit a missing or half-written asset.
# This is accepted as low risk (low probability, low impact).
#
# One-time server setup is documented in deploy/README.md.
#
# Usage (on the server):
#   /opt/ecv7/scripts/docs.sh deploy
#
set -euo pipefail

if [ "${1}" == build ]; then
  cd docs || exit 2
  hugo build
  exit 0
fi

if [ "${1}" == deploy ]; then
  # Where the ecv7 repo is checked out on this server.
  # The web server's docs root points at "$REPO_DIR/docs/public".
  REPO_DIR="/opt/ecv7"

  # deploy-docs.sh invokes this over SSH, and a non-interactive SSH shell has a
  # minimal PATH that omits /usr/local/bin, the Go toolchain, and /snap/bin. Add
  # the common locations so git/hugo/go resolve however they were installed.
  export PATH="/usr/local/bin:/usr/local/go/bin:/snap/bin:$PATH"

  echo ">> Updating checkout in ${REPO_DIR}"
  cd ${REPO_DIR} || exit 2
  git pull --ff-only

  cd ${REPO_DIR}/docs || exit 2
  echo ">> Rebuilding docs/public/"
  hugo --gc --minify
  echo ">> Done"

  exit 0
fi

if [ "${1}" == serve ]; then
  cd docs || exit 2
  hugo serve
  exit 0
fi

echo "usage: scripts/docs.sh [build | deploy | serve]"
exit 2
