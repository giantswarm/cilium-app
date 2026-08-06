#!/usr/bin/env bash
#
# Assert that every container image the chart can render is addressed through
# the Giant Swarm registry.
#
# We used to guarantee this by patching the chart's image helpers to prefix
# .Values.image.registry onto every image (sync/patches/image_registries, gone
# since we dropped that value). The registry now lives in the repository values
# in sync/patches/values/Makefile.giantswarm, exactly as upstream does it. That
# is far less invasive, but it is also silent: when upstream adds or renames an
# image, the new repository simply arrives pointing at quay.io and nothing
# complains until a pod fails to pull. This script is that missing complaint.
#
# It checks two things, because a registry check that renders nothing passes
# just as quietly as one that renders everything:
#
#   1. every rendered image is served from gsoci, or is explicitly listed in
#      sync/unmirrored-images.txt;
#   2. every repository we configure in Makefile.giantswarm (and every
#      allowlisted one) actually shows up in the render, so a scenario that
#      stops covering an image fails instead of narrowing the check.

set -o errexit
set -o nounset
set -o pipefail

repo_dir=$(git rev-parse --show-toplevel) ; readonly repo_dir
cd "${repo_dir}"

readonly chart="./helm/cilium"
readonly required_prefix="gsoci.azurecr.io/giantswarm/"
readonly unmirrored_file="./sync/unmirrored-images.txt"
readonly repos_file="./sync/patches/values/Makefile.giantswarm"

tmp=$(mktemp -d) ; readonly tmp
trap 'rm -rf "${tmp}"' EXIT

# Render the chart several times rather than once: some images sit behind
# mutually exclusive options (encryption.type, tls.auto.method), so no single
# set of values reaches all of them. The union must cover every image.
#
# Keep each scenario's name and flags together; render() aborts on any failure
# so a scenario that stops rendering can never silently shrink the check.
render() {
        local name="$1" ; shift
        if ! helm template cilium "${chart}" "$@" > "${tmp}/${name}.yaml" 2> "${tmp}/${name}.err" ; then
                {
                        echo "verify-images: FAILED -- scenario '${name}' did not render:"
                        cat "${tmp}/${name}.err"
                        echo
                        echo "Fix the scenario's values in ${BASH_SOURCE[0]}; do not drop it."
                } >&2
                exit 1
        fi
}

# Everything that is merely disabled by default.
render optional \
        --set nodeinit.enabled=true \
        --set preflight.enabled=true \
        --set envoy.enabled=true \
        --set hubble.relay.enabled=true \
        --set hubble.ui.enabled=true \
        --set clustermesh.useAPIServer=true \
        --set standaloneDnsProxy.enabled=true \
        --set dnsProxy.proxyPort=10094 \
        --set authentication.mutual.spire.enabled=true \
        --set authentication.mutual.spire.install.enabled=true

# certgen, via the cronJob TLS method on both hubble and clustermesh.
render certgen \
        --set hubble.enabled=true \
        --set hubble.tls.auto.enabled=true \
        --set hubble.tls.auto.method=cronJob \
        --set clustermesh.useAPIServer=true \
        --set clustermesh.apiserver.tls.auto.enabled=true \
        --set clustermesh.apiserver.tls.auto.method=cronJob

# ztunnel, which needs encryption.type set away from its default.
render ztunnel \
        --set encryption.enabled=true \
        --set encryption.type=ztunnel

# The MCS-API CoreDNS autoconfigure job.
render mcsapi \
        --set clustermesh.useAPIServer=true \
        --set clustermesh.mcsapi.enabled=true \
        --set clustermesh.mcsapi.corednsAutoConfigure.enabled=true

# Strip a tag and/or digest to leave the bare repository.
image_repository() {
        local ref="$1"
        ref="${ref%%@*}"
        # Only strip a trailing :tag, never a registry port.
        if [[ "${ref##*/}" == *:* ]] ; then
                ref="${ref%:*}"
        fi
        printf '%s' "$ref"
}

