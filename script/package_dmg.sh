#!/usr/bin/env bash
set -euo pipefail

APP_NAME="ClawPet"
VERSION="${1:-1.0.0}"

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="$ROOT_DIR/dist"
APP_BUNDLE="$DIST_DIR/$APP_NAME.app"
STAGING_DIR="$DIST_DIR/dmg-staging"
DMG_PATH="$DIST_DIR/$APP_NAME-$VERSION.dmg"

"$ROOT_DIR/script/build_and_run.sh" --package

rm -rf "$STAGING_DIR" "$DMG_PATH"
mkdir -p "$STAGING_DIR"
cp -R "$APP_BUNDLE" "$STAGING_DIR/"
ln -s /Applications "$STAGING_DIR/Applications"

codesign --verify --deep --strict --verbose=2 "$APP_BUNDLE"
plutil -lint "$APP_BUNDLE/Contents/Info.plist"

hdiutil create \
  -volname "$APP_NAME $VERSION" \
  -srcfolder "$STAGING_DIR" \
  -ov \
  -format UDZO \
  "$DMG_PATH"

rm -rf "$STAGING_DIR"
hdiutil verify "$DMG_PATH"
echo "$DMG_PATH"
