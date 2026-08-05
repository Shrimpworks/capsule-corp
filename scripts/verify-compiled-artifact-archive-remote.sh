#!/bin/sh
set -eu

archive_repository="https://github.com/Shrimpworks/capsule-experiments.git"
archive_commit="0944ffd8cfd01ec23e4ae99138b0931d56804077"
archive_path="experiments/completed-compiled-artifact-payloads"
verification_root="$(mktemp -d "${TMPDIR:-/tmp}/capsule-artifact-archive.XXXXXX")"

git -C "$verification_root" init --quiet
git -C "$verification_root" remote add origin "$archive_repository"
git -C "$verification_root" fetch --quiet --depth=1 origin "$archive_commit"

fetched_commit="$(git -C "$verification_root" rev-parse FETCH_HEAD)"
if [ "$fetched_commit" != "$archive_commit" ]; then
  printf '%s\n' "archive commit mismatch: expected $archive_commit, got $fetched_commit" >&2
  exit 1
fi

git -C "$verification_root" checkout --quiet --detach FETCH_HEAD
node "$verification_root/$archive_path/scripts/verify.mjs"
