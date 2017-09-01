load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "command.go",
        "file_content.go",
        "file_existence.go",
        "licenses.go",
        "structure.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//structure_tests/types/unversioned:go_default_library",
        "//structure_tests/utils:go_default_library",
    ],
)
