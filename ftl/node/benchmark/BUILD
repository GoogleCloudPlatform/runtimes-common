package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//python:image.bzl",
    "py_image",
)
load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_build",
)

exports_files([
    "Dockerfile",
])

py_library(
    name = "benchmark_lib",
    srcs = glob(["*.py"]),
)

py_binary(
    name = "node_benchmark",
    srcs = [
        ":benchmark_lib",
    ],
    data = ["//ftl:node_builder.par"],
    main = "ftl/node/benchmark/main.py",
    deps = [
        "@containerregistry",
    ],
)

load("@base_images_docker//dockerfile_build:dockerfile_build.bzl", "dockerfile_build")

dockerfile_build(
    name = "benchmark_base",
    base = "//ftl:node_builder_base",
    dockerfile = "//ftl/node/benchmark:Dockerfile",
)

py_image(
    name = "node_benchmark_image",
    srcs = [
        ":benchmark_lib",
    ],
    base = ":benchmark_base",
    main = "ftl/node/benchmark/main.py",
    deps = [
        "@containerregistry",
    ],
)
