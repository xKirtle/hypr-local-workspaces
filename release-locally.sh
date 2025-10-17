#!/usr/bin/env bash
set -euo pipefail

go build -o hypr-local-workspaces ./cmd/hypr-local-workspaces
sudo install -Dm755 "hypr-local-workspaces" "/usr/bin/hypr-local-workspaces"
ls -l /usr/bin | grep hypr-local-workspaces
rm hypr-local-workspaces