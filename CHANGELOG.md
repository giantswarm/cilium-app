# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/giantswarm/cilium-app/compare/v0.11.1...HEAD
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
