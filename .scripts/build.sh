#!/bin/bash
# shellcheck disable=SC2034
set -euo pipefail
IFS=$'\n\t'

export LDFLAGS=(
  "-w" "-s"
  "'-X=main.BuildInfo=$(git show -s --format=$'%H @ %cI' HEAD) @ $(date --iso-8601=seconds)'"
)

function install_packages() {
  local DEBIAN_FRONTEND=noninteractive
  local PACKAGES=(
    "gcc-arm-linux-gnueabihf"
    "gcc-aarch64-linux-gnu"
    "gcc-mingw-w64-x86-64"
  )
  apt-get update >/dev/null
  apt-get install -qq -y "${PACKAGES[@]}" >/dev/null
}

function build() {
  local IFS=" "
  export \
    CGO_ENABLED=1 \
    GOOS="$1" \
    GOARCH="$2" \
    NAME="goip-agent-$3"
  case "$GOOS" in
  linux)
    case "$GOARCH" in
    arm) export CC="arm-linux-gnueabihf-gcc" ;;
    arm64) export CC="aarch64-linux-gnu-gcc" ;;
    esac
    ;;
  windows) export CC="x86_64-w64-mingw32-gcc" ;;
  esac
  go build \
    -ldflags """${LDFLAGS[*]}""" \
    -o "$NAME" \
    ./cmd/goip-agent
  unset CC
}

install_packages
build "linux" "arm" "arm"
build "linux" "arm64" "arm64"
build "linux" "amd64" "x64"
build "windows" "amd64" "x64.exe"
du -hs goip-agent-*
sha1sum goip-agent-*
