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

py_binary(
    name = "reconciletags",
    srcs = glob(["*.py"]),
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
    ],
) 

py_test(
    name = "reconciletags_e2e_test",
    srcs = ["reconciletags_e2e_test.py",],
    deps = [
        "@containerregistry",
        "@httplib2//:httplib2",
    ],
    data = [
        "tiny_docker_image",
    ],

) 
