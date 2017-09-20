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
    data = [
        "e2e_cloudbuild.yaml",
        "e2e_test.json",
        "tiny_docker_image/Dockerfile",
        ":reconciletags.par",
    ],
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

docker_build(
    name = "docker_credential_gcr",
    base = "@python_base//image",
    directory = "/usr/local/bin",
    tars = ["@docker_credential_gcr//file"],
)

#This config file configures the client library to consult the credential helper for authenticating against GCR
docker_build(
    name = "docker_config_json",
    base = ":docker_credential_gcr",
    directory = "//root/.docker",
    files = ["config.json"],
)

py_image(
    name = "reconciletags_image",
    srcs = ["reconciletags.py"],
    base = ":docker_config_json",
    main = "reconciletags.py",
    deps = [
        "@containerregistry",
    ],
)
