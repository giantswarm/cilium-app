#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir

cd "${repo_dir}"

set -x
APP_VERSION=$(cat ./helm/VERSION) 

export APP_VERSION
yq -i e  '.appVersion |= env(APP_VERSION)' ./helm/cilium/Chart.yaml

{ set +x; } 2>/dev/null
