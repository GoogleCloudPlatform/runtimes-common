# About

`docgen` is a tool for generating Markdown documentation for container images.

# How to install

- Install [Bazel, version 0.13.0](https://bazel.build) as the build tool.

- Clone this repo:

``` shell
git clone https://github.com/GoogleCloudPlatform/runtimes-common.git
cd runtimes-common
```

- Build:

``` shell
bazel run //:gazelle
bazel build docgen/scripts/docgen:docgen
```

- Set the path to the built scripts:

``` shell
BAZEL_ARCH=linux_amd64_stripped
export PATH=$PATH:$PWD/bazel-bin/docgen/scripts/docgen/${BAZEL_ARCH}/
```

- Example:

``` shell
docgen --spec_file path/to/your/README.yaml > README.md
```

For an example of `README.yaml` and `README.md` files, see
[mysql-docker repo](https://github.com/GoogleCloudPlatform/mysql-docker).
The yaml data follows the structure defined in
[`docgen.proto`](lib/proto/docgen.proto).
