load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "container_tool_command.go",
        "container_tool_command_base.go",
        "container_tool_list_command.go",
        "ctc_lib.go",
        "help.go",
        "version.go",
    ],
    importpath = "github.com/GoogleCloudPlatform/runtimes-common/ctc_lib",
    visibility = ["//visibility:public"],
    deps = [
        "//ctc_lib/flags:go_default_library",
        "//ctc_lib/util:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
    ],
)

go_test(
    name = "ctc_test",
    srcs = glob(["*_test.go"]),
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
    deps = [
        #":go_default_library",
        "//ctc_lib/flags:go_default_library",
        "//ctc_lib/util:go_default_library",
        "@com_github_spf13_cobra//:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["container_tool_command_test.go"],
    embed = [":go_default_library"],
    importpath = "github.com/GoogleCloudPlatform/runtimes-common/ctc_lib",
    deps = ["//vendor/github.com/spf13/cobra:go_default_library"],
)
