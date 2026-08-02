#!/bin/sh
set -eu

experiment_dir=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
source_dir=${CAPSULE_LIBKRUN_SOURCE:-/private/tmp/capsule-libkrun-v1.19.4}
expected_commit=728df8125077d0db44265f6e997c72b81b65c015
patch_file="$experiment_dir/patches/0001-console-correctness.patch"

if [ ! -d "$source_dir/.git" ]; then
    printf 'missing retained local libkrun checkout: %s\n' "$source_dir" >&2
    exit 2
fi

actual_commit=$(git -C "$source_dir" rev-parse HEAD)
if [ "$actual_commit" != "$expected_commit" ]; then
    printf 'unexpected libkrun commit: got %s, want %s\n' \
        "$actual_commit" "$expected_commit" >&2
    exit 2
fi

task_tmp=$(mktemp -d /private/tmp/capsule-console-correctness.XXXXXX)
trap 'rm -rf "$task_tmp"' EXIT HUP INT TERM
source_copy="$task_tmp/libkrun"
mkdir -p "$source_copy"

git -C "$source_dir" archive "$expected_commit" | tar -x -C "$source_copy"
patch -d "$source_copy" -p1 --batch --forward <"$patch_file"

rustfmt --edition 2021 --check \
    "$source_copy/src/devices/src/virtio/console/device.rs" \
    "$source_copy/src/devices/src/virtio/console/port.rs" \
    "$source_copy/src/devices/src/virtio/console/port_io.rs" \
    "$source_copy/src/devices/src/virtio/console/process_tx.rs"

(
    cd "$source_copy"
    CARGO_NET_OFFLINE=true \
        CARGO_TARGET_DIR="$task_tmp/target" \
        cargo test --offline -p krun-devices --lib
)

printf 'commit=%s\n' "$expected_commit"
printf 'patchSha256=%s\n' "$(shasum -a 256 "$patch_file" | awk '{print $1}')"
printf 'result=PASS\n'
