#!/bin/sh
set -eu

experiment_dir=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
libkrun_dir=${CAPSULE_LIBKRUN_SOURCE:-/private/tmp/capsule-libkrun-v1.19.4}
libkrunfw=${CAPSULE_LIBKRUNFW_LIBRARY:-/private/tmp/capsule-libkrunfw-v5.5.0/libkrunfw/libkrunfw.5.dylib}
build_dir="$experiment_dir/.build"
library_source="$libkrun_dir/target/release/libkrun.1.19.4.dylib"
app="$build_dir/CapsuleRootCustodySpike.app"

for required in "$libkrun_dir/include/libkrun.h" "$library_source" "$libkrunfw"; do
    if [ ! -f "$required" ]; then
        printf 'missing pinned build input: %s\n' "$required" >&2
        exit 2
    fi
done
if [ "$(git -C "$libkrun_dir" rev-parse HEAD)" != \
    728df8125077d0db44265f6e997c72b81b65c015 ]; then
    printf 'unexpected libkrun source revision\n' >&2
    exit 2
fi
if ! grep -q 'const KRUNFW_NAME: &str = "@rpath/libkrunfw.5.dylib";' \
    "$libkrun_dir/src/libkrun/src/lib.rs" || \
    ! grep -q 'strcmp(krun_root_options, "ro,nosuid,nodev") == 0' \
    "$libkrun_dir/src/init_blob/init/init.c"; then
    printf 'pinned libkrun tree is missing the two retained Gate C patches\n' >&2
    exit 2
fi

mkdir -p "$build_dir/lib" "$app/Contents/MacOS/lib"
cp "$library_source" "$build_dir/lib/libkrun.1.19.4.dylib"
cp "$libkrunfw" "$build_dir/lib/libkrunfw.5.dylib"
ln -sfn libkrun.1.19.4.dylib "$build_dir/lib/libkrun.1.dylib"
ln -sfn libkrun.1.dylib "$build_dir/lib/libkrun.dylib"

clang -std=c17 -Wall -Wextra -Werror -Wno-deprecated-declarations \
    -isystem "$libkrun_dir/include" "$experiment_dir/src/runner.c" \
    -L "$build_dir/lib" -lkrun -Wl,-rpath,@executable_path/lib \
    -o "$build_dir/capsule-root-custody-runner"
clang -std=c17 -Wall -Wextra -Werror -Wno-deprecated-declarations -dynamiclib \
    "$experiment_dir/src/open_trace.c" -o "$build_dir/libcapsule-open-trace.dylib"

install_name_tool -id @rpath/libkrun.1.dylib \
    "$build_dir/lib/libkrun.1.19.4.dylib"
install_name_tool -id @rpath/libkrunfw.5.dylib \
    "$build_dir/lib/libkrunfw.5.dylib"
install_name_tool -change libkrun.1.dylib @rpath/libkrun.1.dylib \
    "$build_dir/capsule-root-custody-runner"

cp "$build_dir/capsule-root-custody-runner" \
    "$app/Contents/MacOS/capsule-root-custody-runner"
cp "$build_dir/lib/libkrun.1.19.4.dylib" "$app/Contents/MacOS/lib/"
cp "$build_dir/lib/libkrunfw.5.dylib" "$app/Contents/MacOS/lib/"
ln -sfn libkrun.1.19.4.dylib "$app/Contents/MacOS/lib/libkrun.1.dylib"
ln -sfn libkrun.1.dylib "$app/Contents/MacOS/lib/libkrun.dylib"
cp "$experiment_dir/Info.plist" "$app/Contents/Info.plist"

codesign --force --sign - "$build_dir/lib/libkrunfw.5.dylib"
codesign --force --sign - "$build_dir/lib/libkrun.1.19.4.dylib"
codesign --force --sign - "$build_dir/libcapsule-open-trace.dylib"
codesign --force --sign - --identifier com.capsulecorp.spike.root-custody-runner \
    --entitlements "$experiment_dir/runner-standalone.entitlements" \
    "$build_dir/capsule-root-custody-runner"
codesign --force --sign - "$app/Contents/MacOS/lib/libkrunfw.5.dylib"
codesign --force --sign - "$app/Contents/MacOS/lib/libkrun.1.19.4.dylib"
codesign --force --sign - --options runtime \
    --identifier com.capsulecorp.spike.root-custody-runner \
    --entitlements "$experiment_dir/runner.entitlements" \
    "$app/Contents/MacOS/capsule-root-custody-runner"
codesign --force --sign - --options runtime \
    --entitlements "$experiment_dir/runner.entitlements" "$app"

go_cache=${CAPSULE_GO_CACHE:-/private/tmp/capsule-root-custody-go-cache}
mkdir -p "$go_cache"
(
    cd "$experiment_dir/guest-probe"
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GOCACHE="$go_cache" go build \
        -trimpath -ldflags='-s -w -buildid=' \
        -o "$build_dir/guest-root-digest-linux-arm64" .
)

codesign --verify --strict --verbose=2 "$build_dir/capsule-root-custody-runner"
codesign --verify --deep --strict --verbose=2 "$app"
printf 'runner=%s\n' "$build_dir/capsule-root-custody-runner"
printf 'sandboxedRunner=%s\n' "$app/Contents/MacOS/capsule-root-custody-runner"
printf 'guestProbeSha256=%s\n' \
    "$(shasum -a 256 "$build_dir/guest-root-digest-linux-arm64" | awk '{print $1}')"
