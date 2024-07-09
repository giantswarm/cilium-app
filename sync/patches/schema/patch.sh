#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

readonly script_dir_rel=".${script_dir#"${repo_dir}"}"

set -x
cp "${script_dir_rel}/values.yaml.tmpl" "./vendor/cilium/install/kubernetes/cilium/values.yaml.tmpl"
cp "${script_dir_rel}/Makefile.values" "./vendor/cilium/install/kubernetes/Makefile.values"

cd ./vendor/cilium/install/kubernetes && make && cd -

HELM_TOOLBOX_VERSION="v1.1.0"
HELM_TOOLBOX_SHA="961693f182b9b456ed90e5274ac5df81e4af4343104e252666959cdf9570ce9e"
HELM_TOOLBOX_IMAGE="quay.io/cilium/helm-toolbox:${HELM_TOOLBOX_VERSION}@sha256:${HELM_TOOLBOX_SHA}"
cd ./vendor/cilium/install/kubernetes && docker container run --rm \
        --workdir /src/install/kubernetes \
        --volume "${PWD}/../..:/src" \
        --volume "${PWD}:${PWD}" \
        --user "$(id -u):$(id -g)" \
        "${HELM_TOOLBOX_IMAGE}" helm-schema -c cilium --skip-auto-generation title,description,required,default,additionalProperties
cd -

{ set +x; } 2>/dev/null
