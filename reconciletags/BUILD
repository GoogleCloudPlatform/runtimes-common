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
load("@io_bazel_rules_docker//python:image.bzl", "py_image")

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

load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")

pkg_tar(
    name = "default_creds",
    srcs = [":config.json"],
    mode = "0600",
    package_dir = "/root/.docker",
)

pkg_tar(
    name = "cred_helper",
    srcs = ["@docker_credential_gcr//:docker-credential-gcr"],
    mode = "0555",
    package_dir = "/usr/local/bin",
)

docker_build(
    name = "reconcile_base",
    base = "@python_base//image",
    tars = [
        ":cred_helper.tar",
        ":default_creds.tar",
    ],
)

py_image(
    name = "reconciletags_image",
    srcs = ["reconciletags.py"],
    base = ":reconcile_base",
    main = "reconciletags.py",
    deps = [
        "@containerregistry",
    ],
)
