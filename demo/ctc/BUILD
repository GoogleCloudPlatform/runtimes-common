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

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/GoogleCloudPlatform/runtimes-common/demo/ctc",
    visibility = ["//visibility:private"],
    deps = ["//demo/ctc/cmd:go_default_library"],
)

go_binary(
    name = "ctc",
    embed = [":go_default_library"],
    importpath = "github.com/GoogleCloudPlatform/runtimes-common/demo/ctc",
    visibility = ["//visibility:public"],
)
