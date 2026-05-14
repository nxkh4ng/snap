#!/bin/bash
# build.sh - Build and archive snap binaries for distribution
# Usage: ./build.sh <version-number>
# Example: ./build.sh 0.4.0

set -euo pipefail
shopt -s inherit_errexit
trap 'echo "Error on line $LINENO"' ERR

main() {
    local PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    local DIST_DIR="$PROJECT_DIR/dist"

    if [ $# -eq 0 ]; then
        echo "Usage: ./build.sh <version-number>"
        echo "Example: ./build.sh 0.4.0"
        exit 1
    fi
    local -r VERSION=$1

    rm -rf "$DIST_DIR"
    mkdir -p "$DIST_DIR"
    rm -rf "/tmp/snap-build"

    mkdir -p "/tmp/snap-build/linux_amd64"
    mkdir -p "/tmp/snap-build/linux_arm64"
    GOOS=linux GOARCH=amd64 go build -C "$PROJECT_DIR" -ldflags="-s -w -X github.com/nxkh4ng/snap/cmd.version=$VERSION" -o "/tmp/snap-build/linux_amd64/snap"
    GOOS=linux GOARCH=arm64 go build -C "$PROJECT_DIR" -ldflags="-s -w -X github.com/nxkh4ng/snap/cmd.version=$VERSION" -o "/tmp/snap-build/linux_arm64/snap"

    mkdir -p "/tmp/snap-build/darwin_amd64"
    mkdir -p "/tmp/snap-build/darwin_arm64"
    GOOS=darwin GOARCH=amd64 go build -C "$PROJECT_DIR" -ldflags="-s -w -X github.com/nxkh4ng/snap/cmd.version=$VERSION" -o "/tmp/snap-build/darwin_amd64/snap"
    GOOS=darwin GOARCH=arm64 go build -C "$PROJECT_DIR" -ldflags="-s -w -X github.com/nxkh4ng/snap/cmd.version=$VERSION" -o "/tmp/snap-build/darwin_arm64/snap"

    mkdir -p "/tmp/snap-build/windows_x86_64"
    mkdir -p "/tmp/snap-build/windows_x86"
    GOOS=windows GOARCH=amd64 go build -C "$PROJECT_DIR" -ldflags="-s -w -X github.com/nxkh4ng/snap/cmd.version=$VERSION" -o "/tmp/snap-build/windows_x86_64/snap.exe"
    GOOS=windows GOARCH=386 go build -C "$PROJECT_DIR" -ldflags="-s -w -X github.com/nxkh4ng/snap/cmd.version=$VERSION" -o "/tmp/snap-build/windows_x86/snap.exe"

    tar -czf "$DIST_DIR/snap_linux_amd64.tar.gz" -C "/tmp/snap-build/linux_amd64" "snap"
    tar -czf "$DIST_DIR/snap_linux_arm64.tar.gz" -C "/tmp/snap-build/linux_arm64" "snap"
    tar -czf "$DIST_DIR/snap_darwin_amd64.tar.gz" -C "/tmp/snap-build/darwin_amd64" "snap"
    tar -czf "$DIST_DIR/snap_darwin_arm64.tar.gz" -C "/tmp/snap-build/darwin_arm64" "snap"
    zip -j "$DIST_DIR/snap_windows_x86_64.zip" "/tmp/snap-build/windows_x86_64/snap.exe"
    zip -j "$DIST_DIR/snap_windows_x86.zip" "/tmp/snap-build/windows_x86/snap.exe"

    cd "$DIST_DIR"
    shasum -a 256 *.tar.gz *.zip > "snap_${VERSION}_checksums.txt"
}

main "$@"
