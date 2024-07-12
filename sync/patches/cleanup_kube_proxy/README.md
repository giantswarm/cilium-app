> [!CAUTION]
> When removing this patch make sure to also remove `cleanupKubeProxy: false` from `values.yaml`. It may be generated from `patches/schema` patch.

## How were patches generated?

First, stage the changes (in `./helm`) and the run:

> [!TIP]
> Skip the `-R` flags if the changes were added.

```bash
git --no-pager diff -R helm/cilium/templates/cilium-agent/daemonset.yaml \
        > sync/patches/cleanup_kube_proxy/cilium_agent__daemonset.yaml.patch
```

## What is the patched change?

> [!NOTE]
> There is also `cleanupKubeProxy:` value appended to `values.yaml`.

In case something goes wrong this is the raw change is:

In file `./helm/cilium/templates/cilium-agent/daemonset.yaml` after `wait-for-kube-proxy` container, add:

```
      {{ if and (.Values.cleanupKubeProxy) (not (eq .Values.kubeProxyReplacement "disabled")) }}
      - name: cleanup-kube-proxy-iptables
        image: "{{ include "cilium.image" (list $ .Values.image) }}"
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        securityContext:
          privileged: true
        command:
        - sh
        - -c
        - |
          /usr/sbin/iptables-nft-save | grep -v KUBE | grep -v cali | /usr/sbin/iptables-nft-restore
      {{ end }}
```
