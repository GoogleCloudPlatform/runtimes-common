# About

`docgen` is a tool for generating Markdown documentation for container images.

# How to install

- Install packages for Debian:
```
apt install unzip git gcc gpp libstdc++
```

- Install bazel, version 0.13.0:

```
export BAZEL_VERSION=0.13.0
curl -L -O "https://github.com/bazelbuild/bazel/releases/download/${BAZEL_VERSION}/bazel-${BAZEL_VERSION}-installer-linux-x86_64.sh"
chmod +x "bazel-${BAZEL_VERSION}-installer-linux-x86_64.sh"
"./bazel-${BAZEL_VERSION}-installer-linux-x86_64.sh" --user
aliast bazel="${HOME}/bin/bazel"
bazel version
```

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
