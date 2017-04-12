load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
)

go_binary(
    name = "cloudbuild",
    srcs = ["cloudbuild/main.go"],
    deps = [
        "//versions:go_default_library",
    ],
)

go_binary(
    name = "dockerfiles",
    srcs = ["dockerfiles/main.go"],
    deps = [
        "//versions:go_default_library",
    ],
)
