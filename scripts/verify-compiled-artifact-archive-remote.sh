#!/bin/sh
set -eu

archive_repository="https://github.com/Shrimpworks/capsule-experiments.git"
archive_commit="0944ffd8cfd01ec23e4ae99138b0931d56804077"
archive_path="experiments/completed-compiled-artifact-payloads"
verification_root="$(mktemp -d "${TMPDIR:-/tmp}/capsule-artifact-archive.XXXXXX")"
script_root="$(
  unset CDPATH
  cd -- "$(dirname -- "$0")"
  pwd
)"

case "$verification_root" in
  "${TMPDIR:-/tmp}"/capsule-artifact-archive.??????) ;;
  *)
    printf '%s\n' "unexpected verification root: $verification_root" >&2
    exit 1
    ;;
esac

cleanup() {
  status=$?
  trap - EXIT
  rm -rf -- "$verification_root"
  exit "$status"
}
trap cleanup EXIT
trap 'exit 129' HUP
trap 'exit 130' INT
trap 'exit 143' TERM

git -C "$verification_root" init --quiet
git -C "$verification_root" remote add origin "$archive_repository"
git -C "$verification_root" fetch --quiet --depth=1 origin "$archive_commit"

fetched_commit="$(git -C "$verification_root" rev-parse FETCH_HEAD)"
if [ "$fetched_commit" != "$archive_commit" ]; then
  printf '%s\n' "archive commit mismatch: expected $archive_commit, got $fetched_commit" >&2
  exit 1
fi

git -C "$verification_root" checkout --quiet --detach FETCH_HEAD
# The checkout is immutable archive data. Execute only the verifier reviewed in this repository.
node "$script_root/verify-compiled-artifact-archive.mjs" "$verification_root/$archive_path"
