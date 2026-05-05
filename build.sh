#!/bin/bash

PROJECT_DIR="/home/kh4ng/dev/snap"

build_zip() {
    local os=$1
    local arch=$2
    local name=$3
    local ext=$4
    local binary="snap$ext"
    
    GOOS=$os GOARCH=$arch go build -C "$PROJECT_DIR" -ldflags="-s -w" -o "/tmp/$binary"
    cd /tmp && zip -r "$PROJECT_DIR/dist/${name}.zip" "$binary"
    rm "/tmp/$binary"
}

mkdir -p "$PROJECT_DIR/dist"
build_zip linux amd64 snap_linux_amd64 ""
build_zip linux arm64 snap_linux_arm64 ""
build_zip darwin amd64 snap_darwin_amd64 ""
build_zip darwin arm64 snap_darwin_arm64 ""
build_zip windows amd64 snap_windows_x86_64 ".exe"
