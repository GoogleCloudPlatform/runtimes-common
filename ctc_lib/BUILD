load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = glob(
        ["*.go"],
        exclude = glob(["*_test.go"]),
    ),
    importpath = "github.com/GoogleCloudPlatform/runtimes-common/ctc_lib",
    visibility = ["//visibility:public"],
    deps = [
        "//ctc_lib/flags:go_default_library",
        "//ctc_lib/util:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "ctc_test",
    srcs = glob(["*_test.go"]),
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
    deps = [
        #":go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
        "//ctc_lib/flags:go_default_library",
        "//ctc_lib/util:go_default_library",
    ],
)