# Allowlisted repositories, one per line, comments and blanks stripped.
unmirrored=""
if [[ -f "$unmirrored_file" ]] ; then
        unmirrored=$(grep -v '^[[:space:]]*#' "$unmirrored_file" | tr -d '[:blank:]' | grep -v '^$' || true)
fi

# Every distinct image reference across all renders.
refs=$(cat "${tmp}"/*.yaml \
        | grep -E '^[[:space:]]*image:[[:space:]]' \
        | tr -d '"' \
        | awk '{print $2}' \
        | sort -u)

if [[ -z "$refs" ]] ; then
        echo "verify-images: rendered no image references at all -- the render scenarios are broken" >&2
        exit 1
fi

count=0
offenders=""
tolerated=""
seen_repos=""

while IFS= read -r ref ; do
        [[ -n "$ref" ]] || continue
        count=$((count + 1))
        repo=$(image_repository "$ref")

        # standaloneDnsProxy ships with an empty repository upstream; it has no
        # Giant Swarm image and is unusable until the user sets one.
        if [[ -z "$repo" ]] ; then
                tolerated="${tolerated}(unset repository) ${ref}"$'\n'
                continue
        fi

        seen_repos="${seen_repos}${repo}"$'\n'

        if [[ "$repo" == "${required_prefix}"* ]] ; then
                continue
        fi

        if printf '%s\n' "$unmirrored" | grep -Fxq -- "$repo" ; then
                tolerated="${tolerated}${ref}"$'\n'
                continue
        fi

        offenders="${offenders}${ref}"$'\n'
done <<< "$refs"

echo "verify-images: checked ${count} distinct image references across 4 render scenarios"

if [[ -n "$tolerated" ]] ; then
        echo "verify-images: not served from ${required_prefix}, allowed by ${unmirrored_file}:"
        printf '%s' "$tolerated" | while IFS= read -r l ; do echo "  $l" ; done
fi

if [[ -n "$offenders" ]] ; then
        {
                echo
                echo "verify-images: FAILED -- these images are not addressed through ${required_prefix}:"
                printf '%s' "$offenders" | while IFS= read -r l ; do echo "  $l" ; done
                echo
                echo "Upstream probably added or renamed an image. Either:"
                echo "  - mirror it (giantswarm/retagger) and set its *_REPO in"
                echo "    ${repos_file} to ${required_prefix}<name>, or"
                echo "  - add its repository to ${unmirrored_file} to pull it from upstream on purpose."
        } >&2
        exit 1
fi

# Coverage: every repository we configure must appear in the render. Guards
# against a scenario quietly stopping to cover an image, which would turn the
# check above into a no-op for it.
#
# Matched as a prefix because the operator repository gains a cloud suffix
# (cilium-operator -> cilium-operator-generic).
expected=$( { grep -hoE '^export [A-Z_]+_REPO[[:space:]]*:?=[[:space:]]*\S+' "$repos_file" \
                | awk -F'=' '{print $NF}' ; printf '%s\n' "$unmirrored" ; } \
        | tr -d '[:blank:]' | grep -v '^$' | sort -u)

missing=""
while IFS= read -r want ; do
        [[ -n "$want" ]] || continue
        printf '%s' "$seen_repos" | grep -q "^${want}" || missing="${missing}${want}"$'\n'
done <<< "$expected"

if [[ -n "$missing" ]] ; then
        {
                echo
                echo "verify-images: FAILED -- configured images never appeared in any render:"
                printf '%s' "$missing" | while IFS= read -r l ; do echo "  $l" ; done
                echo
                echo "Either a render scenario no longer reaches this image (fix the scenario"
                echo "in ${BASH_SOURCE[0]}), or the image is gone from the chart (drop its"
                echo "entry from ${repos_file} / ${unmirrored_file})."
        } >&2
        exit 1
fi

echo "verify-images: OK -- all configured images covered and served from ${required_prefix} unless allowlisted"
