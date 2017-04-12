# Description

Versioning tools for Dockerfile source repos.

- `dockerfiles` generates versionsed Dockerfiles from on a common template.
- `cloudbuild` generates a configuration file to build these Dockerfiles using
  [Google Container Builder](https://cloud.google.com/container-builder/docs/).

# Installation

- Install bazel: https://bazel.build
- Clone this repo:

``` shell
git clone github.com/GoogleCloudPlatform/runtimes-common
cd runtimes-common/versioning
```

- Build:

``` shell
bazel build scripts:all
```

- Set the path to the built scripts:

``` shell
export PATH=$PATH:$PWD/bazel-bin/scripts
```

# Create `versions.yaml`

At root of the Dockerfile source repo, add a file called `versions.yaml`.
Follow the format defined in `versions.go`. See an example on
[github](https://github.com/GoogleCloudPlatform/mysql-docker).

Primary folders in the Dockerfile source repo:

- `templates` contains `Dockerfile.template`, which is a Go template for
  generating `Dockerfile`s.
- `tests` contains any tests that should be included in the generated cloud
  build configuration.
- Version folders as defined in `versions.yaml`. The `Dockerfile`s are
  generated into these folders. The folders should also contain all
  supporting files for each version, for example `docker-entrypoint.sh` files.

# Usage of `dockerfiles` command

```console
cd path/to/dockerfile/repo
dockerfiles
```

# Usage of `cloudbuild` command

```console
cd path/to/dockerfile/repo
cloudbuild > cloudbuild.yaml
```

You can use the generated `cloudbuild.yaml` file as followed:

```console
export BUCKET=<your gcs bucket>
gcloud container builds submit . \
  --config=cloudbuild.yaml \
  --verbosity=info \
  --gcs-source-staging-dir="gs://$BUCKET/staging" \
  --gcs-log-dir="gs://$BUCKET/logs"
```

Note: `BUCKET` is typically the name of the active project in your gcloud.
