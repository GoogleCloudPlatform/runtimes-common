# FTL Release Notes

# Version 0.6.2 - 6/21/2018
* [NODE, PHP] Fixed issue where if composer.json or package.json had values but no dependencies, FTL would error [#690](https://github.com/GoogleCloudPlatform/runtimes-common/pull/690)

# Version 0.6.1 - 6/12/2018
* [PYTHON] Fixed issue where if a requirements.txt file had comments but no dependencies FTL would fail [#680](https://github.com/GoogleCloudPlatform/runtimes-common/pull/680)

# Version 0.6.0 - 6/9/2018
* Added functionality tests for all runtimes to ensure images built with FTL are runnable.  Previously tests ensured that packages were installed but this is not sufficient [#672](https://github.com/GoogleCloudPlatform/runtimes-common/pull/672)
* Fixed issue where app file permissions were not being respected by FTL [#672](https://github.com/GoogleCloudPlatform/runtimes-common/pull/672/files#diff-68efcbd4de1f61dcfc12ab3948b88f34R41)
* Changed TTL for dependencies that are not fully specified (e.g. package.json, composer.json, requirements.txt) to be 6hr, previously 1 week [#667](https://github.com/GoogleCloudPlatform/runtimes-common/pull/667)
* [PYTHON] Fix issue where FTL errored when building python images with an empty requirementst.txt file [#669](https://github.com/GoogleCloudPlatform/runtimes-common/pull/669)
* [PYTHON] Added python support for local packages [#674](https://github.com/GoogleCloudPlatform/runtimes-common/pull/674)
* [PHP] Fixed breaking issue with PHP where apps were not able to be run properly due to autoloads.php files not being created correctly.  This was resolved by reverting FTL PHP to ‘phase 1’ caching where app installation is done with a single ‘composer install’ call [#661](https://github.com/GoogleCloudPlatform/runtimes-common/pull/661) [#672](https://github.com/GoogleCloudPlatform/runtimes-common/pull/672/files#diff-f21f9bdff4d4bbeebae09c8ae5f95448R89)


# Version 0.5.0 - 5/30/2018
* Add --sh-c-prefix flag to add a "sh -c" prefix to --entrypoint [#632](https://github.com/GoogleCloudPlatform/runtimes-common/pull/632)
* [NODE] Change Node build context to support npm local packages [#650](https://github.com/GoogleCloudPlatform/runtimes-common/pull/650)
* [PHP/PYTHON] Fixed threading issue with unsafe strptime [#645](https://github.com/GoogleCloudPlatform/runtimes-common/pull/645)
* [PYTHON] Fixed virtualenv issue with Pipfile.lock and parallel python installation [#635](https://github.com/GoogleCloudPlatform/runtimes-common/pull/635)

# Version 0.4.0 - 5/15/2018
* Fix issue where FTL was returning successful on builds w/ errors [#617](https://github.com/GoogleCloudPlatform/runtimes-common/pull/617)
* Support writing error logs to --builder-output-path or $BUILDER_OUTPUT [#613](https://github.com/GoogleCloudPlatform/runtimes-common/pull/613)
* [NODE] Changed node to cache on all descriptor files vs a single file [#630](https://github.com/GoogleCloudPlatform/runtimes-common/pull/630)
* [PYTHON] Parallelized layer uploads done in python builds [#607](https://github.com/GoogleCloudPlatform/runtimes-common/pull/607)
* [PHP] Parallelized layer uploads done in  builds [#612](https://github.com/GoogleCloudPlatform/runtimes-common/pull/612)


# Version 0.3.1 - 4/24/2018
* Fixed issue with python where if the interpreter layer was cached, other layers wouldn't build properly [#601](https://github.com/GoogleCloudPlatform/runtimes-common/pull/601)
* Added cache tests for all runtimes for add one dependency and added cleanup phase in between [#601](https://github.com/GoogleCloudPlatform/runtimes-common/pull/601)

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
