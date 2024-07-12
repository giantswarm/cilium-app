#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

readonly script_dir_rel=".${script_dir#"${repo_dir}"}"

set -x
cp "${script_dir_rel}/values.yaml.tmpl" "./helm/cilium/values.yaml.tmpl"
cp "${script_dir_rel}/Makefile.values" "./helm/Makefile.values"

# We need to copy some Makefile dependencies and edit the Makefile (removing
# "../.."), because we don't have that deep file structure inside ./helm
# directory.
cp ./vendor/cilium/{Makefile.defs,Makefile.quiet,VERSION} ./helm/
# replace "include ../../Makefile.defs" with "include Makefile.defs"
sed -i 's#include \.\./\.\./Makefile.defs#include Makefile.defs#g' ./helm/Makefile
# replace "--workdir /src/install/kubernetes" with "--workdir /src"
sed -i 's#--workdir /src/install/kubernetes#--workdir /src#g' ./helm/Makefile
# replace "--volume $(CURDIR)/../..:/src" with "--volume $(CURDIR):/src"
sed -i 's#\(--volume .*\)\.\./\.\.:/src#\1:/src#g' ./helm/Makefile

cd ./helm && make ; cd -

HELM_TOOLBOX_VERSION="v1.1.0"
HELM_TOOLBOX_SHA="961693f182b9b456ed90e5274ac5df81e4af4343104e252666959cdf9570ce9e"
HELM_TOOLBOX_IMAGE="quay.io/cilium/helm-toolbox:${HELM_TOOLBOX_VERSION}@sha256:${HELM_TOOLBOX_SHA}"
cd ./helm && docker container run --rm \
        --workdir /src \
        --volume "${PWD}:/src" \
        --user "$(id -u):$(id -g)" \
        "${HELM_TOOLBOX_IMAGE}" helm-schema -c cilium --skip-auto-generation title,description,required,default,additionalProperties
cd -

{ set +x; } 2>/dev/null
