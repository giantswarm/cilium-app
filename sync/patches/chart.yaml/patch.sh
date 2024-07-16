#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
#script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

#readonly script_dir_rel=".${script_dir#"${repo_dir}"}"

set -x
sed -i \
        -e '/^version:/c\version: 0.0.0-dev' \
        -e '/^sources:/a\  - https://github.com/giantswarm/cilium-app' \
        -e '/^annotations:/a\  application.giantswarm.io/team: "cabbage"' \
        "./helm/cilium/Chart.yaml"
sed -i \
        -e '/^\* <https:\/\/github.com\/cilium\/cilium>/i\* <https://github.com/giantswarm/cilium-app>' \
        "./helm/cilium/README.md"
{ set +x; } 2>/dev/null
