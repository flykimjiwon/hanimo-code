#!/usr/bin/env bash
set -e

echo "=== modol single binary builder ==="

# 1. Bundle with esbuild
echo "[1/3] Bundling with esbuild..."
node scripts/build-binary.mjs

# 2. Generate SEA blob
echo "[2/3] Generating SEA blob..."
node --experimental-sea-config sea.config.json

# 3. Copy node binary and inject blob
echo "[3/3] Creating binary..."
cp "$(which node)" dist/modol
npx postject dist/modol NODE_SEA_BLOB dist/sea-prep.blob \
  --sentinel-fuse NODE_SEA_FUSE_fce680ab2cc467b6e072b8b5df1996b2

# macOS: re-sign the binary
if [[ "$(uname)" == "Darwin" ]]; then
  echo "Signing binary (macOS)..."
  codesign --sign - dist/modol
fi

echo ""
echo "Binary built: dist/modol"
ls -lh dist/modol
