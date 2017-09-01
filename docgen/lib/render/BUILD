package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "enhanced.go",
        "render.go",
    ],
    deps = [
        "//docgen/lib/proto:go_default_library",
        "//docgen/lib/render/templates:go_default_library",
    ],
)
