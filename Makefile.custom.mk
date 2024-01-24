.PHONY: ensure-schema-gen
ensure-schema-gen:
	@helm schema-gen --help &>/dev/null || helm plugin install https://github.com/mihaisee/helm-schema-gen.git

.PHONY: schema-gen
schema-gen: ensure-schema-gen ## Generates the values schema file
	@cd helm/cilium && helm schema-gen values.yaml > values.schema.json
