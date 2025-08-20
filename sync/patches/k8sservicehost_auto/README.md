# Patch for k8sServiceHost "auto" handling

Remove this once https://github.com/cilium/cilium/pull/41291 has been released

## How is the patch generated?

First, stage the changes (in `./helm`) and the run:

> [!TIP]
> Skip the `-R` flags if the changes were added.

```bash
git --no-pager diff -R helm/cilium/templates/_helpers.tpl \
        > sync/patches/k8sservicehost_auto/_helpers.tpl.patch
```

## What is the patched change?

In case something goes wrong this is the raw change, in file `helm/cilium/templates/_helpers.tpl`, replace:

```
  {{- if and (eq .Values.k8sServiceHost "auto") (lookup "v1" "ConfigMap" $configmapNamespace $configmapName) }}
    {{- $configmap := (lookup "v1" "ConfigMap" $configmapNamespace $configmapName) }}
    {{- $kubeconfig := get $configmap.data "kubeconfig" }}
    {{- $k8sServer := get ($kubeconfig | fromYaml) "clusters" | mustFirst | dig "cluster" "server" "" }}
    {{- $uri := (split "https://" $k8sServer)._1 | trim }}
    {{- (split ":" $uri)._0 | quote }}
```

with:

```
  {{- if eq .Values.k8sServiceHost "auto" }}
    {{- $configmap := (lookup "v1" "ConfigMap" $configmapNamespace $configmapName) }}
    {{- if $configmap }}
      {{- $kubeconfig := get $configmap.data "kubeconfig" }}
      {{- $k8sServer := get ($kubeconfig | fromYaml) "clusters" | mustFirst | dig "cluster" "server" "" }}
      {{- $uri := (split "https://" $k8sServer)._1 | trim }}
      {{- (split ":" $uri)._0 | quote }}
    {{- else }}
      {{- fail (printf "ConfigMap %s/%s not found, please create it or set k8sServiceHost to a valid value" $configmapNamespace $configmapName) }}
    {{- end }}
```
