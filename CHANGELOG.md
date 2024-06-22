# Change Log

## [Unreleased](https://github.com/primalskill/project-template-go-http/compare/v1.3.0...HEAD)

Click on the title link to see the unreleased changes.

## [v1.3.0](https://github.com/primalskill/devcert/compare/v1.2.0...v1.3.0) (2024-06-22)

### Added

- Added changelog file.
- Added roadmap file.
- Uninstall command that removes the CA and all devcert files.

### Changed

- Dropped support for CentOS 8 and added support for CentOS Stream 9 and Rocky Linux.
- Formatted the readme.md file.
- Updated Go version and dependencies. The minimum supported Go version for builds is v1.21.

## [v1.2.0](https://github.com/primalskill/devcert/compare/v1.1.2...v1.2.0) (2023-07-16)

### Added

- Certificate info command.

### Changed

- Updated Go version and dependencies. The minimum supported Go version for builds is v1.20.

## [v1.1.2](https://github.com/primalskill/devcert/compare/v1.1.0...v1.1.2) (2022-01-16)

### Changed

- Fixed buggy regular expression when validating domain names.
- Fixed CLI messages.

## [v1.1.0](https://github.com/primalskill/devcert/compare/v1.0.0...v1.1.0) (2021-11-13)

### Added

- Linux support: Debian, Ubuntu, OpenSUSE, RHEL, CentOS, Fedora, Arch Linux.

### Changed

- Fixed bug adding "multi" suffis to the generated certificates regardless if the certificate contained multiple domains or not.
- Message formatting fixes.

## v1.0.0 (2021-11-11)

### Added

- Initial commit
