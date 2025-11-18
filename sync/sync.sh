#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly dir
cd "${dir}/.."

# Stage 1 sync - intermediate to the ./vendir folder
set -x
vendir sync
helm dependency update helm/cilium/
{ set +x; } 2>/dev/null

# Patches
./sync/patches/eni/patch.sh
./sync/patches/image_registries/patch.sh
./sync/patches/readme/patch.sh
./sync/patches/networkpolicies/patch.sh
./sync/patches/k8sservicehost_auto/patch.sh
./sync/patches/chart_yaml/patch.sh
./sync/patches/readme/patch.sh
./sync/patches/values/patch.sh

# Store diffs
rm -f ./diffs/*
for f in $(git --no-pager diff --no-exit-code --no-color --no-index vendor/cilium/install/kubernetes helm --name-only) ; do
        [[ "$f" == "helm/cilium/Chart.yaml" ]] && continue
        [[ "$f" == "helm/cilium/Chart.lock" ]] && continue
        [[ "$f" == "helm/cilium/README.md" ]] && continue
        [[ "$f" == "helm/cilium/values.schema.json" ]] && continue
        [[ "$f" == "helm/cilium/values.yaml" ]] && continue
        [[ "$f" =~ ^helm/cilium/charts/.* ]] && continue

        base_file="vendor/cilium/install/kubernetes/${f#"helm/"}"
        [[ ! -e $base_file ]] && base_file="vendor/cilium/${f#"helm/"}"
        [[ ! -e $base_file ]] && base_file="/dev/null"

        set +e
        set -x
        git --no-pager diff --no-exit-code --no-color --no-index "$base_file" "${f}" \
                > "./diffs/${f//\//__}.patch" # ${f//\//__} replaces all "/" with "__"

        { set +x; } 2>/dev/null
        set -e
        ret=$?
        if [ $ret -ne 0 ] && [ $ret -ne 1 ] ; then
                exit $ret
        fi
done

# Print upstream changelog
awk '/^##/ { if (++count == 2) exit } count >= 1' vendor/cilium/CHANGELOG.md
