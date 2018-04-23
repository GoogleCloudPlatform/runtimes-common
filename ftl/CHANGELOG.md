# FTL Release Notes

# Version 0.3.0 - 4/22/2018
* Added additional test suite to verify cached image/layers work appropriately for each runtime [#583](https://github.com/GoogleCloudPlatform/runtimes-common/pull/583)
* Added FTL version to logging [#579](https://github.com/GoogleCloudPlatform/runtimes-common/pull/579)
* Fixed an issue where the --no-cache flag was also not uploading artifacts [#582](https://github.com/GoogleCloudPlatform/runtimes-common/pull/582)
* Added --log-dir flag to FTL for writing log files (user/internal) that can be used in subsequent cloudbuild steps [#590](https://github.com/GoogleCloudPlatform/runtimes-common/pull/590)
* Updated all commands FTL shells out for to have better logging (stdout, stderr, return code) [#590](https://github.com/GoogleCloudPlatform/runtimes-common/pull/590)
* Added library utilities to support populating a global cache using FTL[#578](https://github.com/GoogleCloudPlatform/runtimes-common/pull/578)
* [Python] Changed python cache keys to include `python --version` output instead of --python-cmd[#588](https://github.com/GoogleCloudPlatform/runtimes-common/pull/588)
* [Python] Added configurable --virtualenv-dir flag to python[#587](https://github.com/GoogleCloudPlatform/runtimes-common/pull/587)


# Version 0.2.0 - 4/3/2018
* [PHP] fixed composer.lock parsing issue where the deps listed were being parsed incorrectly [#569](https://github.com/GoogleCloudPlatform/runtimes-common/pull/569)
* [Python] Added Pipfile.lock support to Python: using Pipfile.lock allows for per package caching (FTL Phase 2) [#554](https://github.com/GoogleCloudPlatform/runtimes-common/pull/554)
* [Python] Fixed venv directory `/bin/activate` script to have the correct path [#561](https://github.com/GoogleCloudPlatform/runtimes-common/pull/561)
* [Node] changed npm to install from a directory that is constant across builds [#572](https://github.com/GoogleCloudPlatform/runtimes-common/pull/572)

# Version 0.1.1 - 3/6/2018
* fixed error where docker metadata (exposed_ports, etc.) would not be written on an app w/ no dependencies [#]
* added --no-cache and --no-upload flags
* fixed --cache-repository flag to work as expected
* added --exposed-ports=['8090','8091'] flag to have ports exposed in output image
* fixed issue where --entrypoint was not being set properly in result image
* additional logging
* [NODE] removed auto-entrypoint detection as the default it set would override the base image default
* [PHP] added phase 2 implementation to php.  This means faster php builds for apps as packages are still cached when dependencies are changed
* [Python] added phase 1.5 implementation to python.  This means faster python builds as some layer uploading can be skipped for cache layers
* [Python] fixed issue where python was 'pip installing' each run when it should have skipped that step and used the cache
* [Python] additional logging on pip calls
* [Python] new flags --python-cmd, --pip-cmd, and --venv-cmd to support different python versions and builder container setups
* [Python] fixed issue where FTL build would fail if PYTHONPATH was not set in builder
