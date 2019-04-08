load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["root.go"],
    importpath = "github.com/GoogleCloudPlatform/runtimes-common/demo/ctc/cmd",
    visibility = ["//visibility:public"],
    deps = [
        "//ctc_lib:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
    ],
)
