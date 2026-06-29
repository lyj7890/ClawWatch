#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-run}"
APP_NAME="ClawPet"
BUNDLE_ID="com.clawwatch.ClawPet"
MIN_SYSTEM_VERSION="14.0"
APP_VERSION="1.0.0"
BUILD_VERSION="1"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACKAGE_DIR="$ROOT_DIR/macos/ClawPet"
DIST_DIR="$ROOT_DIR/dist"
CACHE_DIR="$ROOT_DIR/.build-cache"
APP_BUNDLE="$DIST_DIR/$APP_NAME.app"
APP_CONTENTS="$APP_BUNDLE/Contents"
APP_MACOS="$APP_CONTENTS/MacOS"
APP_RESOURCES="$APP_CONTENTS/Resources"
APP_BINARY="$APP_MACOS/$APP_NAME"
INFO_PLIST="$APP_CONTENTS/Info.plist"
ICON_FILE="$APP_RESOURCES/ClawPet.icns"

pkill -x "$APP_NAME" >/dev/null 2>&1 || true

export CLANG_MODULE_CACHE_PATH="$CACHE_DIR/clang"
export SWIFTPM_MODULECACHE_OVERRIDE="$CACHE_DIR/swiftpm"

rm -rf "$APP_BUNDLE"
mkdir -p "$APP_MACOS" "$APP_RESOURCES"

if swift build --package-path "$PACKAGE_DIR"; then
  BUILD_BINARY="$(swift build --package-path "$PACKAGE_DIR" --show-bin-path)/$APP_NAME"
  cp "$BUILD_BINARY" "$APP_BINARY"
else
  echo "SwiftPM unavailable; falling back to direct swiftc build." >&2
  swiftc \
    -parse-as-library \
    -target arm64-apple-macosx14.0 \
    -o "$APP_BINARY" \
    $(find "$PACKAGE_DIR/Sources" -name '*.swift' -print)
fi

chmod +x "$APP_BINARY"

ICON_WORK_DIR="$CACHE_DIR/ClawPet.iconset"
ICON_SOURCE="$CACHE_DIR/ClawPet-1024.png"
rm -rf "$ICON_WORK_DIR"
mkdir -p "$ICON_WORK_DIR"
swift "$ROOT_DIR/script/generate_app_icon.swift" "$ICON_SOURCE"
while read -r pixels filename; do
  sips -z "$pixels" "$pixels" "$ICON_SOURCE" --out "$ICON_WORK_DIR/$filename" >/dev/null
done <<'SIZES'
16 icon_16x16.png
32 icon_16x16@2x.png
32 icon_32x32.png
64 icon_32x32@2x.png
128 icon_128x128.png
256 icon_128x128@2x.png
256 icon_256x256.png
512 icon_256x256@2x.png
512 icon_512x512.png
1024 icon_512x512@2x.png
SIZES
iconutil -c icns "$ICON_WORK_DIR" -o "$ICON_FILE"

cat >"$INFO_PLIST" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleExecutable</key>
  <string>$APP_NAME</string>
  <key>CFBundleIdentifier</key>
  <string>$BUNDLE_ID</string>
  <key>CFBundleName</key>
  <string>$APP_NAME</string>
  <key>CFBundleShortVersionString</key>
  <string>$APP_VERSION</string>
  <key>CFBundleVersion</key>
  <string>$BUILD_VERSION</string>
  <key>CFBundleIconFile</key>
  <string>ClawPet</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>LSMinimumSystemVersion</key>
  <string>$MIN_SYSTEM_VERSION</string>
  <key>LSUIElement</key>
  <true/>
  <key>NSPrincipalClass</key>
  <string>NSApplication</string>
</dict>
</plist>
PLIST

codesign --force --deep --sign - "$APP_BUNDLE"

open_app() {
  /usr/bin/open -n "$APP_BUNDLE"
}

case "$MODE" in
  run)
    open_app
    ;;
  --debug|debug)
    lldb -- "$APP_BINARY"
    ;;
  --logs|logs)
    open_app
    /usr/bin/log stream --info --style compact --predicate "process == \"$APP_NAME\""
    ;;
  --telemetry|telemetry)
    open_app
    /usr/bin/log stream --info --style compact --predicate "subsystem == \"$BUNDLE_ID\""
    ;;
  --verify|verify)
    open_app
    sleep 1
    pgrep -x "$APP_NAME" >/dev/null
    ;;
  --package|package)
    ;;
  *)
    echo "usage: $0 [run|--debug|--logs|--telemetry|--verify|--package]" >&2
    exit 2
    ;;
esac
