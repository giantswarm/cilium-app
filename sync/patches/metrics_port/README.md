## How were patches generated?

First, stage the changes (in `./helm`) and the run:

> [!TIP]
> Skip the `-R` flags if the changes were added.

```bash
git --no-pager diff -R helm/cilium/templates/cilium-agent/service.yaml \
        > sync/patches/metrics_port/cilium_agent__service.yaml.patch
```

## What is the patched change?

In case something goes wrong this is the raw change:

In file `./helm/cilium/templates/cilium-agent/service.yaml` this port:

```
  - name: metrics
    port: {{ .Values.prometheus.port }}
    protocol: TCP
    targetPort: prometheus
```