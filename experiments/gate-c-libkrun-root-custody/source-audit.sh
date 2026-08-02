#!/bin/sh
set -eu

libkrun_dir=${CAPSULE_LIBKRUN_SOURCE:-/private/tmp/capsule-libkrun-v1.19.4}
imago_file=${CAPSULE_IMAGO_FILE_SOURCE:-/Users/dsteele/.cargo/registry/src/index.crates.io-1949cf8c6b5b557f/imago-0.2.3/src/file.rs}
expected_commit=728df8125077d0db44265f6e997c72b81b65c015
device_source="$libkrun_dir/src/devices/src/virtio/block/device.rs"
api_source="$libkrun_dir/src/libkrun/src/lib.rs"

actual_commit=$(git -C "$libkrun_dir" rev-parse HEAD)
if [ "$actual_commit" != "$expected_commit" ]; then
    printf 'unexpected libkrun commit: %s\n' "$actual_commit" >&2
    exit 2
fi
for required in "$device_source" "$api_source" "$imago_file"; do
    if [ ! -f "$required" ]; then
        printf 'missing source input: %s\n' "$required" >&2
        exit 2
    fi
done

first_open=$(grep -c '\.open(PathBuf::from(&disk_image_path))' "$device_source")
second_open=$(grep -c 'ImagoFile::open_sync(file_opts)' "$device_source")
stored_path=$(grep -c 'disk_image_path: disk_path.to_string()' "$api_source")
positional_reads=$(grep -c 'libc::preadv(' "$imago_file")
positional_writes=$(grep -c 'libc::pwritev(' "$imago_file")
if [ "$first_open" -ne 1 ] || [ "$second_open" -ne 1 ] || \
    [ "$stored_path" -lt 1 ] || [ "$positional_reads" -lt 1 ] || \
    [ "$positional_writes" -lt 1 ]; then
    printf 'source audit invariant failed\n' >&2
    exit 1
fi

printf 'libkrunCommit=%s\n' "$actual_commit"
printf 'blockConsumer.metadataOpen=%s\n' "$first_open"
printf 'blockConsumer.imagoOpen=%s\n' "$second_open"
printf 'apiStoresOnlySuppliedPath=%s\n' "$stored_path"
printf 'imagoPositionalReadSites=%s\n' "$positional_reads"
printf 'imagoPositionalWriteSites=%s\n' "$positional_writes"
printf 'deviceSourceSha256=%s\n' "$(shasum -a 256 "$device_source" | awk '{print $1}')"
printf 'apiSourceSha256=%s\n' "$(shasum -a 256 "$api_source" | awk '{print $1}')"
printf 'imagoSourceSha256=%s\n' "$(shasum -a 256 "$imago_file" | awk '{print $1}')"
