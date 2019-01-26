	Changelog
	=========
	[0.1.5] - 2019-01-25
	--------------------
	### Fixed
	- Issue #9 - Now supports `SSM.GetParametersByPath`. Naming secrets like paths will allow storing multiple values per Secret.
	- Issue #14 - Changed `log.Fatalf` to `log.Errorf` when there's a permission error on `GetParameter`. The restricted value
	will simply be skipped instead of the Pod entering a crashloop.

	[0.1.4] - 2018-09-11
	--------------------
	### Fixed
	- Add ca-certificates package to final alpine image [Issue #7]
	- Add Volume (Type=hostPath) for /etc/ssl/certs, to ensure AWS roots are available [Issue #7]

	[0.1.3] - 2018-07-24
	--------------------
	### Added
	- Docker multi-stage build (copies only the aws-ssm binary): 330MB -> 12MB

	[0.1.2] - 2018-07-23
	--------------------
	### Security
	- Removed a debug message that was dumping Secrets as %v. Plaintext values were logged
	as a series of byte values

	### Added
	- Reasonable defaults for ENV in Dockerfile

	### Other
	- Testing .travis.yml


	[0.1.1] - 2018-07-23
	--------------------
	### Added
	- Some basic tests

	### Fixed
	- Set Go @ v1.10-alpine

	### Removed [Security]
	- References to kubeconfig in Helm chart (still needed outside the cluster, though)


	[0.1.0] - 2018-07-22
	--------------------
	- Initial Release
