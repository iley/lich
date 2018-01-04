#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

if [[ -z ${GOPATH:-} ]]; then
  export GOPATH="$HOME/go"
fi

REV_SHORT=$(git rev-parse --short HEAD)
DATE=$(date +"%Y-%m-%d %H:%M:%S")
VERSION_STRING="$REV_SHORT built on $DATE"

go build -o lich -ldflags "-X main.version '$VERSION_STRING'"

sudo supervisorctl stop lich
sudo cp lich /opt/lich/lich
sudo supervisorctl start lich
