#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir

cd "${repo_dir}"

# Wire certgen's CA expiry early-warning flag into the certgen job specs.
# Upstream wires this on main, but not yet in the vendored chart version.
#
# certgen never rotates the CA on its own (it always runs with
# --ca-reuse-secret and only mints a CA when the cilium-ca secret is
# missing, see https://github.com/cilium/certgen/issues/500). With
# certgen.enforceCAValidityThroughoutLeavesDuration=true the certgen job
# instead fails as soon as the CA would expire before the requested leaf
# certificate validity ends — i.e. roughly one year before CA expiry with
# 365-day leafs — which pages via the CiliumHubbleCertificateRenewalJobFailed
# alert, leaving plenty of time for a manual CA rotation.
for f in \
        ./helm/cilium/templates/hubble/tls-cronjob/_job-spec.tpl \
        ./helm/cilium/templates/clustermesh-apiserver/tls-cronjob/_job-spec.tpl ; do
        if grep -q 'ca-enforce-validity-throughout-leaves-duration' "$f" ; then
                echo "Skipping $f: upstream already wires the flag, this patch can be removed"
                continue
        fi

        set -x
        sed -i 's/^\([[:space:]]*\)- "--ca-common-name=Cilium CA"$/&\n\1- "--ca-enforce-validity-throughout-leaves-duration={{ .Values.certgen.enforceCAValidityThroughoutLeavesDuration }}"/' "$f"
        { set +x; } 2>/dev/null

        grep -q 'ca-enforce-validity-throughout-leaves-duration' "$f" || {
                echo "Failed to patch ${f}: --ca-common-name anchor line not found" >&2
                exit 1
        }
done
