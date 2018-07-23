	Changelog
	=========

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
