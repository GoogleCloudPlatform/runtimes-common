package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//python:image.bzl",
    "py_image",
)
load("@base_images_docker//dockerfile_build:dockerfile_build.bzl", "dockerfile_build")

exports_files([
    "Dockerfile",
])

py_library(
    name = "benchmark_lib",
    srcs = glob(["*.py"]),
)
