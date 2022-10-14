# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/giantswarm/cilium-app/compare/v0.4.1...HEAD
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
