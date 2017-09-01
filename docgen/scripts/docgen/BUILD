load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["docgen.go"],
    visibility = ["//visibility:private"],
    deps = [
        "//docgen/lib/render:go_default_library",
        "//docgen/lib/spec:go_default_library",
    ],
)

go_binary(
    name = "docgen",
    library = ":go_default_library",
    visibility = ["//visibility:public"],
)
