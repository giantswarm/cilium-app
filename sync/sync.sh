#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly dir
cd "${dir}/.."

# Stage 1 sync - intermediate to the ./vendir folder
set -x
vendir sync --file vendir.s1.yml --lock-file vendir.s1.lock
{ set +x; } 2>/dev/null

# Copy cilium original files so we can do a diff later on.
cp -a ./vendor/cilium ./vendor/cilium.orig

# Patches
./sync/patches/cleanup_kube_proxy/patch.sh
./sync/patches/schema/patch.sh
./sync/patches/eni/patch.sh
./sync/patches/image_registries/patch.sh
./sync/patches/metrics_port/patch.sh
./sync/patches/chart.yaml/patch.sh

# Store diffs
rm -f ./diffs/*
for f in $(git -C ./vendor --no-pager diff --no-exit-code --no-color --no-index cilium{.orig,} --name-only) ; do
        set +e
        set -x
        git -C ./vendor --no-pager diff --no-exit-code --no-color --no-index "cilium.orig/${f#"cilium/"}" "${f}" \
                > "./diffs/${f//\//__}.patch" # ${f//\//__} replaces all "/" with "__"
        ret=$?
        { set +x; } 2>/dev/null
        set -e
        if [ $ret -ne 0 ] && [ $ret -ne 1 ] ; then
                exit $ret
        fi
done

# Stage 2 sync - from ./vendir to ./helm
set -x
vendir sync --file vendir.s2.yml --lock-file vendir.s2.lock
{ set +x; } 2>/dev/null
