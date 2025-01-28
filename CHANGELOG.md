# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.31.0] - 2025-01-28

### Changed

- Upgrade Cilium to [v1.16.6](https://github.com/cilium/cilium/releases/tag/v1.16.6).

## [0.30.0] - 2025-01-14

### Changed

- Upgrade Cilium to [v1.16.5](https://github.com/cilium/cilium/releases/tag/v1.16.5).
- Revert back to upstream default of using image digests for container images.

### Removed

- Delete defaultPolicies and extraPolicies templates.

## [0.29.0] - 2024-11-12

### Changed

- Upgrade Cilium to [v1.16.3](https://github.com/cilium/cilium/releases/tag/v1.16.3).

## [0.28.0] - 2024-09-09

### Changed

- Move provider specific custom CNI configuration to subchart.

## [0.27.0] - 2024-08-20

### Changed

- Upgrade Cilium to [v1.16.1](https://github.com/cilium/cilium/releases/tag/v1.16.1).

## [0.26.0] - 2024-08-08

### Changed

- Upgrade Cilium to [v1.16.0](https://github.com/cilium/cilium/releases/tag/v1.16.0).
- Disable digest in all images.
- Improve security defaults for:
  - Hubble UI
  - Hubble Relay
  - Cilium Operator

## [0.25.2] - 2025-01-22

### Changed

- Upgrade cilium to [v1.15.13](https://github.com/cilium/cilium/releases/tag/v1.15.13).

## [0.25.1] - 2024-06-26

### Changed

- Fix regression setting Policy BPF Max map `policyMapMax` back to 65536 from 16384.

## [0.25.0] - 2024-06-19

### Changed

- Upgrade cilium to [v1.15.6](https://github.com/cilium/cilium/releases/tag/v1.15.6).

## [0.24.0] - 2024-04-30

### Changed

- Upgrade cilium to [v1.15.4](https://github.com/cilium/cilium/releases/tag/v1.15.4).

## [0.23.0] - 2024-04-24

### Added

- Cilium ENI mode for CAPA becomes usable with these changes

  - Add security group tag filter for pod network
  - Select subnets from secondary VPC CIDRs

## [0.22.0] - 2024-03-27

### Added

- Add helm values schema.

### Changed

- Add safe-to-evict annotations to Hubble Relay and UI pods.
- Enable deletion of extra network policies.
- Update team label to `cabbage`

## [0.21.0] - 2024-02-29

### Added

- Support CAPA clusters for ENI mode

### Changed

- Use SocketLB on host namespace only.

## [0.20.1] - 2024-02-27

### Changed

- Revert replacing `null` values.

## [0.20.0] - 2024-02-26

### Changed

- Upgrade cilium to `1.15.1`.
- Replace `null` values in `values.yaml` with its actual defaults. Config values with `null` types in the values schema prevented users from changing its values.

## [0.19.2] - 2024-01-22

### Fixed

- Replace `ToServices`/`ToPorts` combination in CiliumNetworkPolicy because of breakage in Cilium v1.14

## [0.19.1] - 2024-01-18

### Changed

- Set container registry to `gsoci.azurecr.io` in values.yaml.

## [0.19.0] - 2024-01-17

### Changed

- Upgrade cilium to `1.14.5`.
- Set default image registry to `gsoci.azurecr.io` in values.yaml.

## [0.18.0] - 2023-11-20

### Changed

- Upgrade cilium to `1.14.3`.

## [0.17.0] - 2023-11-08

### Changed

- Generate cilium chart using our fork and `vendir`.

## [0.16.0] - 2023-10-25

### Changed

- Disable uninstalling the CNI config files and binary when restarting the agent.

## [0.15.0] - 2023-10-24

### Added
- Add EKS support for cilium in ENI mode.

## [0.14.0] - 2023-10-18

### Changed

- Replace condition for PSP CR installation.

## [0.13.0] - 2023-09-26

### Added

- Support removal of previously-deployed default policies by setting `defaultPolicies.enabled=false` and `defaultPolicies.remove=false`

## [0.12.0] - 2023-09-05

### Added

- Support creating `CiliumNetworkPolicy` manifests that allow egress requests to DNS and proxy hosts

### Changed

- Add missing conditional for PSP rendering of default-policies installer job

## [0.11.2] - 2023-09-04

### Fixed

- Reenable BPF metrics

## [0.11.1] - 2023-09-01

### Changed

- Create custom CNI config depending on provider to allow bigger customization.
- Bump all manifests to upstream version 1.13.6.

## [0.11.0] - 2023-07-10

### Changed

- Increased Policy BPF Max map to 65536 from 16384.
- Enabled cilium_bpf_map_pressure metric.
- Excluding PSS labels from cilium identities/policies.
- Excluding Flux labels from cilium identities/policies.
- Excluding Helm labels from cilium identities/policies.
- Excluding job specific labels from cilium identities/policies.

## [0.10.0] - 2023-05-16

### Changed

- Enable PDB for `cilium-operator`.

## [0.9.3] - 2023-04-19

### Changed

- Revert to NetworkPolicy to allow hubble and hubble-relay egress.

## [0.9.2] - 2023-04-13

### Changed

- Change to CiliumNetworkPolicy to allow hubble and hubble-relay.

## [0.9.1] - 2023-04-13

### Added

- Add network policy to allow exposing hubble UI through ingress.

## [0.9.0] - 2023-03-20

### Changed

- Use `image.registry` value as image registry for all containers in the chart.

## [0.8.0] - 2023-03-08

### Changed

- Bump all manifests to upstream version 1.13.
- Enable Hubble
- Enable Monitoring for Agent, Operator and Hubble

## [0.7.0] - 2023-02-10

### Changed

- Enable LocalRedirectPolicy for node-local-cache and kiam.


## [0.6.1] - 2022-11-22

### Changed

- Align Helm chart ownership and CODEOWNERS file.

## [0.6.0] - 2022-11-07

### Changed

- Allow `world` access for pods in `giantswarm` namespace in default policies.
- Enable CiliumLocalRedirectPolicy

## [0.5.0] - 2022-10-18

### Changed

- Updated all templates with changes from upstream release v1.11.9

## [0.4.2] - 2022-10-14

### Fixed

- Updated healthcheck port to match new detault introduced in v1.11.6

## [0.4.1] - 2022-10-14

### Changed

- Bumped default version to v1.11.9

## [0.4.0] - 2022-10-13

### Changed

- Enable prometheus exporters for `agent`, `operator` by default.

## [0.3.1] - 2022-10-10

### Fixed

- Run `cleanup-kube-proxy-iptables` container in cilium agent in privileged mode.
- Use iptables-nft binaries for `cleanup-kube-proxy-iptables` container.

## [0.3.0] - 2022-10-06

### Added

- Add init container that cleans up iptables rules before starting cilium agent.

## [0.2.6] - 2022-07-26

### Changed

- Instead of allowing egress towards all endpoints, by default only allow access to the api server for all pods in `kube-system` and `giantswarm` namespaces.

## [0.2.5] - 2022-07-25

### Changed

- Use retagged images instead of upstream ones.
- Run the default policies creation job in hostNetwork.

## [0.2.4] - 2022-06-29

### Added

- Added the `cilium-create-default-policies` Job as a post-upgrade hook

## [0.2.3] - 2022-06-21

### Fixed

- Typo in PSP

## [0.2.2] - 2022-06-21

### Fixed

- Added missing PSP property for hubble

## [0.2.1] - 2022-06-06

### Added

- Add NetworkPolicy to allow ingress traffic towards hubble proxy.

## [0.2.0] - 2022-05-02

### Added

- Add Job to create default ingress and egress policies.
- Add PSP for hubble-relay.

## [0.1.1] - 2022-04-07

### Fixed

- Fix the version in notes.

### Added

- PodSecurityPolicies

## [0.1.0] - 2022-03-25

[Unreleased]: https://github.com/giantswarm/cilium-app/compare/v0.31.0...HEAD
[0.31.0]: https://github.com/giantswarm/cilium-app/compare/v0.30.0...v0.31.0
[0.30.0]: https://github.com/giantswarm/cilium-app/compare/v0.29.0...v0.30.0
[0.29.0]: https://github.com/giantswarm/cilium-app/compare/v0.28.0...v0.29.0
[0.28.0]: https://github.com/giantswarm/cilium-app/compare/v0.27.0...v0.28.0
[0.27.0]: https://github.com/giantswarm/cilium-app/compare/v0.26.0...v0.27.0
[0.26.0]: https://github.com/giantswarm/cilium-app/compare/v0.25.2...v0.26.0
[0.25.2]: https://github.com/giantswarm/cilium-app/compare/v0.25.1...v0.25.2
[0.25.1]: https://github.com/giantswarm/cilium-app/compare/v0.25.0...v0.25.1
[0.25.0]: https://github.com/giantswarm/cilium-app/compare/v0.24.0...v0.25.0
[0.24.0]: https://github.com/giantswarm/cilium-app/compare/v0.23.0...v0.24.0
[0.23.0]: https://github.com/giantswarm/cilium-app/compare/v0.22.0...v0.23.0
[0.22.0]: https://github.com/giantswarm/cilium-app/compare/v0.21.0...v0.22.0
[0.21.0]: https://github.com/giantswarm/cilium-app/compare/v0.20.1...v0.21.0
[0.20.1]: https://github.com/giantswarm/cilium-app/compare/v0.20.0...v0.20.1
[0.20.0]: https://github.com/giantswarm/cilium-app/compare/v0.19.2...v0.20.0
[0.19.2]: https://github.com/giantswarm/cilium-app/compare/v0.19.1...v0.19.2
[0.19.1]: https://github.com/giantswarm/cilium-app/compare/v0.19.0...v0.19.1
[0.19.0]: https://github.com/giantswarm/cilium-app/compare/v0.18.0...v0.19.0
[0.18.0]: https://github.com/giantswarm/cilium-app/compare/v0.17.0...v0.18.0
[0.17.0]: https://github.com/giantswarm/cilium-app/compare/v0.16.0...v0.17.0
[0.16.0]: https://github.com/giantswarm/cilium-app/compare/v0.15.0...v0.16.0
[0.15.0]: https://github.com/giantswarm/cilium-app/compare/v0.14.0...v0.15.0
[0.14.0]: https://github.com/giantswarm/cilium-app/compare/v0.13.0...v0.14.0
[0.13.0]: https://github.com/giantswarm/cilium-app/compare/v0.12.0...v0.13.0
[0.12.0]: https://github.com/giantswarm/cilium-app/compare/v0.11.2...v0.12.0
[0.11.2]: https://github.com/giantswarm/cilium-app/compare/v0.11.1...v0.11.2
[0.11.1]: https://github.com/giantswarm/cilium-app/compare/v0.11.0...v0.11.1
[0.11.0]: https://github.com/giantswarm/cilium-app/compare/v0.10.0...v0.11.0
[0.10.0]: https://github.com/giantswarm/cilium-app/compare/v0.9.3...v0.10.0
[0.9.3]: https://github.com/giantswarm/cilium-app/compare/v0.9.2...v0.9.3
[0.9.2]: https://github.com/giantswarm/cilium-app/compare/v0.9.1...v0.9.2
[0.9.1]: https://github.com/giantswarm/cilium-app/compare/v0.9.0...v0.9.1
[0.9.0]: https://github.com/giantswarm/cilium-app/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/giantswarm/cilium-app/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/giantswarm/cilium-app/compare/v0.6.1...v0.7.0
[0.6.1]: https://github.com/giantswarm/cilium-app/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/giantswarm/cilium-app/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/cilium-app/compare/v0.4.2...v0.5.0
[0.4.2]: https://github.com/giantswarm/cilium-app/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/giantswarm/cilium-app/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/giantswarm/cilium-app/compare/v0.3.1...v0.4.0
[0.3.1]: https://github.com/giantswarm/cilium-app/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/cilium-app/compare/v0.2.6...v0.3.0
[0.2.6]: https://github.com/giantswarm/cilium-app/compare/v0.2.5...v0.2.6
[0.2.5]: https://github.com/giantswarm/cilium-app/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/giantswarm/cilium-app/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/giantswarm/cilium-app/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/giantswarm/cilium-app/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/giantswarm/cilium-app/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/cilium-app/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/giantswarm/cilium-app/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/cilium-app/releases/tag/v0.1.0
