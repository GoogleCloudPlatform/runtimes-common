package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//docker/contrib/python:image.bzl",
    "py_image",
)
load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_push",
    "docker_build",
)
load("@subpar//:subpar.bzl", "par_binary")

py_library(
    name = "ftl_lib",
    srcs = glob([
        "__init__.py",
        "common/*.py",
    ]),
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
    ],
)
