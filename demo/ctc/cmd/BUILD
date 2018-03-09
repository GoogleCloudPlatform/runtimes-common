load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = glob(["*.go"]),
    visibility = ["//visibility:public"],
    deps = [
        "//ctc_lib:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)
