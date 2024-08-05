#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
#script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

#readonly script_dir_rel=".${script_dir#"${repo_dir}"}"

set -x

# Apply patches to README.md.gotmpl files. More info in README.md of this patch
# directory.
git apply ./sync/patches/readme/_readme.md.gotmpl.patch

{ set +x; } 2>/dev/null

