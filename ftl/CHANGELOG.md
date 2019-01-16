# FTL Release Notes

# Version 0.14.0 - 1/16/2018
* [PHP] Fixes a performance regression for FTL PHP where `composer update` would be called each run disregarding caching.  This removed these server side `composer update` calls.[#749](https://github.com/GoogleCloudPlatform/runtimes-common/pull/749)
* [PHP] Added `gcp-build` script support to FTL PHP [#752](https://github.com/GoogleCloudPlatform/runtimes-common/pull/752)
* [NODE] Fixed issue where FTL node would mask errro on invalid name in package.json[#755](https://github.com/GoogleCloudPlatform/runtimes-common/pull/755)
* [PYTHON] Fixed issue with recursive requirements.txt files where FTL Python would not work properly with `-r` strings in comments[#753](https://github.com/GoogleCloudPlatform/runtimes-common/pull/753)

# Version 0.13.1 - 12/12/2018
* [NODE] FTL Node caching had an issue that persisted after 0.13.0 due to a release error.  This PR resolves the initial FTL Node caching issue[#743](https://github.com/GoogleCloudPlatform/runtimes-common/pull/743)

# Version 0.13.0 - 11/27/2018
* [NODE] Updated gcp-build script execution to always run for builds instead of only running when a full rebuild is required[#739](https://github.com/GoogleCloudPlatform/runtimes-common/pull/739)
* [NODE] Fixed FTL Node issue where cache Get and Set keys could be different resulting in FTL rebuilding images that should have been cached[#738](https://github.com/GoogleCloudPlatform/runtimes-common/pull/738)

# Version 0.12.0 - 11/15/2018
* [PYTHON] Added --venv-cmd flag to FTL python which allows users to specify a full venv command for FTL to use[#735](https://github.com/GoogleCloudPlatform/runtimes-common/pull/735)

# Version 0.11.0 - 11/09/2018
* Removed verbosity from tar commands by default as it was too verbose in log output [#726](https://github.com/GoogleCloudPlatform/runtimes-common/pull/726)
* Updated FTL json error output to match GCB formatting [#727](https://github.com/GoogleCloudPlatform/runtimes-common/pull/727)
* Added --ttl flag to FTL which allows users to specify the desired TTL for FTL to respect [#732](https://github.com/GoogleCloudPlatform/runtimes-common/pull/732)
* [NODE] Fixed issue with node cache timestamps which caused node to only cache for 6 hrs a day [#731](https://github.com/GoogleCloudPlatform/runtimes-common/pull/731)


# Version 0.10.0 - 10/17/2018
* [PYTHON] Fixed issue where recursive requirements.txt files did not work properly with FTL [#721](https://github.com/GoogleCloudPlatform/runtimes-common/pull/721)
* [NODE] Fixed issue with gcp-build scripts where devDependencies were not being cleaned up properly [#723](https://github.com/GoogleCloudPlatform/runtimes-common/pull/723)
* [PHP] Updated FTL to always create a new composer-lock.json server side from the composer.json file[#724](https://github.com/GoogleCloudPlatform/runtimes-common/pull/724)


# Version 0.9.0 - 9/27/2018
* [PYTHON] Fixed issue where if a user deployed an app that had a setup.cfg, FTL would error [#717](https://github.com/GoogleCloudPlatform/runtimes-common/pull/717)
* [NODE] Changed FTL to generate a package-lock.json file server side if only a package.json is supplied [#715](https://github.com/GoogleCloudPlatform/runtimes-common/pull/715)


# Version 0.8.0 - 8/29/2018
* Added --cache-key-version flag which lets users change the entire cache-key-version (updated everytime FTL version has important cache change) [#701](https://github.com/GoogleCloudPlatform/runtimes-common/pull/701)
* Added --cache-salt flag which lets users add a salt value appended to the cache-key-version [#705](https://github.com/GoogleCloudPlatform/runtimes-common/pull/705)
* Updated FTL to output build errors as json [#702](https://github.com/GoogleCloudPlatform/runtimes-common/pull/702)
* Fixed issue where if cached layer ref was missing FTL would fail [#704](https://github.com/GoogleCloudPlatform/runtimes-common/pull/704)
* [NODE] Added support for yarn builds from a yarn.lock [#709](https://github.com/GoogleCloudPlatform/runtimes-common/pull/709)

# Version 0.7.0 - 7/9/2018
* [NODE, PHP] Fixed bug in FTL where dep files (node_modules, vendor) would be added in two layers instead of one as intended which lead to permissions and app errors [#697](https://github.com/GoogleCloudPlatform/runtimes-common/pull/697)
* Added --version flag to FTl which outputs the FTL version to stdout [#692](https://github.com/GoogleCloudPlatform/runtimes-common/pull/692)

# Version 0.6.2 - 6/21/2018
* Added --additional-directory flag which allows additional folders to be added to an FTL image [#689](https://github.com/GoogleCloudPlatform/runtimes-common/pull/689)
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
* [Python] new flags --python-cmd, --pip-cmd, and --virtualenv-cmd to support different python versions and builder container setups
* [Python] fixed issue where FTL build would fail if PYTHONPATH was not set in builder
