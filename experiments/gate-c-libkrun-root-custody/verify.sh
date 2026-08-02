#!/bin/sh
set -eu

experiment_dir=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
repository_dir=$(CDPATH='' cd -- "$experiment_dir/../.." && pwd)

for script in "$experiment_dir"/*.sh; do
    sh -n "$script"
done
python3 -m py_compile "$experiment_dir/local_custody.py" "$experiment_dir/run_guest.py"
"$experiment_dir/source-audit.sh"
"$experiment_dir/build.sh"
python3 "$experiment_dir/local_custody.py" >/dev/null

if [ "${CAPSULE_RUN_GUEST:-false}" = true ]; then
    python3 "$experiment_dir/run_guest.py" --timeout 60
fi

git -C "$repository_dir" diff --check -- experiments/gate-c-libkrun-root-custody
printf 'P0-1 experiment verification passed\n'
