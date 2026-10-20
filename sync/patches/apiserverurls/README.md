## Why?

`k8s.apiServerURLs: auto` derives the list of API server URLs at install time. It lists the
API server from `k8sServiceHost`/`k8sServicePort` (resolved from `kube-public/cluster-info` when
`auto`), followed by the ClusterIP of the `default/kubernetes` Service. With two or more URLs,
Cilium enables its API server failover mode: agents rotate through the list until their first
sync and then reach the API through the `kubernetes` Service, so they no longer depend on the
API load balancer. The ClusterIP is also reachable while an agent restarts, because the
datapath keeps translating it. Control plane node IPs are not listed: they change on every
control plane rollout, and Flux does not render the chart again when they do.

## How were patches generated?

Apply the change in `./helm` on top of the `eni` patches and run:

```bash
git --no-pager diff helm/cilium/templates/_helpers.tpl \
        > sync/patches/apiserverurls/_helpers.tpl.patch
git --no-pager diff helm/cilium/templates/cilium-configmap.yaml \
        > sync/patches/apiserverurls/cilium-configmap.yaml.patch
git --no-pager diff helm/cilium/templates/cilium-agent/daemonset.yaml \
        > sync/patches/apiserverurls/cilium_agent__daemonset.yaml.patch
```

## What is the patched change?

- `_helpers.tpl`: add the `k8sAPIServerURLs` helper.
- `cilium-configmap.yaml` and `cilium-agent/daemonset.yaml`: render
  `{{ include "k8sAPIServerURLs" . }}` instead of `{{ .Values.k8s.apiServerURLs }}`.
