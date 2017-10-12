package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//docker/contrib/python:image.bzl",
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
        "@mock//:mock",
    ],
)

py_library(
    name = "node_lib",
    srcs = glob(["node/*.py"]),
    deps = [":ftl_lib"],
)

# The node base image sets a default CMD, which causes issues if you run our image with no CMD.
# So, reset it to [""] before building our py_image.
docker_build(
    name = "node_builder_base",
    base = "@node_base//image",
    cmd = [""],
    entrypoint = "",
)

py_binary(
    name = "node_builder",
    srcs = [":node_lib"],
    main = "node/main.py",
    deps = ["@containerregistry"],
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
