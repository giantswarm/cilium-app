# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/giantswarm/cilium-app/compare/v0.2.5...HEAD
[0.2.5]: https://github.com/giantswarm/cilium-app/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/giantswarm/cilium-app/compare/v0.2.3...v0.2.4
[0.2.3]: https://github.com/giantswarm/cilium-app/compare/v0.2.2...v0.2.3
[0.2.2]: https://github.com/giantswarm/cilium-app/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/giantswarm/cilium-app/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/giantswarm/cilium-app/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/giantswarm/cilium-app/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/cilium-app/releases/tag/v0.1.0
