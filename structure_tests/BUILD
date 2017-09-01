licenses(["notice"])

package(default_visibility = ["//visibility:public"])

sh_binary(
    name = "ext_run",
    srcs = ["ext_run.sh"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["types.go"],
    visibility = ["//visibility:private"],
    deps = [
        "//structure_tests/types/v1:go_default_library",
        "//structure_tests/types/v2:go_default_library",
    ],
)

go_binary(
    name = "structure_tests",
    library = ":go_default_library",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["structure_test.go"],
    library = ":go_default_library",
    deps = ["@com_github_ghodss_yaml//:go_default_library"],
)
