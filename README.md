runtimes-common
=============

Common tools and scripts for building GCP runtimes.

* If you're looking for the container structure tests, check out our [new dedicated repo](https://github.com/GoogleCloudPlatform/container-structure-test).

You'll most likely need the `bazel` tool to build the code in this repository.
Follow these instructions to install and configure [bazel](https://bazel.build/).

We provide a pre-commit git hook for convenience.
Please install this before sending any commits via:

```shell
ln -s $(pwd)/hack/hooks/* .git/hooks/