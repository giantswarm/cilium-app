If there is an regression the patches can be created like this:

```
$ git --no-pager diff --no-color -R vendor/cilium/install/kubernetes/cilium/templates/_helpers.tpl > ./customization/patches/image_registries/_helpers.tpl.patch

$ git --no-pager diff --no-color patch -R vendor/cilium/install/kubernetes/cilium/templates/cilium-operator/_helpers.tpl > ./customization/patches/image_registries/_cilium_operator__helpers.tpl.patch
```

Just in case this is the desired definition of `cilium.image` that should be present in `_helpers.tpl` and `cilium-operator/_helpers.tpl`:

```
{{- define "cilium.image" -}}
{{- if not (kindIs "slice" .) }}
{{- (fail (printf "required list, but got %q" (kindOf .))) }}
{{- end }}
{{- if (ne (len .) 2) }}
{{- (fail (printf "required list of 2 arguments, but got %d" (len .))) }}
{{- end }}
{{- $ := index . 0 }}
{{- with index . 1 }}
{{- $digest := (.useDigest | default false) | ternary (printf "@%s" .digest) "" -}}
{{- if .override -}}
{{- printf "%s" .override -}}
{{- else -}}
{{- printf "%s/%s:%s%s" $.Values.image.registry .repository .tag $digest -}}
{{- end -}}
{{- end -}}
{{- end -}}
```
