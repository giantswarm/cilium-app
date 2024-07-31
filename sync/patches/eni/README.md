## How were patches generated?

First, stage the changes (in `./helm`) and the run:

> [!TIP]
> Skip the `-R` flags if the changes were added.

```bash
git --no-pager diff -R helm/cilium/templates/cilium-agent/daemonset.yaml \
        > sync/patches/eni/cilium_agent__daemonset.yaml.patch
git --no-pager diff -R helm/cilium/templates/cilium-configmap.yaml \
        > sync/patches/eni/cilium-configmap.yaml.patch
```

## What is the patched change?

In case something goes wrong this is the raw change:


In file `./helm/cilium/templates/cilium-agent/daemonset.yaml` add the env vars below to `cilium-agent` and `config` containers:

```
        - name: CILIUM_CNI_CHAINING_MODE
          valueFrom:
            configMapKeyRef:
              name: cilium-config
              key: cni-chaining-mode
              optional: true
        - name: CILIUM_CUSTOM_CNI_CONF
          valueFrom:
            configMapKeyRef:
              name: cilium-config
              key: custom-cni-conf
              optional: true
```


In file `./helm/cilium/templates/cilium-configmap.yaml` replace:

```
{{- if .Values.cni.customConf  }}
  # legacy: v1.13 and before needed cni.customConf: true with cni.configMap
  write-cni-conf-when-ready: {{ .Values.cni.hostConfDirMountPath }}/05-cilium.conflist
{{- end }}
```

with:

```
  write-cni-conf-when-ready: {{ .Values.cni.hostConfDirMountPath }}/21-cilium.conflist
```
