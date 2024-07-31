#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
#script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

#readonly script_dir_rel=".${script_dir#"${repo_dir}"}"

set -x

# Apply patches to _helpers.tpl files. More info in README.md of this patch
# directory.
git apply ./sync/patches/image_registries/_helpers.tpl.patch
git apply ./sync/patches/image_registries/_cilium_operator__helpers.tpl.patch

{ set +x; } 2>/dev/null

echo

set -x

# This replaces lines like:
#
#    - image: {{ include "cilium.image" .Values.preflight.image | quote }}
#    + image: {{ include "cilium.image" (list $ .Values.preflight.image) | quote }}
#
# and:
#
#    - image: {{ include "cilium.operator.image" .Values.preflight.image | quote }}
#    + image: {{ include "cilium.oprator.image" (list $ .Values.preflight.image) | quote }}
#
find ./helm/cilium/templates -type f -name '*.yaml' -exec \
        sed -i 's/\({{ include "cilium[^"]*\.image"[[:space:]]\+\)\([^ ]*\)/\1(list $ \2)/' "{}" \;

{ set +x; } 2>/dev/null
