{{- define "cilium.operator.cloud" -}}
{{- $cloud := "generic" -}}
{{- if .Values.eni.enabled -}}
  {{- $cloud = "aws" -}}
{{- else if .Values.azure.enabled -}}
  {{- $cloud = "azure" -}}
{{- else if .Values.alibabacloud.enabled -}}
  {{- $cloud = "alibabacloud" -}}
{{- end -}}
{{- $cloud -}}
{{- end -}}

{{- define "cilium.operator.imageDigestName" -}}
{{- $imageDigest := (.Values.operator.image.useDigest | default false) | ternary (printf "@%s" .Values.operator.image.genericDigest) "" -}}
{{- if .Values.eni.enabled -}}
  {{- $imageDigest = (.Values.operator.image.useDigest | default false) | ternary (printf "@%s" .Values.operator.image.awsDigest) "" -}}
{{- else if .Values.azure.enabled -}}
  {{- $imageDigest = (.Values.operator.image.useDigest | default false) | ternary (printf "@%s" .Values.operator.image.azureDigest) "" -}}
{{- else if .Values.alibabacloud.enabled -}}
  {{- $imageDigest = (.Values.operator.image.useDigest | default false) | ternary (printf "@%s" .Values.operator.image.alibabacloudDigest) "" -}}
{{- end -}}
{{- $imageDigest -}}
{{- end -}}

{{/*
Return cilium operator image
*/}}
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
