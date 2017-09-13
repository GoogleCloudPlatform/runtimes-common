package(default_visibility = ["//visibility:public"])

load("@subpar//:subpar.bzl", "par_binary")
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_prefix",
)
load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_push",
    "docker_build",
)
load("@io_bazel_rules_docker//docker/contrib/python:image.bzl", "py_image")

par_binary(
    name = "reconciletags",
    srcs = ["reconciletags.py"],
    main = "reconciletags.py",
    visibility = ["//visibility:public"],
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
    ],
)

par_binary(
    name = "config_integrity",
    srcs = ["config_integrity_test.py"],
    main = "config_integrity_test.py",
    visibility = ["//visibility:public"],
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
    ],
)

par_binary(
    name = "data_integrity",
    srcs = ["data_integrity_test.py"],
    main = "data_integrity_test.py",
    visibility = ["//visibility:public"],
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
    ],
)

py_test(
    name = "reconciletags_test",
    srcs = ["reconciletags_test.py"],
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
        "@mock",
    ],
)

py_test(
    name = "reconciletags_e2e_test",
    srcs = ["reconciletags_e2e_test.py"],
    data = ["e2e_test.json", "tiny_docker_image/Dockerfile", "cloudbuild.yaml", ":reconciletags.par"],
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
        "@mock",
    ],
)

py_test(
    name = "config_integrity_test",
    srcs = ["config_integrity_test.py"],
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
    ],
)

py_test(
    name = "data_integrity_test",
    srcs = ["data_integrity_test.py"],
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
    ],
)

py_image(
    name = "reconciletags_image",
    srcs = ["reconciletags.py"],
    base = "@python_base//image",
    main = "reconciletags.py",
    deps = [
        "@containerregistry",
    ],
)
