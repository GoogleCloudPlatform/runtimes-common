package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//python:image.bzl",
    "py_image",
)

load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_build",
)
exports_files(["requirements.txt"])

py_library(
    name = "benchmark_lib",
    srcs = glob(["*.py"]),
)

py_binary(
    name = "node_benchmark",
    srcs = ["//ftl:node_lib", "//ftl:node_builder", ":benchmark_lib"],
    main = "ftl/node/benchmark/main.py",
    deps = [
        "@containerregistry",
    ],
)
load("@ftl_benchmark_deps//:requirements.bzl", "requirement")


py_image(
    name = "node_benchmark_image",
    srcs = ["//ftl:node_lib", "//ftl:node_builder", ":benchmark_lib"],
    base = "//ftl:node_builder_base",
    main = "ftl/node/benchmark/main.py",
    deps = [
        "@containerregistry",
        requirement("google-cloud-bigquery"),
    ],
)

