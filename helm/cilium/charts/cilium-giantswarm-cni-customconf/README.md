# cilium-giantswarm-cni-customconf

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square)

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| eni.excludeInterfaceTags | object | `{}` | The tags used to exclude interfaces from IP allocation. (spec.eni.exclude-interface-tags) See https://docs.cilium.io/en/stable/network/concepts/ipam/eni/#eni-allocation-parameters |
| eni.firstInterfaceIndex | int | `1` | The index of the first ENI to use for IP allocation. (spec.eni.first-interface-index) See https://docs.cilium.io/en/stable/network/concepts/ipam/eni/#eni-allocation-parameters |
| eni.securityGroupTags | object | `{}` | The list tags which will be used to filter the security groups to attach to any ENI that is created and attached to the instance. (spec.eni.security-group-tags) See https://docs.cilium.io/en/stable/network/concepts/ipam/eni/#eni-allocation-parameters |
| eni.subnetTags | object | `{}` | The tags used to select the AWS subnets for IP allocation. (spec.eni.subnet-tags) See https://docs.cilium.io/en/stable/network/concepts/ipam/eni/#eni-allocation-parameters |

