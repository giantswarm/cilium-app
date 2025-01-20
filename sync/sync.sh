#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly dir
cd "${dir}/.."

# Stage 1 sync - intermediate to the ./vendir folder
set -x
vendir sync
{ set +x; } 2>/dev/null

# Patches
./sync/patches/cleanup_kube_proxy/patch.sh
./sync/patches/eni/patch.sh
./sync/patches/image_registries/patch.sh
./sync/patches/metrics_port/patch.sh
./sync/patches/readme/patch.sh
./sync/patches/values/patch.sh

# Store diffs
rm -f ./diffs/*
for f in $(git --no-pager diff --no-exit-code --no-color --no-index vendor/cilium/install/kubernetes helm --name-only) ; do
        [[ "$f" == "helm/cilium/Chart.yaml" ]] && continue
        [[ "$f" == "helm/cilium/README.md" ]] && continue
        [[ "$f" == "helm/cilium/values.schema.json" ]] && continue
        [[ "$f" == "helm/cilium/values.yaml" ]] && continue
        set +e
        set -x
        git --no-pager diff --no-exit-code --no-color --no-index "vendor/cilium/install/kubernetes/${f#"helm/"}" "${f}" \
                > "./diffs/${f//\//__}.patch" # ${f//\//__} replaces all "/" with "__"
        ret=$?
        { set +x; } 2>/dev/null
        set -e
        if [ $ret -ne 0 ] && [ $ret -ne 1 ] ; then
                exit $ret
        fi
done
