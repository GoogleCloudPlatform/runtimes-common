## FTL

"FTL" stands for "faster than light", and represents a strategy for constructing container images quickly.

In this context, the "speed of light" is considered to be the time taken to do a standard "docker build" followed by a "docker push" to a registry.

By constructing the container image layers cleverly and reproducibly, we can use the registry as a cache and speed up the build/push steps of many common language package managers.

This repository currently contains Cloud Build steps and binaries for Node.js, Python and PHP languages and package managers. 

## Usage

The typical usage of an FTL binary is:

```shell
$ ftl.par --directory=$dir --base=$base --image=$img
```

This command can be read as "Build the source code in directory `$dir` into an image named `$img`, based on the image `$base`.

These binaries **do not depend on Docker**, and construct images directly in the registry.

## Developing
To run the FTL integration tests, run the following command locally from the root directory:

```shell
python ftl/ftl_node_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 .
python ftl/ftl_php_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 .
gcloud container builds submit --config ftl/ftl_python_integration_tests.yaml .
```
