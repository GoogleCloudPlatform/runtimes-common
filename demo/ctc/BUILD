load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "ctc_demo",
    srcs = ["main.go"],
    visibility = ["//visibility:public"],
    deps = [
        "//ctc_lib:go_default_library",
        "//demo/ctc/cmd:go_default_library",
    ],
)
