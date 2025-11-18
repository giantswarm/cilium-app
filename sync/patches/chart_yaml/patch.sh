#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir

cd "${repo_dir}"

set -x
APP_VERSION=$(yq e '.directories[] | select(.path == "vendor").contents[] | select(.path == "cilium").git.ref' vendir.yml)

yq -i e ".appVersion |= \"${APP_VERSION#v}\"" ./helm/cilium/Chart.yaml

{ set +x; } 2>/dev/null
