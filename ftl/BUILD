package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//python:image.bzl",
    "py_image",
)
load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_build",
)

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

py_library(
    name = "test_lib",
    srcs = ["common/test_util.py"],
    deps = [
        "//testing/lib:containerregistry_mock_lib",
        "@mock//:mock",
    ],
)

py_test(
    name = "context_test",
    srcs = ["common/context_test.py"],
    deps = [
        ":ftl_lib",
        "@mock//:mock",
    ],
)

py_test(
    name = "builder_test",
    srcs = ["common/builder_test.py"],
    deps = [
        ":ftl_lib",
        "@mock//:mock",
    ],
)

py_test(
    name = "cache_test",
    srcs = ["common/cache_test.py"],
    deps = [
        ":ftl_lib",
        "//testing/lib:containerregistry_mock_lib",
        "@mock//:mock",
    ],
)

py_test(
    name = "node_builder_test",
    srcs = ["node/builder_test.py"],
    data =
        glob(
            ["node/testdata/**"],
        ),
    main = "node/builder_test.py",
    deps = [
        ":node_lib",
        ":test_lib",
    ],
)

py_library(
    name = "node_lib",
    srcs = glob(["node/*.py"]),
    deps = [":ftl_lib"],
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

# The node base image sets a default CMD, which causes issues if you run our image with no CMD.
# So, reset it to [""] before building our py_image.
docker_build(
    name = "node_builder_base",
    base = "@base_images_docker//reproducible/ubuntu:node_8_5_0_ubuntu_16_0_4_build",
    cmd = [""],
    entrypoint = "",
    tars = [
        ":cred_helper.tar",
        ":default_creds.tar",
    ],
)

load("@subpar//:subpar.bzl", "par_binary")

par_binary(
    name = "node_builder",
    srcs = [":node_lib"],
    main = "node/main.py",
    deps = [
        "//testing/lib:containerregistry_mock_lib",
        "@containerregistry",
    ],
)

py_image(
    name = "node_builder_image",
    srcs = [":node_lib"],
    base = ":node_builder_base",
    main = "node/main.py",
    deps = [
        "@containerregistry",
    ],
)

py_test(
    name = "node_main_test",
    srcs = ["node/main_test.py"],
    data = [":node_builder_base.tar"],
    main = "node/main_test.py",
    deps = [
        ":node_lib",
        "//testing/lib:containerregistry_mock_lib",
        "@mock//:mock",
    ],
)
