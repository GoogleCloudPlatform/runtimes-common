# About

`docgen` is a tool for generating Markdown documentation for container images.

# How to

This uses [`bazel`](http://bazel.build) as the build tool.

Compile and run the tool as followed:

```
bazel build scripts:docgen
bazel-bin/scripts/docgen --spec_file path/to/your/README.yaml > README.md
```

For an example of `README.yaml` and `README.md` files, see
[mysql-docker repo](https://github.com/GoogleCloudPlatform/mysql-docker).
The yaml data follows the structure defined in
[`docgen.proto`](lib/proto/docgen.proto).
