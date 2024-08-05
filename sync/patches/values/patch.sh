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
cp "${script_dir_rel}/Makefile.giantswarm" "./helm/Makefile.giantswarm"

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

# We also need to remove changes to Chart.yaml because it's maintained by us.
# The regexp makes sure the target starts with a letter so it doesn't match
# .PHONY.
sed -i 's/\(^[-_a-z]\+:.*\) update-chart/\1/g' ./helm/Makefile

# We also need to include our own Makefile.giantswarm into the main Makefile 
# to override image names
sed -i '/include $(MAKEFILE_VALUES)/a include Makefile.giantswarm' ./helm/Makefile

cd ./helm && make ; cd -

{ set +x; } 2>/dev/null
