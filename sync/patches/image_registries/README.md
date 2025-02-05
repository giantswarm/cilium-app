If there is an regression the patches can be created like this:

```
$ git --no-pager diff --no-color -R helm/cilium/templates/_helpers.tpl > ./sync/patches/image_registries/_helpers.tpl.patch

$ git --no-pager diff --no-color patch -R helm/cilium/templates/cilium-operator/_helpers.tpl > ./sync/patches/image_registries/_cilium_operator__helpers.tpl.patch
```

Just in case this is the desired definition of `cilium.image` that should be present in `_helpers.tpl`:

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
{{- $tag := .tag | default "" | eq "" | ternary "" (printf ":%s" .tag) -}}
{{- if .override -}}
{{- printf "%s" .override -}}
{{- else -}}
{{- printf "%s/%s%s%s" $.Values.image.registry .repository $tag $digest -}}
{{- end -}}
{{- end -}}
{{- end -}}
```

Helper for `cilium-operator/_helpers.tpl`:

```
{{- define "cilium.operator.image" -}}
{{- if not (kindIs "slice" .) }}
{{- (fail (printf "required list, but got %q" (kindOf .))) }}
{{- end }}
{{- if (ne (len .) 2) }}
{{- (fail (printf "required list of 2 arguments, but got %d" (len .))) }}
{{- end }}
{{- $ := index . 0 }}
{{- with index . 1 }}
{{- $cloud := include "cilium.operator.cloud" . }}
{{- $imageDigest := include "cilium.operator.imageDigestName" . }}
{{- $tag := .Values.operator.image.tag | default "" | eq "" | ternary "" (printf ":%s" .Values.operator.image.tag) }}
{{- if .Values.operator.image.override -}}
{{- printf "%s" .Values.operator.image.override -}}
{{- else -}}
{{- printf "%s/%s-%s%s%s%s" $.Values.image.registry .Values.operator.image.repository $cloud .Values.operator.image.suffix $tag $imageDigest -}}
{{- end -}}
{{- end -}}
{{- end -}}
```

Also all calls to `include "cilium.image"` should be replaced (done with `sed` in `patch.sh`) like so:

```diff
-          image: {{ include "cilium.image" .Values.preflight.image | quote }}
+          image: {{ include "cilium.image" (list $ .Values.preflight.image) | quote }}
```
