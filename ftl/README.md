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

As an example, we will demonstrate using the Node FTL builder to create a container for a node application from a node app's source code:
Assume we are deploying the node source for the app https://github.com/JustinBeckwith/cloudcats.  First download the node ftl.par file from https://storage.googleapis.com/gcp-container-tools/ftl/node/latest/ftl.par.  Then we run the ftl.par file pointing to our application:
```shell
$ ftl.par --directory=$HOME/cloudcats/web --base=gcr.io/google-appengine/nodejs:latest --image=gcr.i
o/my-project/cloudcats-node-app:latest
```
## FTL Authentication Requirements:

FTL requires authentication with a container registry in order to properly deploy your application. The following requirements must be met for FTL to authenticate with common container registries:

1. Credentials:
In order to authenticate with a container registry, you will need to provide valid credentials. These credentials will typically include a username and password, or an access key. Depending on the container registry you are using, the specific method for providing these credentials may vary. For example, some container registries may require you to generate an access key and then use that key to authenticate, while others may require a username and password.

2. Configuration:
Once you have valid credentials, you must configure FTL with the correct container registry URL and any additional configuration options required for authentication. This configuration is typically done through the ftl.par file, which is used to deploy your application. The ftl.par file contains a variety of configuration options, including the container registry URL, the image name and tag, and any authentication-related configuration options.

3. Permissions:
In order for FTL to properly authenticate with the container registry and deploy your application, your credentials must have sufficient permissions to pull and push images to the container registry. This typically requires at least read and write access to the container registry. If you do not have the proper permissions, you may receive an error when attempting to deploy your application.

In order to properly authenticate with a container registry using FTL, you must provide valid credentials, configure FTL with the correct container registry URL and any additional authentication-related options, and ensure that your credentials have sufficient permissions to pull and push images to the container registry. By following above requirements, you can ensure a successful deployment of your application using FTL.

## Releases
Currently FTL is released in .par format for each supported runtime.  The latest release is v0.18.0, changelog [here](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/ftl/CHANGELOG.md)

### node

[v0.18.0](https://storage.googleapis.com/gcp-container-tools/ftl/node/node-v0.18.0/ftl.par)

[HEAD](https://storage.googleapis.com/gcp-container-tools/ftl/node/latest/ftl.par)

Specific version (based on git $COMMIT_SHA)
`https://storage.googleapis.com/gcp-container-tools/ftl/node/$COMMIT_SHA/ftl.par`

### python

[v0.18.0](https://storage.googleapis.com/gcp-container-tools/ftl/python/python-v0.18.0/ftl.par)

[HEAD](https://storage.googleapis.com/gcp-container-tools/ftl/python/latest/ftl.par)

Specific version (based on git $COMMIT_SHA)
`https://storage.googleapis.com/gcp-container-tools/ftl/python/$COMMIT_SHA/ftl.par`

### php
[v0.18.0](https://storage.googleapis.com/gcp-container-tools/ftl/php/php-v0.18.0/ftl.par)

[HEAD](https://storage.googleapis.com/gcp-container-tools/ftl/php/latest/ftl.par)

Specific version (based on git $COMMIT_SHA)
`https://storage.googleapis.com/gcp-container-tools/ftl/php/$COMMIT_SHA/ftl.par`

## Building and Running
FTL is built using bazel so bazel must be installed.  NOTE: FTL requires a bazel version of 0.19.1 due to syntax changes in later versions.  To build artifacts with FTL, use `bazel build` and then one of the bazel rules specified in a BUILD file.  The most common rules are `//ftl:node_builder`, `//ftl:python_builder`, and `//ftl:python_builder`.  

To build ftl.par using Bazel, the user can run the following command in the terminal:
```
bazel build //ftl:ftl.par
```
To run FTL locally, `bazel run` can be used, passing flag args to the command.  An example is below:
```
bazel run //ftl:python_builder -- \
  --base=gcr.io/google-appengine/python:latest \
  --name=gcr.io/aprindle-vm/python-ftl-v50:latest \
  --directory=$(pwd)/ftl/python/testdata/packages_test  \
  --virtualenv-dir=$HOME/env \
  --verbosity=INFO
```
FTL also supports a `--tar_base_image_path=$TARGET_PATH` flag if users do not which to upload to a registry


## Developing - Integration Tests
To run the FTL integration tests, run the following command locally from the root directory:

```shell
python ftl/ftl_<RUNTIME={node,php,python}>_integration_tests_yaml.py | gcloud builds submit --config /dev/fd/0 .
```

## FTL Runtime Design Documents
[php](https://docs.google.com/document/d/1cbf3DUpNQxdmhxo2AEhp-L34_BEsuNF8rEVGI8z7Esg/edit?usp=sharing)

[python](https://docs.google.com/document/d/1pXfg6pLPpQoIb5_E6PeVWHLttL2YgUPX6ysoqmyVzik/edit?usp=sharing)
