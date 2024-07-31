#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
script_dir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) ; readonly script_dir

cd "${repo_dir}"

readonly script_dir_rel=".${script_dir#"${repo_dir}"}"

set -x
git apply "${script_dir_rel}/cilium_agent__daemonset.yaml.patch"

#
# Because patches/schema/patch.sh overwrites values.yaml.tmpl, the changes below
# had to be added there. If `patches/schema` is removed, please uncomment the
# lines below.
#
## Add cleanupKubeProxy: to values.yaml
#cat <<'EOF' > "./helm/cilium/values.yaml.tmpl"
## If true, it adds an initContainer to cilium-agent pods that cleans up any legacy kube-proxy iptables rules from the node before running cilium.
## Only makes sense when `kubeProxyReplacement` is enabled (i.e. not set to 'disabled').
#cleanupKubeProxy: false
#EOF
{ set +x; } 2>/dev/null
